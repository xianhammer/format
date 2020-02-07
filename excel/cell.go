package excel

import (
	"strconv"
	"time"

	"github.com/xianhammer/format/xml"
)

// NOTICE In sheet.go the appendSheet is implemented. This func copy each cell.
//        If Cell reference anything but basic values, remember to update/change
//        appendSheet accordingly!

// TODO Apply styling/formatting - somehow.

// Converting from (internal) excel time to golang time.
// Only works if dates are after 1900.
// Based on http://www.cpearson.com/excel/datetime.htm
// cpearson: "As long as all your dates later than 1900-Mar-1, this [Lotus 123 bug] should be of no concern"

var excel1900Epoc = time.Date(1899, time.December, 30, 0, 0, 0, 0, time.UTC).Unix()

type Cell struct {
	value string
	xf    *cellXf
	type_ Type
}

func (c *Cell) From(other *Cell) {
	c.value = other.value
	c.xf = other.xf
	c.type_ = other.type_
}

func (c *Cell) Value(ss *SharedStrings, applyStyle bool) (cell string) {
	if ss != nil && c.type_ == String {
		idx := 0
		for i, l := 0, len(c.value); i < l && (c.value[i]-'0') < 10; i++ {
			idx = 10*idx + int(c.value[i]&0x0f)
		}
		cell = ss.strings[idx]
	} else {
		cell = c.value
	}

	if applyStyle && c.xf != nil && c.xf.nf != nil && c.xf.nf.formatter != nil {
		cell = c.xf.nf.formatter(c.xf.nf, []byte(cell))
	}

	return
}

func (c *Cell) SetValue(ss *SharedStrings, v string) (out *Cell) {
	if c.type_ == String {
		c.value = ss.add(v)
	} else {
		c.value = v
	}
	return c
}

func (c *Cell) SetCellXf(xf *cellXf) (out *Cell) {
	c.xf = xf
	return c
}

func (c *Cell) SetType(t Type) (out *Cell) {
	c.type_ = t
	return c
}

func (c *Cell) toXMLBuilder(b *xml.Builder, s *Sheet, row, column int) {
	b.Tag([]byte("c"))
	defer b.EndTag() // End c

	b.Attr([]byte("r"), []byte(FormatDimension(row, column)))
	if c.xf != nil {
		b.Attr([]byte("s"), []byte(strconv.Itoa(c.xf.index)))
	}

	if c.value == "" {
		return
	}

	if c.type_ == String {
		b.Attr([]byte("t"), []byte("s"))
	}
	// if c.type_ == Formula {
	// 	b.Attr([]byte("t"), []byte{"str"})
	// } else {
	// 	b.Attr([]byte("t"), []byte{c._type})
	// } // TODO Set format/type...

	b.Tag([]byte("v"))
	b.Text([]byte(c.value))
	b.EndTag() // End v
}
