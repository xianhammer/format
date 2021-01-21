package csv

import (
	"io"
	"os"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var (
	ForExcel = unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
)

func NewEncodedWriter(wIn io.WriteCloser, encoder encoding.Encoding) (w io.WriteCloser) {
	if encoder == nil {
		w = wIn
	} else {
		w = transform.NewWriter(wIn, encoder.NewEncoder())
	}
	return

}

func CreateFile(p string, encoder encoding.Encoding) (w io.WriteCloser, err error) {
	f, err := os.Create(p)
	if err == nil {
		w = NewEncodedWriter(f, encoder)
	}
	return
}
