package excel

import (
	"bytes"
	"io"
	"testing"

	"github.com/xianhammer/format/xml"
)

const sampleSheetXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
	<dimension ref="A1:U2"/>
	<sheetViews>
		<sheetView workbookViewId="0"/>
	</sheetViews>
	<sheetFormatPr defaultRowHeight="15"/>
	<cols>
		<col customWidth="1" max="1" min="1" width="16.5703125"/>
		<col customWidth="1" max="2" min="2" width="16.7109375"/>
		<col customWidth="1" max="3" min="3" width="13.42578125"/>
		<col customWidth="1" max="4" min="4" width="10.28515625"/>
		<col customWidth="1" max="6" min="5" width="10.140625"/>
		<col customWidth="1" max="7" min="7" width="9.85546875"/>
		<col customWidth="1" max="9" min="8" width="19.7109375"/>
		<col customWidth="1" max="10" min="10" width="22.85546875"/>
		<col customWidth="1" max="11" min="11" width="28.28515625"/>
		<col customWidth="1" max="12" min="12" width="12.42578125"/>
		<col customWidth="1" max="13" min="13" width="20"/>
		<col customWidth="1" max="14" min="14" width="14.7109375"/>
		<col customWidth="1" max="15" min="15" width="10.5703125"/>
		<col customWidth="1" max="16" min="16" width="8.5703125"/>
		<col customWidth="1" max="17" min="17" width="7.42578125"/>
		<col customWidth="1" max="18" min="18" width="11.7109375"/>
		<col customWidth="1" max="19" min="19" width="12.140625"/>
		<col customWidth="1" max="20" min="20" width="15.28515625"/>
		<col customWidth="1" max="21" min="21" width="14"/>
	</cols>
	<sheetData>
		<row r="1" spans="1:21">
			<c r="A1" s="1" t="s">
				<v>4</v>
			</c>
			<c r="B1" s="1" t="s">
				<v>508</v>
			</c>
			<c r="C1" s="1" t="s">
				<v>509</v>
			</c>
			<c r="D1" s="1" t="s">
				<v>510</v>
			</c>
			<c r="E1" s="1" t="s">
				<v>511</v>
			</c>
			<c r="F1" s="1" t="s">
				<v>512</v>
			</c>
			<c r="G1" s="1" t="s">
				<v>513</v>
			</c>
			<c r="H1" s="1" t="s">
				<v>514</v>
			</c>
			<c r="I1" s="1" t="s">
				<v>515</v>
			</c>
			<c r="J1" s="1" t="s">
				<v>516</v>
			</c>
			<c r="K1" s="1" t="s">
				<v>517</v>
			</c>
			<c r="L1" s="1" t="s">
				<v>518</v>
			</c>
			<c r="M1" s="1" t="s">
				<v>519</v>
			</c>
			<c r="N1" s="1" t="s">
				<v>520</v>
			</c>
			<c r="O1" s="1" t="s">
				<v>521</v>
			</c>
			<c r="P1" s="1" t="s">
				<v>522</v>
			</c>
			<c r="Q1" s="1" t="s">
				<v>523</v>
			</c>
			<c r="R1" s="1" t="s">
				<v>524</v>
			</c>
			<c r="S1" s="1" t="s">
				<v>525</v>
			</c>
			<c r="T1" s="1" t="s">
				<v>526</v>
			</c>
			<c r="U1" s="1" t="s">
				<v>527</v>
			</c>
		</row>
		<row r="2" spans="1:21">
			<c r="A2" s="4">
				<v>12345678</v>
			</c>
			<c r="B2" s="5" t="s">
				<v>33</v>
			</c>
			<c r="C2" t="s">
				<v>33</v>
			</c>
			<c r="D2">
				<v>0</v>
			</c>
			<c r="E2">
				<v>0</v>
			</c>
			<c r="F2">
				<v>10691</v>
			</c>
			<c r="G2">
				<v>10691</v>
			</c>
			<c r="H2" t="s">
				<v>528</v>
			</c>
			<c r="I2" t="s">
				<v>529</v>
			</c>
			<c r="N2" s="5"/>
			<c r="O2" s="5"/>
			<c r="P2" s="5"/>
			<c r="R2" s="5"/>
			<c r="T2" s="5"/>
		</row>
	</sheetData>
	<autoFilter ref="A1:U2"/>
	<pageMargins bottom="0.75" footer="0.3" header="0.3" left="0.7" right="0.7" top="0.75"/>
</worksheet>`

const sampleSheetXMLSmall = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<worksheet>
	<dimension ref="B3:C3"/>
	<cols>
		<col customWidth="1" max="1" min="1" width="16.5703125"/>
		<col customWidth="1" max="2" min="2" width="16.7109375"/>
	</cols>
	<sheetData>
		<row r="1">
			<c r="A1" s="1" t="s">
				<v>4</v>
			</c>
			<c r="B1" s="1" t="s">
				<v>508</v>
			</c>
		</row>
	</sheetData>
</worksheet>`

const sampleSheetXMLSmallError = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<worksheet>
	<dimension ref="B3_C3"/>
	<cols>
		<col customWidth="1" max="1" min="1" width="16.5703125"/>
		<col customWidth="1" max="2" min="2" width="16.7109375"/>
	</cols>
	<sheetData>
		<row r="1">
			<c r="A1" s="1" t="s">
				<v>4</v>
			</c>
			<c r="B1" s="1" t="s">
				<v>508</v>
			</c>
		</row>
	</sheetData>
</worksheet>`

func sheetTest(t *testing.T, r io.Reader, expect Dimension, expectErr error) {
	sheet := new(Sheet)
	saxer := new(saxSheet)
	saxer.sheet = sheet

	tok := xml.NewTokenizer(saxer)
	_, err := tok.ReadFrom(r)

	if err != io.EOF {
		t.Errorf("Expected error [%v], got [%v]", io.EOF, err)
	}
	if saxer.Err != expectErr {
		t.Errorf("Expected sax error [%v], got [%v]", expectErr, saxer.Err)
	}

	if saxer.Err != nil {
		return
	}

	if sheet.Dimension.ColumnStart != expect.ColumnStart {
		t.Errorf("Expected dimension, cell-start [%v], got [%v]", expect.ColumnStart, sheet.Dimension.ColumnStart)
	}
	if sheet.Dimension.RowStart != expect.RowStart {
		t.Errorf("Expected dimension, row-start [%v], got [%v]", expect.RowStart, sheet.Dimension.RowStart)
	}
	if sheet.Dimension.ColumnEnd != expect.ColumnEnd {
		t.Errorf("Expected dimension, cell-end [%v], got [%v]", expect.ColumnEnd, sheet.Dimension.ColumnEnd)
	}
	if sheet.Dimension.RowEnd != expect.RowEnd {
		t.Errorf("Expected dimension, row-end [%v], got [%v]", expect.RowEnd, sheet.Dimension.RowEnd)
	}

	if sheet.Dimension.Rows() != expect.Rows() {
		t.Errorf("Expected row dimension [%v], got [%v]", expect.Rows(), sheet.Dimension.Rows())
	}
	if sheet.Dimension.Columns() != expect.Columns() {
		t.Errorf("Expected column dimension [%v], got [%v]", expect.Columns(), sheet.Dimension.Columns())
	}

	if len(sheet.Rows) != expect.Rows() {
		t.Errorf("Expected row count [%v], got [%v]", expect.Rows(), len(sheet.Rows))
	}
}

func TestSheet(t *testing.T) {
	sheetTest(t, bytes.NewBufferString(sampleSheetXML), Dimension{1, 1, 21, 2}, nil)
}

func TestSheetSmall(t *testing.T) {
	sheetTest(t, bytes.NewBufferString(sampleSheetXMLSmall), Dimension{2, 3, 3, 3}, nil)
}

func TestSheetSmallError(t *testing.T) {
	sheetTest(t, bytes.NewBufferString(sampleSheetXMLSmallError), Dimension{2, 3, 3, 3}, ErrInvalidDimension)
}
