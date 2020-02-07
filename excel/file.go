package excel

import (
	"archive/zip"
	"io"
)

type File zip.File

func (f *File) Open() (io.ReadCloser, error) {
	return (*zip.File)(f).Open()
}

func (f *File) QuerySelectorAll(tagname string) (elements []Attributes, err error) {
	r, err := f.Open()
	if err != nil {
		return
	}

	return NodeAttributes(r, []byte(tagname))
}
