package parse

import (
	"testing"
)

func TestInteger(t *testing.T) {
	// func Integer(b []byte) (i int64, n int)
	tests := []struct {
		input      string
		expect     int64
		expectSize int
	}{
		{"0", 0, 1},
		{"0a", 0, 1},
		{"0-", 0, 1},
		{"1234", 1234, 4},
		{"1234--", 1234, 4},
	}

	for testID, test := range tests {
		got, gotSize := Integer([]byte(test.input))
		if got != test.expect {
			t.Errorf("[test=%d] Expected [%v], got [%v]\n", testID, test.expect, got)
		}
		if gotSize != test.expectSize {
			t.Errorf("[test=%d] Expected [%v], got [%v]\n", testID, test.expectSize, gotSize)
		}
	}
}

func TestFloat(t *testing.T) {
	// func Integer(b []byte) (i int64, n int)
	tests := []struct {
		input      string
		expect     float64
		expectSize int
	}{
		{"0", 0, 1},
		{"0a", 0, 1},
		{"0-", 0, 1},
		{"0.1", 0.1, 3},
		{"9.1", 9.1, 3},
		{"9e1", 90, 3},
		{"9.5e1", 95, 5},
		{"9.5e4", 95000, 5},
		{"314159e-5", 3.1415900000000003, 9},
	}

	for testID, test := range tests {
		got, gotSize := Float([]byte(test.input))
		if got != test.expect {
			t.Errorf("[test=%d] Expected [%v], got [%v]\n", testID, test.expect, got)
		}
		if gotSize != test.expectSize {
			t.Errorf("[test=%d] Expected [%v], got [%v]\n", testID, test.expectSize, gotSize)
		}
	}
}
