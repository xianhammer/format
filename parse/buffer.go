package parse

import "io"

type Buffer struct {
	idx, size int64
	data      []byte
}

// NewBuffer creates a new bounded buffer for tokenizers.
// The buffer support byte appending
func NewBuffer(size int) *Buffer {
	b := new(Buffer)
	b.size = int64(size)
	b.data = make([]byte, size)
	return b
}

// Push a byte on the buffer. Return ErrOutOfBounds if overflow.
func (b *Buffer) Push(c byte) (err error) {
	if b.idx == b.size {
		return ErrOutOfBounds
	}
	b.data[b.idx] = c
	b.idx++
	return
}

// FetchData return the data currently written then clear the buffer.
func (b *Buffer) FetchData() (d []byte) {
	d = b.data[:b.idx]
	b.idx = 0
	return
}

// GetData return the data currently written.
func (b *Buffer) GetData() []byte {
	return b.data[:b.idx]
}

// Clear the buffer.
func (b *Buffer) Clear() {
	b.idx = 0
}

// Empty return true if the buffer is empty.
func (b *Buffer) Empty() bool {
	return b.idx == 0
}

// Full return true if the buffer is full.)
func (b *Buffer) Full() bool {
	return b.idx == b.size
}

func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error) {
	n0, err := r.Read(b.data[b.idx:])
	n = int64(n0)
	b.idx += n
	return
}

func (b *Buffer) WriteTo(w io.Writer) (n int64, err error) {
	n0, err := w.Write(b.data[:b.idx])
	n = int64(n0)
	return
}
