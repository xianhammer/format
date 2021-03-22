package cfb

import (
	"io"
)

type Fat struct {
	s      []Sector
	offset uint32
	doc    *Document
}

func (f *Fat) add(sect Sector) {
	f.s = append(f.s, sect)
}

func (f *Fat) Sectors() uint32 {
	return uint32(len(f.s))
}

func (f *Fat) Len() uint32 {
	return f.Sectors() * f.doc.sectorSize
}

func (f *Fat) Seek(offset int64, whence int) (n int64, err error) {
	size := f.Len()
	pos := f.offset
	switch whence {
	case io.SeekStart:
		pos = uint32(offset)
	case io.SeekCurrent:
		pos += uint32(offset)
	case io.SeekEnd:
		pos = size - uint32(offset)
	}

	if 0 <= pos && pos < size {
		f.offset = pos
		n = int64(pos)
	} else {
		err = ErrSeekIndex
	}

	return
}

// func (f *Fat) Read(p []byte) (n int, err error) {
// 	sectorSize := s.doc.sectorSize

// 	sID := s.offset / sectorSize
// 	offset := s.offset % sectorSize

// 	if sID >= uint32(len(s.s)) {
// 		return 0, io.EOF
// 	}

// 	n = copy(p, s.s[sID][offset:])
// 	s.offset += uint32(n)
// 	return
// }

func (f *Fat) Close() (err error) {
	return
}
