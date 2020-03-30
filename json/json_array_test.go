package json

import (
	"strings"
	"testing"
)

func TestArray(t *testing.T) {
	for testID, test := range []struct {
		input     string
		expectErr error
	}{
		{"[]", nil},
		{"[1]", nil},
		{"[1,2,3,4]", nil},
	} {
		in := strings.NewReader(test.input)
		emit := new(testsax)
		err := Parse(in, emit, nil)
		if err != test.expectErr {
			t.Errorf("Test %d: Expected error %v, got %v", testID, test.expectErr, err)
			// } else if !test.expectValue(emit) {
			// 	t.Errorf("Test %d: Got unexpected type/value\n%v", testID, emit)
		}
	}
}
