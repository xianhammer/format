package excel

import (
	"testing"
)

func TestParseDimension(t *testing.T) {
	// <dimension ref="A1:U89"/>
	tests := []struct {
		input         string
		dims          []int
		rows, columns int
		err           error
	}{
		{"A1:U89", []int{1, 1, 21, 89}, 89, 21, nil},
		{"A1:B2", []int{1, 1, 2, 2}, 2, 2, nil},
		{"A1_B2", []int{0, 0, 0, 0}, 0, 0, ErrInvalidDimension},
		{"AZ11:BB1342", []int{52, 11, 54, 1342}, 1332, 3, nil},
	}

	for testId, expect := range tests {
		d, err := ParseDimension([]byte(expect.input))
		if err != expect.err {
			t.Errorf("Test [%d]: Expected error [%v], got [%v]", testId, expect.err, err)
		}
		if err != nil {
			continue
		}

		if d.ColumnStart != expect.dims[0] {
			t.Errorf("Test [%d]: Expected dimension, cell-start [%v], got [%v]", testId, expect.dims[0], d.ColumnStart)
		}
		if d.RowStart != expect.dims[1] {
			t.Errorf("Test [%d]: Expected dimension, row-start [%v], got [%v]", testId, expect.dims[1], d.RowStart)
		}
		if d.ColumnEnd != expect.dims[2] {
			t.Errorf("Test [%d]: Expected dimension, cell-end [%v], got [%v]", testId, expect.dims[2], d.ColumnEnd)
		}
		if d.RowEnd != expect.dims[3] {
			t.Errorf("Test [%d]: Expected dimension, row-end [%v], got [%v]", testId, expect.dims[3], d.RowEnd)
		}

		if d.Rows() != expect.rows {
			t.Errorf("Test [%d]: Expected row dimension [%v], got [%v]", testId, expect.rows, d.Rows())
		}
		if d.Columns() != expect.columns {
			t.Errorf("Test [%d]: Expected column dimension [%v], got [%v]", testId, expect.columns, d.Columns())
		}
	}
}
