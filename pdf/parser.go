package pdf

import (
	"fmt"
	"io"
)

// Parse a file as PDF returning, on success, a PDF Document.
func Parse(r io.ReadSeeker) (doc Document, err error) {
	// Reading the file header is not relevant, unless the ReadFileFooter is a success.
	// However, if the content is streamed, it IS relevant to check if the file(header) is valid.

	// readHeader, major, minor, err := ReadFileHeader(r)
	// _, _, _, err = ReadFileHeader(r)
	// if err != nil {
	// 	return
	// }

	// doc, err = New(major, minor)
	// if err != nil {
	// 	return
	// }

	// "fileHeader" bytes are actually read, so seek back to immediately after header.
	// Currently, this is not relevant, as the filefooter is the next read.
	// r.Seek(readHeader, io.SeekBegin)

	// The ReadFileFooter does NOT encompas the Seek method. This is due to PDF
	// may contain several footer structures!
	headersize, footersize := fileHeader, fileTrailer
	filesize, _ := r.Seek(0, io.SeekEnd)
	if filesize < footersize {
		footersize = filesize
	}
	if filesize < headersize {
		headersize = filesize
	}

	// Create reusable buffer.
	var buffer []byte
	if headersize > footersize {
		buffer = make([]byte, headersize)
	} else {
		buffer = make([]byte, footersize)
	}

	// Process the file footer
	r.Seek(-footersize, io.SeekEnd)
	/*idxTrailer*/ _, startxref, err := ReadFooter(r, buffer[:footersize])
	if err != nil && err != io.EOF {
		return
	}

	// Process the file header
	r.Seek(0, io.SeekStart)
	_, major, minor, err := ReadHeader(r, buffer[:headersize])
	if err != nil && err != io.EOF {
		return
	}
	doc.Version.Set(major, minor)

	// Process the trailer dictionary
	// r.Seek(idxTrailer, io.SeekStart)
	// readTrailer, trailer, err := ReadDictionary(r, buffer[:])
	// fmt.Printf("readTrailer, trailer = %d, %v\n", readTrailer, trailer)

	// Process the start XRef
	doc.StartXRef = NewXRef(startxref)
	r.Seek(startxref, io.SeekStart)
	n, err := doc.StartXRef.Read(r)
	fmt.Printf("0: n, err  = %v, %v\n", n, err)
	if err != nil {
		return
	}
	r.Seek(startxref+n, io.SeekStart)

	trailer := NewObject()
	n, err = trailer.Read(r)
	fmt.Printf("1: n, err  = %v, %v\n", n, err)
	fmt.Printf("   trailer = %v\n", trailer)

	// if err != nil && err != io.EOF {
	// 	return
	// }

	return
}
