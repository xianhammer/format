package pdf

import (
	"bytes"
	"io"
)

// ReadHeader parse a stream for PDF header information.
func ReadHeader(r io.Reader, b []byte) (read, major, minor int, err error) {
	// Read file header and validate as PDF.
	// b := make([]byte, headersize)
	read, err = r.Read(b)
	if err != nil || err == io.EOF {
		return
	}
	return ParseHeader(b[:read])
}

// ReadFooter parse a stream for PDF footer information.
func ReadFooter(r io.Reader, b []byte) (trailer, startXRef int64, err error) {
	read, err := r.Read(b)
	if err != nil && err != io.EOF {
		return
	}

	return ParseFooter(b[:read])
}

// ReadDictionary reads a stream to parse a PDF dictionary. Keeps reading until
// a full doctionary object is read.
/*
func ReadDictionary(r io.Reader, b []byte) (read int64, err error) {
	N := 5
	reservedZone := len(b) - N
	// N bytes are reserved for buffered parsing - the maximum "keyword" size in an object.
	// This allow for simpler processing as accessing eg. b[i+1] will succeed.

	n, err := r.Read(b) // Fill buffer
	read = int64(n)
	if err != nil || err == io.EOF {
		return
	}

	buffer := new(bytes.Buffer)
	level := 1
	i, n := 0, 0

	fill := func() bool {
		copy(b[reservedZone:i], b[:])
		n, err = r.Read(b[i-reservedZone:])
		read += int64(n)
		i = 0
		return err == nil && err != io.EOF
	}

	// idxStart := 0
	var process, processDefault func()
	processName := func() {
		for ; characters[b[i]]&(c_whitespace|c_delimiter) == 0; i++ {
			if i >= reservedZone && fill() == false {
				break
			}
			buffer.WriteByte(b[i])
		}
		i--
		process = processDefault
	}
	processElement := func() {
		// if !characters[b[i]]&(c_whitespace|c_delimiter)
		// for ; characters[b[i]]&(c_whitespace|c_delimiter) == 0; i++ {
		// 	if i >= reservedZone && fill() == false {
		// 		break
		// 	}
		// 	buffer.WriteByte(b[i])
		// }
		// i--
		// process = processDefault
	}

	processDefault = func() {
		switch b[i] {
		case '<':
			if b[i+1] == '<' {
				level++
				i++
			}
		case '>':
			if b[i+1] == '>' {
				level--
				i++
			}
		case '/':
			if level == 0 {
				buffer.Reset()
				process = processName
			} // else store offset

		default:
			if level == 0 && characters[b[i]]&c_delimiter != 0 {
				process = processElement
			}
			// if characters[b[i]]&c_whitespace != 0 {
			// 	buffer.Reset()
			// }
		}
	}

	process = processDefault
	for ; level > 0; i++ {
		if i >= reservedZone && fill() == false {
			break
		}

		process()
	}

	return
}
*/
// ParseHeader parses a byte slice for the PDF file header.
// Return major and minor parts of the version number found.
func ParseHeader(b []byte) (read, major, minor int, err error) {
	l := len(b)
	if b[0] != '%' || b[1] != 'P' || b[2] != 'D' || b[3] != 'F' || b[4] != '-' {
		err = ErrInvalidFileHeader
		return
	}

	i := 5 // Offset for major version number
	for ; i < l && ('0' <= b[i] && b[i] <= '9'); i++ {
		major = 10*major + int(b[i]-'0')
	}

	if i == l || b[i] != '.' {
		err = ErrInvalidFileHeader
		return
	}

	for i++; i < l && ('0' <= b[i] && b[i] <= '9'); i++ {
		minor = 10*minor + int(b[i]-'0')
	}

	// Skip EOL
	for ; i < l && characters[b[i]]&c_eol != 0; i++ {
	}

	read = i
	return
}

/* --- EXAMPLE FOOTER ---
trailer
<</Size 8/Root 1 0 R>>
startxref
1541
%%EOF
*/

// ParseFooter parse a byte slice for PDF footer information.
// If trailer index is 0, no trailer keyword was found.
// An error is returned if either the startXRef or %%EOF are not found.
func ParseFooter(b []byte) (trailer, startXRef int64, err error) {
	l := len(b)

	idxTrailer := bytes.Index(b[:], []byte(keyword_trailer))
	if idxTrailer >= 0 {
		idxTrailer += len(keyword_trailer)
		for ; idxTrailer < l && characters[b[idxTrailer]]&c_eol != 0; idxTrailer++ { // Skip EOL characters
		}
	} else {
		idxTrailer = 0
	}
	trailer = int64(idxTrailer)

	idxStartXRef := bytes.Index(b[idxTrailer:], []byte(keyword_startxref))
	if idxStartXRef >= 0 {
		idxStartXRef += idxTrailer + len(keyword_startxref)
		for ; idxStartXRef < l && characters[b[idxStartXRef]]&c_eol != 0; idxStartXRef++ { // Skip EOL characters
		}
		for ; idxStartXRef < l && characters[b[idxStartXRef]]&c_digit != 0; idxStartXRef++ {
			startXRef = 10*startXRef + int64(b[idxStartXRef]-'0')
		}
		for ; idxStartXRef < l && characters[b[idxStartXRef]]&c_eol != 0; idxStartXRef++ { // Skip EOL characters
		}

		if bytes.Index(b[idxStartXRef:], []byte(keyword_eof)) < 0 {
			err = ErrMissingEOF
		}
	} else {
		err = ErrMissingStartXRef
	}

	return
}
