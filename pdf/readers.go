package pdf

import (
	"bufio"
	"bytes"
	"io"
)

// func ReadWhile(r *bufio.Reader, mask Mask, readAhead int, characters []Mask) (body *bytes.Buffer, err error) {
// 	return ReadUntil(r, c_negate^mask, readAhead, characters)
// }

func ReadUntil(r *bufio.Reader, mask Mask, readAhead int, characters []Mask) (body *bytes.Buffer, err error) {
	body = new(bytes.Buffer)

	var b []byte
	for i := 0; i >= len(b); i = 0 {
		b, err = r.Peek(readAhead)
		if err != nil && err != io.EOF {
			break
		}

		for ; i < len(b) && characters[b[i]]&mask == 0; i++ {
		}

		body.Write(b[:i])
		r.Discard(i)
	}
	return
}

func Skip(r io.Reader, toSkip int64) (n int64, err error) {
	p := make([]byte, 1024)
	for skipped := 0; toSkip > 0 && err == nil; n += int64(skipped) {
		if int64(len(p)) > toSkip {
			skipped, err = r.Read(p[:toSkip])
		} else {
			skipped, err = r.Read(p)
		}
		toSkip -= int64(skipped)
	}
	return
}

type ByteCountReader struct {
	count int64
	r     io.Reader
}

// NewByteCountReader implement a io.Reader counting the number of bytes passed through.
func NewByteCountReader(r io.Reader, offset int64) *ByteCountReader {
	return &ByteCountReader{offset, r}
}

// Count return number of bytes read so far.
func (bc *ByteCountReader) Count() int64 {
	return bc.count
}

// Move the current count 'offset' ahead (or back)
func (bc *ByteCountReader) Move(offset int64) *ByteCountReader {
	bc.count += offset
	return bc
}

// Read - satisfy io.Reader interface.
func (bc *ByteCountReader) Read(p []byte) (n int, err error) {
	n, err = bc.r.Read(p)
	bc.count += int64(n)
	return
}
