package excel

import (
	"fmt"
	"strconv"

	"github.com/xianhammer/format/xml"
)

type cellXf struct {
	FontId            int
	FillId            int
	BorderId          int
	XfId              int
	ApplyNumberFormat int
	ApplyFont         int
	ApplyFill         int
	ApplyBorder       int
	QuotePrefix       int
	nf                *numFmt
	index             int
	numFmtId          string
	uniqueID          string
}

func NewCellXf(nf *numFmt) (f *cellXf) {
	f = new(cellXf)
	f.nf = nf
	return
}

func (xf *cellXf) setUniqueID() {
	var code string
	if xf.nf == nil {
		code = ""
	} else {
		code = xf.nf.Code
	}
	xf.uniqueID = fmt.Sprintf("%s_%d_%d", code, xf.ApplyNumberFormat, xf.ApplyFont) // TODO Add more...
}

func (f *cellXf) toXMLBuilder(b *xml.Builder /*, idx int*/) {
	b.Tag([]byte("xf"))

	b.Attr([]byte("numFmtId"), []byte(f.nf.numFmtId))
	b.Attr([]byte("fontId"), []byte(strconv.Itoa(f.FontId)))
	b.Attr([]byte("fillId"), []byte(strconv.Itoa(f.FillId)))
	b.Attr([]byte("borderId"), []byte(strconv.Itoa(f.BorderId)))
	b.Attr([]byte("xfId"), []byte(strconv.Itoa(f.XfId)))

	if f.QuotePrefix > 0 {
		b.Attr([]byte("quotePrefix"), []byte(strconv.Itoa(f.QuotePrefix)))
	}
	if f.ApplyNumberFormat > 0 {
		b.Attr([]byte("applyNumberFormat"), []byte(strconv.Itoa(f.ApplyNumberFormat)))
	}
	if f.ApplyFont > 0 {
		b.Attr([]byte("applyFont"), []byte(strconv.Itoa(f.ApplyFont)))
	}
	if f.ApplyFill > 0 {
		b.Attr([]byte("applyFill"), []byte(strconv.Itoa(f.ApplyFill)))
	}
	if f.ApplyBorder > 0 {
		b.Attr([]byte("applyBorder"), []byte(strconv.Itoa(f.ApplyBorder)))
	}
	b.EndTag() // End <xf>
}
