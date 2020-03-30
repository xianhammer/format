package oxmsg

import (
	"io"
	"math"

	"github.com/xianhammer/format/cbf"
)

type Entry struct {
	*cbf.DirectoryEntry
	property        *Property
	isKnown         bool
	interpretedName string
}

func newEntry(d *cbf.DirectoryEntry) (e *Entry) {
	e = new(Entry)
	e.DirectoryEntry = d

	var err error
	if e.property, e.isKnown, err = ParseProperty(d.Name()); err == nil {
		e.interpretedName = e.property.Name
	} else {
		e.interpretedName = d.Name()
	}
	return
}

func (e *Entry) Name() string {
	return e.interpretedName
}

func (e *Entry) Uint8() (v uint8, err error) {
	s, err := e.Stream()
	if err != nil {
		return
	}
	return s.ReadUint8()
}

func (e *Entry) Uint16() (v uint16, err error) {
	s, err := e.Stream()
	if err != nil {
		return
	}
	return s.ReadUint16()
}

func (e *Entry) Uint32() (v uint32, err error) {
	s, err := e.Stream()
	if err != nil {
		return
	}
	return s.ReadUint32()
}

func (e *Entry) Float32() (v float32, err error) {
	s, err := e.Stream()
	if err != nil {
		return
	}

	bits, err := s.ReadUint32()
	if err == nil {
		v = math.Float32frombits(bits)
	}
	return
}

func (e *Entry) Float64() (v float64, err error) {
	s, err := e.Stream()
	if err != nil {
		return
	}

	bits, err := s.ReadUint64()
	if err == nil {
		v = math.Float64frombits(bits)
	}
	return
}

func (e *Entry) String() (v string, err error) {
	s, err := e.Stream()
	if err != nil {
		return
	}

	if e.property.Type == PtypString {
		return s.ReadUnicode()
	}
	return s.ReadString()
}

func (e *Entry) TypedReader() (r io.Reader, err error) {
	stream, err := e.Stream()
	if err != nil || e.property == nil {
		return nil, err
	}
	// if err != nil || e.property == nil {
	// 	r = stream
	// 	return
	// }

	switch e.property.Type {
	case PtypString:
		r = stream.AsUnicode()
	case PtypString8:
		r = stream
	default:
		r = stream
		// TODO Implement other type readers
	}
	return
}
