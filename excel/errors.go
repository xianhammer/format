package excel

import "errors"

var (
	ErrUnknownFile          = errors.New("Missing _rels/.rels file")
	ErrMissingRels          = errors.New("Missing _rels/.rels file")
	ErrMissingWorkbook      = errors.New("Missing workbook file")
	ErrMissingSharedstrings = errors.New("Missing sharedstrings file")
	ErrMissingStyles        = errors.New("Missing styles file")

	ErrInvalidSpans     = errors.New("Invalid spans attribute format")
	ErrInvalidDimension = errors.New("Invalid dimension attribute format")
	ErrInvalidReference = errors.New("Invalid reference attribute format")

	ErrDuplicateSheet = errors.New("Duplicate sheet")
	ErrUnknownSheet   = errors.New("Unknown sheet")

	ErrArgumentInconsictency = errors.New("Inconsistent arguments")
)
