package pdf

import (
	"bufio"
	"io"
)

type XRefEntry struct {
	offset     int
	generation int
	kind       byte
}

type XRef struct {
	offset int64       // offset of xref table in PDF file.
	Entry  []XRefEntry // map[int]XRefEntry
}

func NewXRef(offset int64) *XRef {
	return &XRef{offset, nil} //make(map[int]XRefEntry)}
}

func (xref *XRef) Read(r io.Reader) (n int64, err error) {
	const LineSize = 20
	br := bufio.NewReader(r)

	var b []byte

	// SKIP "xref<eol>"
	b, err = br.Peek(LineSize)

	i := 0
	for ; i < LineSize && !(b[i] == '\n' || b[i] == '\r'); i++ {
	}
	for ; i < LineSize && (b[i] == '\n' || b[i] == '\r'); i++ {
	}
	br.Discard(i)

	for {
		b, err = br.Peek(LineSize)
		if err != nil {
			break
		}

		var indexOffset, count, n0 int
		for i = 0; i < LineSize && (b[i]-'0') < 10; i++ {
			indexOffset = 10*indexOffset + int(b[i]&0x0f)
		}

		if i >= LineSize || b[i] != ' ' {
			break
		}
		for i++; i < LineSize && (b[i]-'0') < 10; i++ {
			count = 10*count + int(b[i]&0x0f)
		}

		if i >= LineSize || !(b[i] == '\n' || b[i] == '\r') {
			break
		}

		for i++; i < LineSize && (b[i] == '\n' || b[i] == '\r'); i++ {
		}
		br.Discard(i)

		diff := (indexOffset + count) - len(xref.Entry)
		if diff > 0 {
			xref.Entry = append(xref.Entry, make([]XRefEntry, diff)...)
		}

		buffer := make([]byte, 20*count)
		n0, err = io.ReadFull(br, buffer)
		if err != nil {
			break
		}

		n += int64(n0)
		for j := 0; j < count; j++ {
			var eOffset, eGeneration int
			k, b0 := 0, buffer[j*20:]
			for ; k < 10 && (b0[k]-'0') < 10; k++ {
				eOffset = 10*eOffset + int(b0[k]&0x0f)
			}
			if k != 10 || b0[10] != ' ' {
				err = ErrInvalidXRefEntry
				break
			}
			for k++; k < 16 && (b0[k]-'0') < 10; k++ {
				eGeneration = 10*eGeneration + int(b0[k]&0x0f)
			}
			if k != 16 || b0[16] != ' ' {
				err = ErrInvalidXRefEntry
				break
			}

			// TODO? Check EOL...
			xref.Entry[indexOffset+j] = XRefEntry{eOffset, eGeneration, b0[17]}
		}
	}

	return
}
