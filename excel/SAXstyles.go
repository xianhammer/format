package excel

import (
	"bytes"
	"strconv"

	"github.com/xianhammer/format/xml"
)

type saxStyles struct {
	xml.Partial
	inNumFmts          bool
	currentNumFmt      *numFmt
	inCellXfs          bool
	currentCellXf      *cellXf
	inCellStyleXfs     bool
	currentCellStyleXf *cellStyleXf
	target             *Styles
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
	} else if s.inCellStyleXfs {
		if bytes.Equal(name, []byte("xf")) {
			s.currentCellStyleXf = NewCellStyleXf(nil)
			s.target.AddCellStyleXf(s.currentCellStyleXf)
		} else {
			s.inCellStyleXfs = !bytes.Equal(name, []byte("/cellStyleXfs"))
		}
	} else if s.inNumFmts {
		if bytes.Equal(name, []byte("numFmt")) {
			s.currentNumFmt = new(numFmt)
		} else {
			s.inNumFmts = !bytes.Equal(name, []byte("/numFmts"))
		}
	} else if bytes.Equal(name, []byte("cellXfs")) {
		s.inCellXfs = true
	} else if bytes.Equal(name, []byte("cellStyleXfs")) {
		s.inCellStyleXfs = true
	} else if bytes.Equal(name, []byte("numFmts")) {
		s.inNumFmts = true
	}
}

// <numFmt numFmtId="8" formatCode="#,##0.00\ "kr.";[Red]\-#,##0.00\ "kr.""/>

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
	} else if s.inCellStyleXfs {
		if bytes.Equal(name, []byte("numFmtId")) {
			s.currentCellStyleXf.numFmtId = string(value)
		} else if bytes.Equal(name, []byte("fontId")) {
			s.currentCellStyleXf.FontId = toInt(value)
		} else if bytes.Equal(name, []byte("fillId")) {
			s.currentCellStyleXf.FillId = toInt(value)
		} else if bytes.Equal(name, []byte("borderId")) {
			s.currentCellStyleXf.BorderId = toInt(value)
		} else if bytes.Equal(name, []byte("xfId")) {
			s.currentCellStyleXf.XfId = toInt(value)
		} else if bytes.Equal(name, []byte("applyNumberFormat")) {
			s.currentCellStyleXf.ApplyNumberFormat = toInt(value)
		} else if bytes.Equal(name, []byte("applyFont")) {
			s.currentCellStyleXf.ApplyFont = toInt(value)
		} else if bytes.Equal(name, []byte("applyFill")) {
			s.currentCellStyleXf.ApplyFill = toInt(value)
		} else if bytes.Equal(name, []byte("applyBorder")) {
			s.currentCellStyleXf.ApplyBorder = toInt(value)
		} else if bytes.Equal(name, []byte("applyProtection")) {
			s.currentCellStyleXf.ApplyProtection = toInt(value)
		} else if bytes.Equal(name, []byte("quotePrefix")) {
			s.currentCellStyleXf.QuotePrefix = toInt(value)
		}
	}
}
