package excel

import (
	"bufio"
	"fmt"
	"io"

	"github.com/xianhammer/format/xml"
)

// Sheet of an Excel workbook, providing access to Cells and Rows.
type Sheet struct {
	Attributes Attributes
	Dimension  Dimension
	Rows       [][]Cell

	size     int64
	err      error
	workbook *Workbook
	file     *File
}

func newSheet(workbook *Workbook, name string) (s *Sheet) {
	s = new(Sheet)
	s.workbook = workbook
	s.Attributes = make(Attributes)
	s.Attributes["name"] = name

	s.Dimension.RowStart = 1
	s.Dimension.RowEnd = 1
	s.Dimension.ColumnStart = 1
	s.Dimension.ColumnEnd = 1

	return
}

func (s *Sheet) appendSheet(src *Sheet, addCols []ColumnUpdate) {
	styleTgt := s.workbook.styles
	styleSrc := src.workbook.styles

	styleTgt.merge(styleSrc)

	offset := len(s.Rows)
	s.Rows = append(s.Rows, make([][]Cell, len(src.Rows))...)
	appendRows := s.Rows[offset:]

	extend := len(addCols)

	sharedStringsSrc := src.workbook.sharedstrings
	sharedStringsTgt := s.workbook.sharedstrings

	formatMap := make(map[string]*cellXf)
	for _, xf := range styleTgt.cellXfs {
		formatMap[xf.nf.Code] = xf
	}

	for row := range src.Rows {
		srcRow := src.Rows[row]

		appendRows[row] = make([]Cell, len(srcRow)+extend)

		tgt := appendRows[row]
		for col := range srcRow {
			tgt[col].type_ = srcRow[col].type_

			if xf := srcRow[col].xf; xf != nil {
				tgt[col].xf = formatMap[xf.nf.Code] // This should never fail as all formats from src are already merged into target.
			}

			value := srcRow[col].Value(sharedStringsSrc, false)
			tgt[col].SetValue(sharedStringsTgt, value)
		}

		tgt = tgt[len(srcRow):]
		for col := range addCols {
			tgt[col].xf = styleTgt.cellXfs[0]
			tgt[col].type_ = String
			newValue := addCols[col](row, col, &tgt[col])
			tgt[col].SetValue(sharedStringsTgt, newValue)
		}
	}

	s.refresh()
}

func (s *Sheet) open() (n int64, err error) {
	if s.file == nil {
		return s.size, s.err
	}

	defer func() {
		s.size = n
		s.err = err
		s.file = nil
	}()

	r, err := s.file.Open()
	if err != nil {
		return
	}

	return s.ReadFrom(r)
}

func (s *Sheet) refresh() {
	s.Dimension.RowEnd = (s.Dimension.RowStart - 1) + len(s.Rows)
	// fmt.Printf("s.Dimension.RowEnd=%d (s.Dimension.RowStart=%d)\n", s.Dimension.RowEnd, s.Dimension.RowStart)
	if len(s.Rows) > 0 {
		// TODO Need to validate the -1
		s.Dimension.ColumnEnd = (s.Dimension.ColumnStart + len(s.Rows[0]) - 1)
	}
	// s.Dimension.RowEnd = s.Dimension.RowStart + len(w.Rows)
}

type ColumnUpdate func(row, column int, c *Cell) (newValue string)

// AddColumns add new column(s) to the current sheet.
func (s *Sheet) AddColumns(cols []ColumnUpdate) {
	sharedStrings := s.workbook.sharedstrings
	styleTgt := s.workbook.styles

	extend := len(cols)
	for row := range s.Rows {
		offset := len(s.Rows[row])
		s.Rows[row] = append(s.Rows[row], make([]Cell, extend)...)

		target := s.Rows[row][offset:]
		for i := 0; i < extend; i++ {
			target[i].xf = styleTgt.cellXfs[0]
			target[i].type_ = String
			newValue := cols[i](row, offset+i, &target[i])
			target[i].SetValue(sharedStrings, newValue)
		}
	}
	s.refresh()
}

// SharedStrings return a reference to the sharedstrings(.xml)
func (s *Sheet) SharedStrings() (ss *SharedStrings) {
	return s.workbook.sharedstrings
}

// Styles return a reference to the style(.xml)
func (s *Sheet) Styles() (ss *Styles) {
	return s.workbook.styles
}

// Row return indexed row.
func (s *Sheet) Row(row int) (r Row) {
	return s.Rows[row]
}

// Cell return specified cell.
func (s *Sheet) Cell(row, col int) (c *Cell) {
	return &s.Rows[row][col]
}

// CellByRef return cell by Excel reference, like "C3"
func (s *Sheet) CellByRef(ref string) (cell *Cell, err error) {
	if col, row, err := ParseReference([]byte(ref)); err == nil {
		cell = s.Cell(row-1, col-1)
	}
	return
}

// ColumnByTitle
func (s *Sheet) ColumnByTitle(row int, title string) (column int, err error) {
	if 0 < row && row < len(s.Rows) {
		err = fmt.Errorf("ColumnByTitle: Row index %d is out of range [%d:%d]", row, 0, len(s.Rows))
		return
	}

	column = -1
	sharedStrings := s.workbook.sharedstrings
	for i, c := range s.Rows[row] {
		cellValue := c.Value(sharedStrings, false)
		if cellValue == title {
			column = i
			break
		}
	}

	return
}

// ReadFrom implement the io.ReaderFrom interface.
func (s *Sheet) ReadFrom(r io.Reader) (n int64, err error) {
	saxer := new(saxSheet)
	saxer.sheet = s

	t := xml.NewTokenizer(saxer)

	n, err = t.ReadFrom(bufio.NewReader(r))
	if saxer.Err != nil && (err == nil || err == io.EOF) {
		err = saxer.Err
	}

	return
}

// WriteTo implement the io.WriterTo interface.
// TODO Values below are fixed on purpose (In current version at least).
func (s *Sheet) WriteTo(w io.Writer) (n int64, err error) {
	b := xml.NewBuilder(w)
	defer func() {
		err = b.Close()
	}()

	b.Tag([]byte("worksheet"))
	b.Attr([]byte("xmlns"), []byte(Spreadsheet))
	b.Attr([]byte("xmlns:r"), []byte(RelationshipsDoc))
	b.Attr([]byte("xmlns:mc"), []byte(MarkupCompatibility))
	b.Attr([]byte("xmlns:x14ac"), []byte(X14ac))
	b.Attr([]byte("mc:Ignorable"), []byte("x14ac"))

	defer b.EndTag() // End worksheet

	// <dimension ref="A1:U1"/>
	sqRef := s.Dimension.String()
	b.Tag([]byte("dimension"))
	b.Attr([]byte("ref"), []byte(sqRef))
	b.EndTag() // End dimension

	b.Tag([]byte("sheetViews"))

	b.Tag([]byte("sheetView"))
	b.Attr([]byte("tabSelected"), []byte("1"))
	b.Attr([]byte("workbookViewId"), []byte("0"))
	b.EndTag() // End sheetView

	b.EndTag() // End sheetViews

	b.Tag([]byte("sheetFormatPr"))
	b.Attr([]byte("defaultRowHeight"), []byte("15")) // Fixed value, on purpose (in current version)
	b.EndTag()

	// <sheetData>
	b.Tag([]byte("sheetData"))

	rowIndex := s.Dimension.RowStart
	columnIndex := s.Dimension.ColumnStart
	span := fmt.Sprintf("%d:%d", s.Dimension.ColumnStart, s.Dimension.ColumnEnd)
	for y := range s.Rows {
		// 	<row r="1" s="1" spans="1:NN">
		b.Tag([]byte("row"))
		b.Attr([]byte("r"), []byte(FormatDimension(rowIndex+y, 0)))
		b.Attr([]byte("spans"), []byte(span))

		row := s.Rows[y][:]
		for x := range row {
			row[x].toXMLBuilder(b, s, rowIndex+y, columnIndex+x)
		}

		b.EndTag() // End row
	}

	b.EndTag() // End sheetData

	// <autoFilter ref="A1:U1"/>
	b.Tag([]byte("autoFilter"))
	b.Attr([]byte("ref"), []byte(sqRef))
	b.EndTag() // End dimension

	b.Tag([]byte("pageMargins"))
	b.Attr([]byte("left"), []byte("0.7"))
	b.Attr([]byte("right"), []byte("0.7"))
	b.Attr([]byte("top"), []byte("0.75"))
	b.Attr([]byte("bottom"), []byte("0.75"))
	b.Attr([]byte("header"), []byte("0.3"))
	b.Attr([]byte("footer"), []byte("0.3"))
	b.EndTag()

	b.Tag([]byte("pageSetup"))
	b.Attr([]byte("paperSize"), []byte("9"))
	b.Attr([]byte("orientation"), []byte("portrait"))
	b.Attr([]byte("r:id"), []byte("rId1")) // TODO # 1 should probably be sheet id or something...
	b.EndTag()

	return
}

func (s *Sheet) AppendEmptyRow(width int) (r []Cell) {
	r = make([]Cell, width)
	s.Rows = append(s.Rows, r)

	for i := 0; i < width; i++ {
		// r[i].SetCellXf(NewCellXf(FormatDefault))
		r[i].SetType(String)
	}

	s.refresh()
	return
}

func (s *Sheet) AppendRow() (r *newrow) {
	r = new(newrow)
	r.sheet = s
	return
}

type newrow struct {
	sheet   *Sheet
	cells   []*Cell
	indeces []int
	maxidx  int
}

func (r *newrow) Close() (err error) {
	cells := make([]Cell, r.maxidx+1)
	r.sheet.Rows = append(r.sheet.Rows, cells)

	r.sheet.refresh()
	if len(cells) > r.sheet.Dimension.ColumnEnd {
		r.sheet.Dimension.ColumnEnd = len(cells)
	}

	for i, c := range r.cells {
		cells[r.indeces[i]] = *c
	}

	return
}

func (r *newrow) Set(idx int, c *Cell) {
	r.cells = append(r.cells, c)
	r.indeces = append(r.indeces, idx)
	if idx > r.maxidx {
		r.maxidx = idx
	}

	if c.type_ == String {
		c.value = r.sheet.workbook.sharedstrings.add(c.value)
	}

	return
}
