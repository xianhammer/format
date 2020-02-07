package xml

import (
	"io"

	"github.com/xianhammer/format/parse"
)

// Receiver interface for SAX style callback.
type Callback interface {
	// Tag is called when a tag start '<' is met. The name contain any namespace value too, e.g "ns:tag"
	// Note, this is also called on end tags, such as "</br>" with name being "/br".
	Tag(name []byte)
	// TagEnd is called when the closing '>' is met. If the byte prior to '>' is '/' autoclose is true.
	TagEnd(autoclose bool)
	// Attribute is called when a full attribute name, value pair has been parsed.
	Attribute(tag, name, value []byte)

	// Text represent any value not inside '<' and '>'.
	Text(value []byte)
	// Comment content is returned.
	Comment(value []byte)
	// Section represent CDATA section content.
	Section(name []byte)
}

// Partial is a simple implementation of the Callback interface.
// No methods does anything, so this object can be used to struct extend
type Partial int

func (p Partial) Tag(name []byte)                   {}
func (p Partial) TagEnd(autoclose bool)             {}
func (p Partial) Attribute(tag, name, value []byte) {}
func (p Partial) Text(value []byte)                 {}
func (p Partial) Comment(value []byte)              {}
func (p Partial) Section(name []byte)               {}

// Tokenizer contain data for XML stream tokenizing.
type Tokenizer struct {
	receiver        Callback
	previous        byte          // Previous byte read.
	readUntil       byte          // Marker used to skip bytes until the given is met.
	specialTagLevel int           // Count bytes like '-' (Comment) and '[' (CDATA etc).
	colons          int           // Count colons met
	TagName         *parse.Buffer // TagName buffer for tag names (including any prefixed namespace and colon).
	AttrName        *parse.Buffer // AttrName buffer for attribute names (including any prefixed namespace and colon).
	Value           *parse.Buffer // Value buffer for attribute values.
}

// NewTokenizer creates a new tokenizer with the given Callback based receiver.
// The three (internal) buffers, TagName, attrName and Value are initialised to sizes 128, 128 resp. 32K.
func NewTokenizer(receiver Callback) (t *Tokenizer) {
	t = new(Tokenizer)
	t.receiver = receiver
	t.readUntil = '<'

	t.TagName = parse.NewBuffer(128)
	t.AttrName = parse.NewBuffer(128)
	t.Value = parse.NewBuffer(32 * 1024)
	return
}

// ReadFrom tokenize a stream of bytes as XML items.
// The read function is (intended to be) recoverable and can be re-called if an error is returned.
func (t *Tokenizer) ReadFrom(r io.Reader) (n int64, err error) {
	b := []byte{t.previous}

	next := func() {
		if err == nil {
			n++
			_, err = r.Read(b)
		}
	}

	var bIdentifier *parse.Buffer
	if t.previous == 0 {
		n--
		bIdentifier = t.TagName
		next()
	}

	for ; err == nil; next() {
		if t.readUntil != 0 {
			if b[0] != t.readUntil {
				err = t.Value.Push(b[0])
				continue
			}

			if t.specialTagLevel >= 1 {
				t.specialTagLevel--
				if t.specialTagLevel != 0 {
					continue
				}
			}

			t.previous = 0
			switch t.readUntil {
			case '"':
				t.receiver.Attribute(t.TagName.GetData(), t.AttrName.FetchData(), t.Value.FetchData())
				t.readUntil = 0
				continue

			case '-':
				t.receiver.Comment(t.Value.FetchData())
				t.readUntil = 0
				continue

			case ']':
				t.receiver.Text(t.Value.FetchData())
				t.readUntil = 0
				continue

			case '>':
				t.receiver.Text(t.Value.FetchData())

			case '<':
				if !t.Value.Empty() {
					t.receiver.Text(t.Value.FetchData())
				}
			}
			t.readUntil = 0
		}

		switch b[0] {
		case '?': // ProcessingInstruction
			if t.previous == '<' {
				err = t.TagName.Push(b[0])
			}

		case '!': // Comment, CharData|CDATA, DOCTYPE
			if t.previous != '<' {
				err = ErrIllegalCharacter
			}

			t.specialTagLevel = 1 // Mark the tokenizer may be in a comment section

		case '-', '[':
			if t.specialTagLevel == 0 {
				err = ErrIllegalCharacter
			} else if t.specialTagLevel == 1 {
				t.specialTagLevel++
			} else if b[0] == '[' {
				t.receiver.Section(t.TagName.FetchData())
				t.readUntil = ']'
			} else {
				t.Value.Clear()
				t.readUntil = '-'
			}

		case '<':
			t.TagName.Clear()
			t.AttrName.Clear()
			t.Value.Clear()
			t.colons = 0
			bIdentifier = t.TagName

		case '>':
			if bIdentifier == t.TagName && !t.TagName.Empty() {
				if t.previous == ':' {
					err = ErrUnterminatedIdentifier
					continue
				}
				t.receiver.Tag(t.TagName.GetData())
			}
			t.receiver.TagEnd(t.previous == '/')

			t.Value.Clear()
			t.readUntil = '<' // Read Text element

		case '/':
			if t.previous == '<' {
				err = t.TagName.Push(b[0])
			} else if t.previous == ':' {
				err = ErrUnterminatedIdentifier
			}

		case '=':
			if t.AttrName.Empty() {
				err = ErrIllegalCharacter
			}
			t.Value.Clear()

		case '"':
			if t.previous != '=' {
				err = ErrIllegalCharacter
			}

			t.readUntil = '"' // Read attr-value

		case ':':
			if t.colons >= 1 || bIdentifier.Empty() {
				err = ErrBadIdentifier
			} else {
				t.colons++
				err = bIdentifier.Push(b[0])
			}

		default:
			if t.previous == '=' {
				err = ErrIllegalCharacter
			} else if cs := charset[b[0]]; cs&csWhitespace != 0 {
				// Met a whitespace - transmit any tag (start) registered
				if bIdentifier == t.TagName && !bIdentifier.Empty() {
					if t.specialTagLevel == 1 {
						t.receiver.Section(bIdentifier.FetchData())
						t.readUntil = '>'
					} else {
						t.receiver.Tag(t.TagName.GetData())
					}
				}
				t.AttrName.Clear()
				bIdentifier = t.AttrName
				t.colons = 0
				t.specialTagLevel = 0
			} else if cs&csIdentifierStart == 0 && t.previous == ':' {
				err = ErrBadIdentifier
			} else if cs&csIdentifierFollow != 0 && !bIdentifier.Empty() {
				err = bIdentifier.Push(b[0])
			} else if cs&csIdentifierStart != 0 {
				// Met a start identifier char.
				t.colons = 0
				err = bIdentifier.Push(b[0])
			} else {
				err = ErrIllegalCharacter
			}
		}

		t.previous = b[0]
	}

	if err != io.EOF {
		return
	}

	if t.readUntil != '<' {
		if t.readUntil == '"' {
			err = ErrUnterminatedString
		} else {
			err = ErrUnterminatedTag
		}
		return
	}

	if !t.Value.Empty() {
		t.receiver.Text(t.Value.FetchData())
	}

	return
}
