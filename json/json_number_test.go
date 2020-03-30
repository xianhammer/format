package json

import (
	"strings"
	"testing"
)

func TestNumbers(t *testing.T) {
	for testID, test := range []struct {
		input       string
		expectErr   error
		expectValue func(s *testsax) bool
	}{
		// Zero values - all should return int64 = 0
		{"0", nil, isNumberInteger(0)},
		{"-0", nil, isNumberInteger(0)},
		{"0.0", nil, isNumberInteger(0)},
		{"-0.0", nil, isNumberInteger(0)},

		{"0E1", nil, isNumberInteger(0)},
		{"0E2", nil, isNumberInteger(0)},
		{"0E5", nil, isNumberInteger(0)},
		{"-0E1", nil, isNumberInteger(0)},
		{"-0E2", nil, isNumberInteger(0)},
		{"-0E5", nil, isNumberInteger(0)},
		{"0.0E1", nil, isNumberInteger(0)},
		{"0.0E2", nil, isNumberInteger(0)},
		{"0.0E5", nil, isNumberInteger(0)},
		{"-0.0E1", nil, isNumberInteger(0)},
		{"-0.0E2", nil, isNumberInteger(0)},
		{"-0.0E5", nil, isNumberInteger(0)},

		{"0E-1", nil, isNumberInteger(0)},
		{"0E-2", nil, isNumberInteger(0)},
		{"0E-5", nil, isNumberInteger(0)},
		{"-0E-1", nil, isNumberInteger(0)},
		{"-0E-2", nil, isNumberInteger(0)},
		{"-0E-5", nil, isNumberInteger(0)},
		{"0.0E-1", nil, isNumberInteger(0)},
		{"0.0E-2", nil, isNumberInteger(0)},
		{"0.0E-5", nil, isNumberInteger(0)},
		{"-0.0E-1", nil, isNumberInteger(0)},
		{"-0.0E-2", nil, isNumberInteger(0)},
		{"-0.0E-5", nil, isNumberInteger(0)},

		// Integer values
		{"1", nil, isNumberInteger(1)},
		{"21", nil, isNumberInteger(21)},
		{"321", nil, isNumberInteger(321)},

		{"1E0", nil, isNumberInteger(1)},
		{"21E0", nil, isNumberInteger(21)},
		{"321E0", nil, isNumberInteger(321)},

		{"1E-0", nil, isNumberInteger(1)},
		{"21E-0", nil, isNumberInteger(21)},
		{"321E-0", nil, isNumberInteger(321)},

		{"-1", nil, isNumberInteger(-1)},
		{"-91", nil, isNumberInteger(-91)},
		{"-981", nil, isNumberInteger(-981)},

		{"-1E0", nil, isNumberInteger(-1)},
		{"-91E0", nil, isNumberInteger(-91)},
		{"-981E0", nil, isNumberInteger(-981)},

		{"-1E-0", nil, isNumberInteger(-1)},
		{"-91E-0", nil, isNumberInteger(-91)},
		{"-981E-0", nil, isNumberInteger(-981)},

		// Float values
		{"1.0", nil, isNumberInteger(1)},
		{"21.0", nil, isNumberInteger(21)},
		{"321.0", nil, isNumberInteger(321)},
		{"-1.0", nil, isNumberInteger(-1)},
		{"-91.0", nil, isNumberInteger(-91)},
		{"-981.0", nil, isNumberInteger(-981)},

		{"0.1", nil, isNumberFloat(0.1)},
		{"1.1", nil, isNumberFloat(1.1)},
		{"21.1", nil, isNumberFloat(21.1)},
		{"321.1", nil, isNumberFloat(321.1)},
		{"-0.1", nil, isNumberFloat(-0.1)},
		{"-1.1", nil, isNumberFloat(-1.1)},
		{"-91.1", nil, isNumberFloat(-91.1)},
		{"-981.1", nil, isNumberFloat(-981.1)},

		// Float vs int
		{"1E1", nil, isNumberInteger(10)},
		{"21E1", nil, isNumberInteger(210)},
		{"321E1", nil, isNumberInteger(3210)},
		{"3E-2", nil, isNumberFloat(0.03)},
		{"21E-1", nil, isNumberFloat(2.1)},
		{"321E-1", nil, isNumberFloat(32.1)},
		//*/
	} {
		in := strings.NewReader(test.input)
		emit := new(testsax)
		err := Parse(in, emit, nil)
		if err != test.expectErr {
			t.Errorf("Test %d [%s]: Expected error %v, got %v", testID, test.input, test.expectErr, err)
		} else if !test.expectValue(emit) {
			t.Errorf("Test %d [%s]: Got unexpected type/value\n%v", testID, test.input, emit)
		}
	}
}
