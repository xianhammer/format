package parse

import (
	"strings"
	"testing"
)

func TestBuffer(t *testing.T) {
	bufferSize := 5

	b := NewBuffer(bufferSize)
	if b == nil {
		t.Fatalf("Buffer creation error, expected pointer, got [%v]\n", b)
	}

	if len(b.GetData()) != 0 {
		t.Errorf("Non-empty buffer created.\n")
	}
	if !b.Empty() {
		t.Errorf("Non-empty buffer created.\n")
	}
	if b.Full() {
		t.Errorf("Created buffer registered as full.\n")
	}

	err := b.Push(1)
	if err != nil {
		t.Errorf("Buffer error, expected [%v], got [%v].\n", nil, err)
	}
	if len(b.GetData()) != 1 {
		t.Errorf("Buffer length error, expected [%v], got [%v].\n", 1, len(b.GetData()))
	} else if b.GetData()[0] != 1 {
		t.Errorf("Buffer content error, expected [%v], got [%v].\n", 1, b.GetData()[0])
	}
	if b.Empty() {
		t.Errorf("Empty buffer met, expected content.\n")
	}
	if b.Full() {
		t.Errorf("Non-empty buffer registered as full.\n")
	}

	err = b.Push(1)
	if err != nil {
		t.Errorf("Buffer error, expected [%v], got [%v].\n", nil, err)
	}
	if len(b.GetData()) != 2 {
		t.Errorf("Buffer length error, expected [%v], got [%v].\n", 2, len(b.GetData()))
	}
	if b.Empty() {
		t.Errorf("Empty buffer met, expected content.\n")
	}
	if b.Full() {
		t.Errorf("Non-empty buffer registered as full.\n")
	}

	fetched := b.FetchData()
	if len(fetched) != 2 {
		t.Errorf("Buffer length error, expected [%v], got [%v].\n", 2, len(fetched))
	}
	if !b.Empty() {
		t.Errorf("Empty buffer met, expected content.\n")
	}
	if b.Full() {
		t.Errorf("Non-empty buffer registered as full.\n")
	}

	b.Clear()
	if !b.Empty() {
		t.Errorf("Non-empty buffer met.\n")
	}
	if b.Full() {
		t.Errorf("Non-empty buffer registered as full.\n")
	}

	for i := 0; i < bufferSize; i++ {
		b.Push(byte(i))
	}
	if b.Empty() {
		t.Errorf("Filled buffer is empty.\n")
	}
	if !b.Full() {
		t.Errorf("Filled buffer is not full.\n")
	}

	for i, v := range b.GetData() {
		if i != int(v) {
			t.Errorf("Buffer value error, expected [%v], got [%v].\n", i, v)
		}
	}

	err = b.Push(1)
	if err != ErrOutOfBounds {
		t.Errorf("Buffer error, expected [%v], got [%v].\n", ErrOutOfBounds, err)
	}
}

func TestBuffer1(t *testing.T) {
	input := "some io.Reader stream to be read\n"

	b := NewBuffer(len(input))
	if b == nil {
		t.Fatalf("Buffer creation error, expected pointer, got [%v]\n", b)
	}

	r := strings.NewReader(input)
	n, err := b.ReadFrom(r)
	if err != nil {
		t.Fatalf("Buffer readfrom error [%v]\n", err)
	}
	if n != int64(len(input)) {
		t.Fatalf("Buffer readfrom read length error, expected [%v], got [%v]\n", len(input), n)
	}

	var b0 strings.Builder
	n, err = b.WriteTo(&b0)
	if err != nil {
		t.Fatalf("Buffer writeto error [%v]\n", err)
	}
	if n != int64(len(input)) {
		t.Fatalf("Buffer writeto write length error, expected [%v], got [%v]\n", len(input), n)
	}
	if b0.String() != input {
		t.Fatalf("Buffer writeto write error, expected [%v], got [%v]\n", input, b0.String())
	}
}
