package excel

import (
	"fmt"
	"strconv"

	"github.com/xianhammer/format/xml"
)

type cellStyleXf struct {
	FontId            int
	FillId            int
	BorderId          int
	XfId              int
	ApplyNumberFormat int
	ApplyFont         int
	ApplyFill         int
	ApplyBorder       int
	ApplyProtection   int
	QuotePrefix       int
	nf                *numFmt
	index             int
	numFmtId          string
	uniqueID          string
}

func NewCellStyleXf(nf *numFmt) (f *cellStyleXf) {
	f = new(cellStyleXf)
	f.nf = nf
	return
}

func (xf *cellStyleXf) setUniqueID() {
	var code string
	if xf.nf == nil {
		code = ""
	} else {
		code = xf.nf.Code
	}
	xf.uniqueID = fmt.Sprintf("%s_%d_%d", code, xf.ApplyNumberFormat, xf.ApplyFont) // TODO Add more...
}

func (xf *cellStyleXf) toXMLBuilder(b *xml.Builder /*, idx int*/) {
	b.Tag([]byte("xf"))

	b.Attr([]byte("numFmtId"), []byte(xf.nf.numFmtId))
	b.Attr([]byte("fontId"), []byte(strconv.Itoa(xf.FontId)))
	b.Attr([]byte("fillId"), []byte(strconv.Itoa(xf.FillId)))
	b.Attr([]byte("borderId"), []byte(strconv.Itoa(xf.BorderId)))
	b.Attr([]byte("xfId"), []byte(strconv.Itoa(xf.XfId)))

	if xf.QuotePrefix > 0 {
		b.Attr([]byte("quotePrefix"), []byte(strconv.Itoa(xf.QuotePrefix)))
	}
	if xf.ApplyNumberFormat > 0 {
		b.Attr([]byte("applyNumberFormat"), []byte(strconv.Itoa(xf.ApplyNumberFormat)))
	}
	if xf.ApplyFont > 0 {
		b.Attr([]byte("applyFont"), []byte(strconv.Itoa(xf.ApplyFont)))
	}
	if xf.ApplyFill > 0 {
		b.Attr([]byte("applyFill"), []byte(strconv.Itoa(xf.ApplyFill)))
	}
	if xf.ApplyBorder > 0 {
		b.Attr([]byte("applyBorder"), []byte(strconv.Itoa(xf.ApplyBorder)))
	}
	if xf.ApplyProtection > 0 {
		b.Attr([]byte("applyProtection"), []byte(strconv.Itoa(xf.ApplyProtection)))
	}
	b.EndTag() // End <xf>
}
