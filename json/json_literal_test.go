package json

import (
	"strings"
	"testing"
)

func TestLiterals(t *testing.T) {
	for testID, test := range []struct {
		input       string
		expectErr   error
		expectValue func(s *testsax) bool
	}{
		{"null", nil, isNull},
		{"true", nil, isTrue},
		{"false", nil, isFalse},

		{"_ull", ErrUnexpectedInput, accept},
		{"n_ll", ErrUnexpectedInput, accept},
		{"nu_l", ErrUnexpectedInput, accept},
		{"nul_", ErrUnexpectedInput, accept},
		{"null_", ErrUnexpectedInput, accept},

		{"_rue", ErrUnexpectedInput, accept},
		{"t_ue", ErrUnexpectedInput, accept},
		{"tr_e", ErrUnexpectedInput, accept},
		{"tru_", ErrUnexpectedInput, accept},
		{"true_", ErrUnexpectedInput, accept},

		{"_alse", ErrUnexpectedInput, accept},
		{"f_lse", ErrUnexpectedInput, accept},
		{"fa_se", ErrUnexpectedInput, accept},
		{"fal_e", ErrUnexpectedInput, accept},
		{"fals_", ErrUnexpectedInput, accept},
		{"false_", ErrUnexpectedInput, accept},
	} {
		in := strings.NewReader(test.input)
		emit := new(testsax)
		err := Parse(in, emit, nil)
		if err != test.expectErr {
			t.Errorf("Test %d: Expected error %v, got %v", testID, test.expectErr, err)
		} else if !test.expectValue(emit) {
			t.Errorf("Test %d: Got unexpected type/value\n%v", testID, emit)
		}
	}
}
