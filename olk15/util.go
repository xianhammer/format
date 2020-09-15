package olk15

import (
	"bytes"
	"unicode/utf16"
	"unicode/utf8"
)

// import "golang.org/x/text/encoding/charmap"

func DecodeUTF8(b []byte) string {
	u16s := []uint16{0}
	b8buf := []byte{0, 0, 0, 0}

	ret := &bytes.Buffer{}
	lb := len(b)
	for i := 0; i < lb; i += 2 {
		// Big endian: u16s[0] = uint16(b[i+1]) + (uint16(b[i]) << 8)
		u16s[0] = uint16(b[i]) + (uint16(b[i+1]) << 8)
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		ret.Write(b8buf[:n])
	}

	return ret.String()
}

func defaultDecoder(b []byte) (s string) {
	return string(b)
}
