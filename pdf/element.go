package pdf

import (
	"bufio"
	"fmt"
	"io"
)

type Kind int

const ( // See to table I.1 in pdf_reference_1-7.pdf
	Null Kind = iota
	Integer
	Real
	Boolean
	Name
	String
	Dictionary
	Array
	Stream
	Reference // Not in table I.1 - used for (unresolved) element references.
)

type Element struct {
	// kind       Kind
	start, end int64
}

func (e *Element) String() string {
	return fmt.Sprintf("<%d, %d>", e.start, e.end)
}

type Object struct {
	Entry map[string]*Element
}

func NewObject() (o *Object) {
	o = new(Object)
	o.Entry = make(map[string]*Element)
	return
}

func (o *Object) Read(r io.Reader) (n int64, err error) {
	// For referenced objects, the xref position point at the object ID in [object-id] [generation] "obj" ...
	// However, the "trailer" object has no id+generation, only the "trailer" keyword.
	// To embrace all cases, this read simply skip until first '<<'.
	br := bufio.NewReader(r)

	var b []byte
	var c byte
	b, err = br.ReadBytes('<')
	n = int64(len(b))
	if err != nil {
		return
	}

	c, err = br.ReadByte()
	n++
	if err != nil {
		return
	}
	if c != '<' {
		return int64(len(b) + 1), ErrMissingObjectStart
	}
	n++

	var levels []int
	var literalName []byte
	var currentValue *Element
	var previousValue *Element

	levels = append(levels, 0) // Root value
	inLiteralName := false
	inHexString := false
	for err == nil && len(levels) > 0 {
		c, err = br.ReadByte()
		n++
		if err != nil {
			break
		}

		if inLiteralName {
			inLiteralName = characters[c]&c_literalNameTerm == 0
			if inLiteralName {
				literalName = append(literalName, c)
				continue
			}

			currentValue = new(Element)
			o.Entry[string(literalName)] = currentValue
			currentValue.start = n
		} else if inHexString {
			inHexString = c != '>'
			continue
		}

		if currentValue != nil {
			previousValue = currentValue
			currentValue = nil
		}

		switch c {
		case '>': // Can mark end of hex string too.
			c, err = br.ReadByte()
			n++
			if err == nil && c == '>' {
				previousTerm := levels[len(levels)-1]
				levels = levels[:len(levels)-1]
				literalName = literalName[:previousTerm]
				if previousValue != nil {
					previousValue.end = n
					previousValue = nil
				}
			} else if err == nil {
				err = ErrUnexpectedCharacter
			}

		case '<': // Can mark start of hex-string too.
			c, err = br.ReadByte()
			n++
			if err == nil && c == '<' {
				levels = append(levels, len(literalName))
				if previousValue != nil {
					previousValue.end = n
					previousValue = nil
				}
			} else if err == nil {
				inHexString = true
			}

		case '/':
			resetTerm := levels[len(levels)-1]
			literalName = append(literalName[:resetTerm], c)
			inLiteralName = true

			if previousValue != nil {
				previousValue.end = n
				previousValue = nil
			}
		}
	}

	return
}
