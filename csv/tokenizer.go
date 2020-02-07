package csv

import (
	"io"

	"github.com/xianhammer/format/parse"
)

type Callback interface {
	Field(column, line uint, value []byte)
}

type Tokenizer struct {
	Quote            byte
	FieldSeparator   byte
	RecordSeparator  byte //= '\n'
	LineComment      byte //= '#'
	IgnoreLinesUntil uint
	receiver         Callback
	Buffer           *parse.Buffer
}

func NewTokenizer(receiver Callback) (t *Tokenizer) {
	t = new(Tokenizer)
	t.receiver = receiver
	t.Quote = '"'
	t.RecordSeparator = '\n'
	t.FieldSeparator = ','
	t.IgnoreLinesUntil = 0
	t.LineComment = 0
	return
}

func (t *Tokenizer) ReadFrom(r io.Reader) (n int64, err error) {
	if t.Buffer == nil {
		t.Buffer = parse.NewBuffer(1024)
	}

	b := []byte{0}

	var line, column uint
	for {
		if _, err = r.Read(b); err != nil {
			return
		}

		n++
		if b[0] == t.Quote {
			for quote := 1; err == nil && b[0] == t.Quote; quote++ {
				if _, err = r.Read(b); err != nil {
					return
				}

				n++
				if quote%2 == 0 {
					t.Buffer.Push(t.Quote)
				}
			}
		}

		switch b[0] {
		case t.LineComment:
			for err == nil && b[0] != t.RecordSeparator {
				if _, err = r.Read(b); err != nil {
					return
				}
				n++
			}
			line++

		case t.FieldSeparator:
			t.receiver.Field(column, line, t.Buffer.FetchData())
			column++

		case t.RecordSeparator:
			t.receiver.Field(column, line, t.Buffer.FetchData())
			line++
			column = 0

		default:
			t.Buffer.Push(b[0])
		}
	}
}
