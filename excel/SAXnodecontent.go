package excel

import (
	"bufio"
	"bytes"
	"io"

	"github.com/xianhammer/format/xml"
)

type saxNodeContent struct {
	xml.Partial
	tag    []byte
	values []string
	active bool
}

var empty = []byte{}

func (s *saxNodeContent) appendValue(value []byte) {
	s.values = append(s.values, string(value))
	s.active = false
}

func (s *saxNodeContent) Tag(name []byte) {
	switch name[0] {
	case '/': // Close tag
		if s.active {
			s.appendValue(empty)
		}
	default: // Open tag
		s.active = bytes.Equal(s.tag, name)
	}
}

func (s *saxNodeContent) TagEnd(autoclose bool) {
	if autoclose && s.active {
		s.appendValue(empty)
	}
}

func (s *saxNodeContent) Text(value []byte) {
	if s.active {
		s.appendValue(value)
	}
}

func NodeContent(r io.Reader, tag []byte) (contents []string, err error) {
	s := new(saxNodeContent)
	s.tag = tag

	t := xml.NewTokenizer(s)
	if _, err = t.ReadFrom(bufio.NewReader(r)); err == io.EOF {
		err = nil
	}

	return s.values, err
}
