package json

import (
	"testing"
)

func TestUtil_unmarshal(t *testing.T) {
	tests := []struct {
		input      string
		expect     interface{}
		expectSize int
	}{
		{"0", float64(0), 1}, // JSON Numbers are float64
		{"1", float64(1), 1},
		{"-1", float64(-1), 2},
		{"123", float64(123), 3},
		{"123a", float64(123), 3},

		{"123.0", float64(123.0), 5},
		{"123.33", float64(123.33), 6},
		{"123.66", float64(123.66), 6},
		{"-123.66", float64(-123.66), 7},

		{"123e-2", float64(1.23), 6},
		{"123e+2", float64(12300), 6},
		{"123e2", float64(12300), 5},

		{"-123e-2", float64(-1.23), 7},
		{"-123e+2", float64(-12300), 7},
		{"-123e2", float64(-12300), 6},
	}

	for testID, test := range tests {
		got, gotSize, err := Parse([]byte(test.input), nil)
		if err != nil {
			t.Errorf("[test=%d] Expected error [%v], got [%v]\n", testID, nil, err)
		}
		if !Equal(got, test.expect) {
			t.Errorf("[test=%d] Expected [%v], got [%v]\n", testID, test.expect, got)
		}
		if gotSize != test.expectSize {
			t.Errorf("[test=%d] Expected [%v], got [%v]\n", testID, test.expectSize, gotSize)
		}
	}
}
