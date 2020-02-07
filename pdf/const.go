package pdf

import (
	"errors"
)

// keywords
const (
	keyword_trailer   string = "trailer"
	keyword_startxref        = "startxref"
	keyword_eof              = "%%EOF"
)

var (
	ErrUnexectedMarker     = errors.New("Unexpected marker")
	ErrUnexpectedToken     = errors.New("Unexpected token")
	ErrUnexpectedCharacter = errors.New("Unexpected character")
	ErrNumberTooLong       = errors.New("Number too long")
	ErrBufferTooSmall      = errors.New("Buffer is too small")
	ErrInvalidFileHeader   = errors.New("Invalid file header")
	ErrInvalidFileTrailer  = errors.New("Invalid file trailer")
	ErrInvalidVersion      = errors.New("Invalid file version")
	ErrInvalidXRef         = errors.New("Invalid cross-reference table")
	ErrInvalidXRefEntry    = errors.New("Invalid cross-reference entry")
	ErrInvalidStartXRef    = errors.New("Invalid " + keyword_startxref)
	ErrInvalidDictionary   = errors.New("Not a dictionary")
	ErrMissingTrailer      = errors.New("Missing " + keyword_trailer)
	ErrMissingStartXRef    = errors.New("Missing " + keyword_startxref)
	ErrMissingEOF          = errors.New("Missing " + keyword_eof)
	ErrMissingObjectStart  = errors.New("Missing <<")
	// ErrInvalidXRefSection = errors.New("Invalid cross-reference section")

	// rxXrefFirstLine = regexp.MustCompile(`\n([0-9]+) ([0-9]+)\n`)
	// rxFileHeader    = regexp.MustCompile(`%PDF-([0-9]+)\.([0-9]+)(\r\n|\n\n|\n)`)

	CurrentDefaultVersion = Version{1, 7}
)

const (
	// See Appendix C (Implementation Limits) in PDF reference 1.7
	maxSizeNameObject = 127
	maxSizeNumber     = 32 // Not specified in doc, but set to a sufficiently high value to represent real.
	defaultReadAhead  = 80 // App. line length
	fileHeader        = int64(20)
	fileTrailer       = int64(1024)
)

// https://www.prepressure.com/pdf/basics/version
const (
	MimePdf    string = "application/pdf"
	MimeXPdf   string = "application/x-pdf"
	MimeXBZPdf string = "application/x-bzpdf"
	MimeXGZPdf string = "application/x-gzpdf"
)
