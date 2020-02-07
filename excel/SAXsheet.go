package excel

import (
	"bytes"

	"github.com/xianhammer/format/xml"
)

type attributeType int

const (
	ignore    attributeType = 0
	cell                    = 1
	dimension               = 2
)

type saxSheet struct {
	xml.Partial
	Err         error
	sheet       *Sheet
	cellCount   int
	cellIndex   int
	cellXf      *cellXf
	cellType    byte
	row         []Cell
	acquireText bool
	attribute   attributeType
}

func (s *saxSheet) Tag(name []byte) {
	if s.Err != nil {
		return
	}

	switch name[0] {
	case '/': // Close tag
		if bytes.Equal(name, []byte("/row")) {
			s.row = nil
		}
		s.attribute = ignore

	case 'v':
		s.acquireText = len(name) == 1

	case 'c':
		if len(name) == 1 {
			s.attribute = cell
		}

	case 'r':
		if bytes.Equal(name, []byte("row")) {
			s.row = make([]Cell, s.cellCount)
			s.sheet.Rows = append(s.sheet.Rows, s.row)
		}

	case 'd':
		if bytes.Equal(name, []byte("dimension")) {
			s.attribute = dimension
		}
	}
}

func (s *saxSheet) Attribute(tag, name, value []byte) {
	switch s.attribute {
	case ignore:

	case cell:
		if len(name) != 1 {
			return
		}

		switch name[0] {
		case 'r': // Reference
			s.cellIndex = 0
			for i, l := 0, len(value); i < l && (value[i]-'A') < azBase; i++ {
				s.cellIndex = azBase*s.cellIndex + int(value[i]&0x1f)
			}
			s.cellIndex-- // Make it zero-offset.

			// Ignore the rowID - expected to be the same as the one found in "row" attributes!
		case 't': // Type
			if s.cellType = value[0]; bytes.Equal(value, []byte("str")) {
				s.cellType = 'f' // Special flag for formula (str) cell type
			}

		case 's': // Style
			var style int
			for i, l := 0, len(value); i < l && (value[i]-'0') < 10; i++ {
				style = 10*style + int(value[i]&0x0f)
			}
			s.cellXf = s.sheet.workbook.styles.GetCellXf(style)
		}

	case dimension:
		if bytes.Equal(name, []byte("ref")) {
			s.sheet.Dimension, s.Err = ParseDimension(value)
			s.cellCount = s.sheet.Dimension.Columns()
		}
	}
}

func (s *saxSheet) Text(value []byte) {
	if s.acquireText {
		s.acquireText = false
		s.row[s.cellIndex].value = string(value)
		s.row[s.cellIndex].xf = s.cellXf
		s.row[s.cellIndex].type_ = Type(s.cellType)
		s.cellXf = nil
		s.cellType = 0
	}
}
