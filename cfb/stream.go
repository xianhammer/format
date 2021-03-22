package cfb

import (
	"bytes"
	"encoding/binary"
	"io"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var (
	win16be  = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	utf16bom = unicode.BOMOverride(win16be.NewDecoder())
)

type Stream struct {
	s          []Sector
	offset     uint32
	size       uint32
	sectorSize uint32
	byteorder  binary.ByteOrder
}

func NewStream(sectorSize, streamSize uint32) (s *Stream) {
	return &Stream{
		sectorSize: sectorSize,
		size:       streamSize,
		byteorder:  binary.LittleEndian,
	}
}

func (s *Stream) add(sect Sector, addSize bool) {
	s.s = append(s.s, sect)
	if addSize {
		s.size += uint32(len(sect))
	}
}

func (s *Stream) Sectors() uint32 {
	return uint32(len(s.s))
}

func (s *Stream) Len() uint32 {
	if s.size > 0 {
		return s.size
	}
	return uint32(len(s.s)) * s.sectorSize
}

func (s *Stream) Seek(offset int64, whence int) (n int64, err error) {
	size := s.Len()
	pos := s.offset
	switch whence {
	case io.SeekStart:
		pos = uint32(offset)
	case io.SeekCurrent:
		pos += uint32(offset)
	case io.SeekEnd:
		pos = size - uint32(offset)
	}

	if 0 <= pos && pos < size {
		s.offset = pos
		n = int64(pos)
	} else {
		err = ErrSeekIndex
	}

	return
}

func (s *Stream) Read(dst []byte) (n int, err error) {
	if s.offset >= s.size {
		return 0, io.EOF
	}

	sID := s.offset / s.sectorSize
	if 0 < s.size && s.size <= s.offset {
		return 0, io.EOF
	}

	src := s.s[sID][s.offset%s.sectorSize:]
	if (sID+1) == uint32(len(s.s)) && s.size-s.offset < s.sectorSize {
		src = src[:s.size-s.offset]
	}

	n = copy(dst, src)

	s.offset += uint32(n)
	return
}

func (s *Stream) Close() (err error) {
	return
}

func (s *Stream) AsUnicode() (r io.Reader) {
	return transform.NewReader(s, utf16bom)
}

func (s *Stream) ReadUnicode() (v string, err error) {
	r := transform.NewReader(s, utf16bom)

	var b bytes.Buffer
	if _, err = io.Copy(&b, r); err == nil {
		v = b.String()
	}
	return
}

func (s *Stream) ReadString() (v string, err error) {
	var b bytes.Buffer
	if _, err = io.Copy(&b, s); err == nil {
		v = b.String()
	}
	return
}

func (s *Stream) ReadUint8() (v uint8, err error) {
	b := []byte{0}
	if _, err = io.ReadFull(s, b); err == nil {
		v = b[0]
	}
	return
}

func (s *Stream) ReadUint16() (v uint16, err error) {
	b := []byte{0, 0}
	if _, err = io.ReadFull(s, b); err == nil {
		v = s.byteorder.Uint16(b)
	}
	return
}

func (s *Stream) ReadUint32() (v uint32, err error) {
	b := []byte{0, 0, 0, 0}
	if _, err = io.ReadFull(s, b); err == nil {
		v = s.byteorder.Uint32(b)
	}
	return
}

func (s *Stream) ReadUint64() (v uint64, err error) {
	b := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	if _, err = io.ReadFull(s, b); err == nil {
		v = s.byteorder.Uint64(b)
	}
	return
}
