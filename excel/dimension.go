package excel

import (
	"bytes"
	"fmt"
	"strconv"
)

var oftenFound = []byte("A1:")

type Dimension struct {
	ColumnStart, RowStart, ColumnEnd, RowEnd int
}

func (d Dimension) Columns() int {
	return d.ColumnEnd - d.ColumnStart + 1
}

func (d Dimension) Rows() int {
	return d.RowEnd - d.RowStart + 1
}

func (d Dimension) Start() (s string) {
	return FormatDimension(d.RowStart, d.ColumnStart)
}

func (d Dimension) End() (s string) {
	// return FormatDimension(d.RowEnd-1, d.ColumnEnd)
	return FormatDimension(d.RowEnd, d.ColumnEnd)
}

func (d Dimension) String() (s string) {
	return d.Start() + ":" + d.End()
}

func (d Dimension) workbookReference() (s string) {
	value := d.ColumnStart - 1
	msb, lsb := value/azBase, value%azBase
	if msb != 0 {
		s = fmt.Sprintf("$%c%c$%d:", 'A'+(msb-1), 'A'+lsb, d.RowStart)
	} else {
		s = fmt.Sprintf("$%c$%d:", 'A'+lsb, d.RowStart)
	}

	value = d.ColumnEnd - 1
	msb, lsb = value/azBase, value%azBase
	if msb != 0 {
		s += fmt.Sprintf("$%c%c$%d", 'A'+(msb-1), 'A'+lsb, d.RowEnd)
	} else {
		s += fmt.Sprintf("$%c$%d", 'A'+lsb, d.RowEnd)
	}

	return
}

func ParseDimension(value []byte) (d Dimension, err error) {
	i := 0
	if bytes.Equal(value, oftenFound) {
		i = 3
		d.ColumnStart = 1
		d.RowStart = 1
	} else {
		for l := len(value); i < l && (value[i]-'A') < azBase; i++ {
			d.ColumnStart = azBase*d.ColumnStart + int(value[i]&0x1f)
		}

		for l := len(value); i < l && (value[i]-'0') < 10; i++ {
			d.RowStart = 10*d.RowStart + int(value[i]&0x0f)
		}

		if value[i] == ':' {
			i++
		} else {
			err = ErrInvalidDimension
		}
	}

	for l := len(value); i < l && (value[i]-'A') < azBase; i++ {
		d.ColumnEnd = azBase*d.ColumnEnd + int(value[i]&0x1f)
	}

	for l := len(value); i < l && (value[i]-'0') < 10; i++ {
		d.RowEnd = 10*d.RowEnd + int(value[i]&0x0f)
	}
	return
}

func FormatDimension(row, column int) (s string) {
	if column > 0 {
		value := column - 1
		msb, lsb := value/azBase, value%azBase
		if msb != 0 {
			s = fmt.Sprintf("%c%c", 'A'+(msb-1), 'A'+lsb)
		} else {
			s = fmt.Sprintf("%c", 'A'+lsb)
		}
	}

	if row > 0 {
		s += strconv.Itoa(row)
	}

	return
}
