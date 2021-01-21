package csv

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type writer struct {
	Quote           []byte
	RecordSeparator []byte
	FieldSeparator  []byte
	ValueTrue       []byte
	ValueFalse      []byte
	writer          io.Writer
}

func NewWriter(w io.Writer) (t *writer) {
	t = new(writer)
	t.Quote = []byte{'"'}
	t.RecordSeparator = []byte{'\n'}
	t.FieldSeparator = []byte{','}
	t.ValueTrue = []byte("true")
	t.ValueFalse = []byte("false")
	t.writer = w
	return t
}

func (w *writer) WriteRow(cells ...interface{}) (err error) {
	for i, l := 0, len(cells); i < l; i++ {
		if i > 0 {
			w.writer.Write(w.FieldSeparator)
		}

		switch v := cells[i].(type) {
		case string:
			if strings.Contains(v, "\"") {
				v = strings.ReplaceAll(v, "\"", "\"\"")
			}
			fmt.Fprintf(w.writer, "\"%s\"", v)

		case time.Time:
			fmt.Fprintf(w.writer, "\"%s\"", v)

		case bool:
			if v {
				_, err = w.writer.Write(w.ValueTrue)
			} else {
				_, err = w.writer.Write(w.ValueFalse)
			}

		default:
			fmt.Fprintf(w.writer, "%s", v)
		}
	}
	_, err = w.writer.Write(w.RecordSeparator)
	return
}
