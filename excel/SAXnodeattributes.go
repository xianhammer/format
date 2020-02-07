package excel

import (
	"bufio"
	"bytes"
	"io"

	"github.com/xianhammer/format/xml"
)

type saxNodeAttributes struct {
	xml.Partial
	tag     []byte
	current Attributes
	values  Elements
}

func (s *saxNodeAttributes) Tag(name []byte) {
	switch name[0] {
	case '/': // Close tag
		s.current = nil
	case '?': // ProcessingInstruction - ignore.
		s.current = nil
	default: // Open tag
		if bytes.Equal(s.tag, name) {
			s.current = make(map[string]string)
			s.values = append(s.values, s.current)
		}
	}
}

func (s *saxNodeAttributes) TagEnd(autoclose bool) {
	s.current = nil
}

func (s *saxNodeAttributes) Attribute(tag, name, value []byte) {
	if s.current != nil {
		s.current[string(name)] = string(value)
	}
}

func NodeAttributes(r io.Reader, tag []byte) (attributes Elements, err error) {
	s := new(saxNodeAttributes)
	s.tag = tag

	t := xml.NewTokenizer(s)
	if _, err = t.ReadFrom(bufio.NewReader(r)); err == io.EOF {
		err = nil
	}

	return s.values, err
}
