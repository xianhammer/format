package excel

import (
	"bufio"
	"bytes"
	"io"

	"github.com/xianhammer/format/xml"
)

type saxChildAttributes struct {
	saxNodeAttributes
	parentTag []byte
	active    bool
}

func (s *saxChildAttributes) Tag(name []byte) {
	switch name[0] {
	case '/': // Close tag
		if s.active && bytes.Equal(s.parentTag, name[1:]) {
			s.active = false
		}
	case '?': // ProcessingInstruction - ignore.
	default: // Open tag
		if !s.active && bytes.Equal(s.parentTag, name) {
			s.active = true
		}
	}
}

func (s *saxChildAttributes) TagEnd(autoclose bool) {
	if s.active {
		s.saxNodeAttributes.TagEnd(autoclose)
	}
}

func (s *saxChildAttributes) Attribute(tag, name, value []byte) {
	if s.active {
		s.saxNodeAttributes.Attribute(tag, name, value)
	}
}

func ChildAttributes(r io.Reader, tag, parentTag []byte) (attributes Elements, err error) {
	s := new(saxChildAttributes)
	s.parentTag = parentTag
	s.tag = tag

	t := xml.NewTokenizer(s)
	if _, err = t.ReadFrom(bufio.NewReader(r)); err == io.EOF {
		err = nil
	}

	return s.values, err
}
