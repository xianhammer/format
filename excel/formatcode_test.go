package excel

import (
	"testing"
)

func TestFormatcode(t *testing.T) {
	tests := []struct {
		input  string
		output format
	}{
		{`dd\-mm\-yyyy\ hh:mm:ss`, format{"02-01-2006 15:04:05", 't'}},
		{`dd\-mm\-yyyy`, format{"02-01-2006", 't'}},
		{`dd-mm-yyyy hh:mm:ss`, format{"02-01-2006 15:04:05", 't'}},
		{`dd-mm-yyyy`, format{"02-01-2006", 't'}},
		{`dd-MM-yyyy hh:mm:ss`, format{"02-01-2006 15:04:05", 't'}},
		{`dd-MM-yyyy`, format{"02-01-2006", 't'}},
		{`d-M-yy`, format{"2-1-06", 't'}},
		{`dd-MM-yyyy Z`, format{"02-01-2006 Z07", 't'}},
		{`dd-MM-yyyy ZZ`, format{"02-01-2006 Z0700", 't'}},
		{`dd-MM-yyyy Z:Z`, format{"02-01-2006 Z07:00", 't'}},
		{`dd-MM-yyyy TZName`, format{"02-01-2006 MST", 't'}},
		{`##\ ##\ ##\ ##`, format{"## ## ## ##", 't'}},
		{`## ## ## ##`, format{"## ## ## ##", 't'}},
	}

	for testId, expect := range tests {
		got := Formatcode(testId, expect.input)
		if expect.output.format != got.format {
			t.Errorf("Test[%d]: Expected output [%v], got [%v]", testId, expect.output.format, got.format)
		}
		if testId != got.key {
			t.Errorf("Test[%d]: Expected output [%v], got [%v]", testId, testId, got.key)
		}
	}
}
