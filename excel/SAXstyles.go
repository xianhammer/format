package excel

import (
	"bytes"
	"strconv"

	"github.com/xianhammer/format/xml"
)

type saxStyles struct {
	xml.Partial
	inNumFmts     bool
	currentNumFmt *numFmt
	inCellXfs     bool
	currentCellXf *cellXf
	target        *Styles
}

func toInt(b []byte) (i int) {
	i, _ = strconv.Atoi(string(b))
	return
}

func (s *saxStyles) Tag(name []byte) {
	if s.inCellXfs {
		if bytes.Equal(name, []byte("xf")) {
			s.currentCellXf = NewCellXf(nil)
			s.target.AddCellXf(s.currentCellXf)
		} else {
			s.inCellXfs = !bytes.Equal(name, []byte("/cellXfs"))
		}
	} else if s.inNumFmts {
		if bytes.Equal(name, []byte("numFmt")) {
			s.currentNumFmt = new(numFmt)
		} else {
			s.inNumFmts = !bytes.Equal(name, []byte("/numFmts"))
		}
	} else if bytes.Equal(name, []byte("cellXfs")) {
		s.inCellXfs = true
	} else if bytes.Equal(name, []byte("numFmts")) {
		s.inNumFmts = true
	}
}

func (s *saxStyles) Attribute(tag, name, value []byte) {
	if s.inNumFmts {
		if bytes.Equal(name, []byte("numFmtId")) {
			s.currentNumFmt.numFmtId = string(value)
			s.target.AddNumFmt(s.currentNumFmt)
		} else if bytes.Equal(name, []byte("formatCode")) {
			s.currentNumFmt.SetCode(string(value))
		}
	} else if s.inCellXfs {
		if bytes.Equal(name, []byte("numFmtId")) {
			s.currentCellXf.numFmtId = string(value)
		} else if bytes.Equal(name, []byte("fontId")) {
			s.currentCellXf.FontId = toInt(value)
		} else if bytes.Equal(name, []byte("fillId")) {
			s.currentCellXf.FillId = toInt(value)
		} else if bytes.Equal(name, []byte("borderId")) {
			s.currentCellXf.BorderId = toInt(value)
		} else if bytes.Equal(name, []byte("xfId")) {
			s.currentCellXf.XfId = toInt(value)
		} else if bytes.Equal(name, []byte("applyNumberFormat")) {
			s.currentCellXf.ApplyNumberFormat = toInt(value)
		} else if bytes.Equal(name, []byte("applyFont")) {
			s.currentCellXf.ApplyFont = toInt(value)
		} else if bytes.Equal(name, []byte("applyFill")) {
			s.currentCellXf.ApplyFill = toInt(value)
		} else if bytes.Equal(name, []byte("applyBorder")) {
			s.currentCellXf.ApplyBorder = toInt(value)
		} else if bytes.Equal(name, []byte("quotePrefix")) {
			s.currentCellXf.QuotePrefix = toInt(value)
		}
	}
}
