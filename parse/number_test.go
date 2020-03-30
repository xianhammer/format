package parse

import (
	"testing"
)

func TestHex(t *testing.T) {
	// func Hex(b []byte) (i uint64, n int)
	tests := []struct {
		input      string
		expect     uint64
		expectSize int
	}{
		{"0", 0, 1},
		{"0a", 0x0A, 2},
		{"0g", 0, 1},
		{"0-", 0, 1},
		{"1234", 0x1234, 4},
		{"1234--", 0x1234, 4},

		{"g012", 0, 0},
		{"f012", 0xF012, 4},
		{"f012g", 0xF012, 4},
		{"F012g", 0xF012, 4},
		{"A0g", 0xA0, 4},

		{"3DA408B9", 0x3DA408B9, 8},
		{"FFFFFFFF5B", 0xFFFFFFFF5B, 10},
	}

	for testID, test := range tests {
		got, gotSize := Hex([]byte(test.input))
		if got != test.expect {
			t.Errorf("[test=%d] Expected [%v], got [%v]\n", testID, test.expect, got)
		}
		if gotSize != test.expectSize {
			t.Errorf("[test=%d] Expected [%v], got [%v]\n", testID, test.expectSize, gotSize)
		}
	}
}

func TestInteger(t *testing.T) {
	// func Integer(b []byte) (i uint64, n int)
	tests := []struct {
		input      string
		expect     uint64
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
	// Float(b []byte) (f float64, n int)
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
