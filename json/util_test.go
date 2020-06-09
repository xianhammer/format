package json

import (
	"testing"
)

func TestUtil_unmarshal(t *testing.T) {
	tests := []struct {
		input     string
		expect    interface{}
		expectErr error
	}{
		{`0`, new(int), nil},
		{`true`, new(bool), nil},
		{`"abc"`, new(string), nil},

		{`[0]`, &[]int{}, nil},
		{`[true]`, &[]bool{}, nil},
		{`["abc"]`, &[]string{}, nil},

		{`[0,-1,+2,1e1,2e+2,3e-3]`, &[]float32{}, nil},
		{`[0,-1,+2,1e1,2e+2,3e-3]`, &[]float64{}, nil},

		{`[0,-1,+2,1e1,2e+2]`, &[]int{}, nil},
		{`[0,-1,+2,1e1,2e+2]`, &[]int8{}, nil},
		{`[0,-1,+2,1e1,2e+2]`, &[]int16{}, nil},
		{`[0,-1,+2,1e1,2e+2]`, &[]int32{}, nil},
		{`[0,-1,+2,1e1,2e+2]`, &[]int64{}, nil},

		{`[0,-1,+2,1e1,2e+2]`, &[]uint{}, nil},
		{`[0,-1,+2,1e1,2e+2]`, &[]uint8{}, nil},
		{`[0,-1,+2,1e1,2e+2]`, &[]uint16{}, nil},
		{`[0,-1,+2,1e1,2e+2]`, &[]uint32{}, nil},
		{`[0,-1,+2,1e1,2e+2]`, &[]uint64{}, nil},

		{`{"a":"hello"}`, &struct {
			a string
		}{}, nil},

		{`{"a":"hello","b":{"c":-2,"d":[9,3.141]}}`, &struct {
			a string
			b struct {
				c int
				d []float32
			}
		}{}, nil},

		{`{"a":"hello","e":true,"b":{"c":-2,"d":[9,3.141]}}`, &struct {
			a string
			b struct {
				c int
				d []float32
			}
			e bool
		}{}, nil},

		{`0`, new(string), ErrInvalidType},
		{`true`, new(string), ErrInvalidType},
		{`"abc"`, new(int), ErrInvalidType},

		{`[0]`, []int{}, ErrMustBePointer},
		{`[true]`, &[]int{}, ErrInvalidType},

		{`[0,"abc"]`, &[]int{}, ErrInvalidType},
		{`[0,"abc"]`, &[]int8{}, ErrInvalidType},
		{`[0,"abc"]`, &[]int16{}, ErrInvalidType},
		{`[0,"abc"]`, &[]int32{}, ErrInvalidType},
		{`[0,"abc"]`, &[]int64{}, ErrInvalidType},

		{`[0,"abc"]`, &[]uint{}, ErrInvalidType},
		{`[0,"abc"]`, &[]uint8{}, ErrInvalidType},
		{`[0,"abc"]`, &[]uint16{}, ErrInvalidType},
		{`[0,"abc"]`, &[]uint32{}, ErrInvalidType},
		{`[0,"abc"]`, &[]uint64{}, ErrInvalidType},

		{`[0,"abc"]`, &[]float32{}, ErrInvalidType},
		{`[0,"abc"]`, &[]float64{}, ErrInvalidType},

		{`[0,"abc"]`, &[]string{}, ErrInvalidType},
		{`[0,"abc"]`, &[]bool{}, ErrInvalidType},

		{`[0,0]`, &[]complex64{}, ErrInvalidType},

		{`{"a":"hello"}`, &struct {
			a bool
		}{}, ErrInvalidType},
	}

	for testID, test := range tests {
		got, _, err := Parse([]byte(test.input), nil)
		if err != nil {
			t.Errorf("[test=%d] Expected error [%v], got [%v]\n", testID, nil, err)
		}

		err = Unmarshal(got, test.expect)
		if err != test.expectErr {
			t.Errorf("[test=%d] Expected error [%v], got [%v]\n", testID, test.expectErr, err)
		}
	}
}

func TestUtil_equal(t *testing.T) {
	tests := []struct {
		input     string
		expect    interface{}
		expectErr error
	}{

		{"true", float64(-1), nil},
		{"-1", true, nil},

		{"[1]", []interface{}{interface{}(1), 2}, nil},
		{"[1]", []interface{}{interface{}(2)}, nil},
		{"[1,2]", []interface{}{interface{}(1)}, nil},

		{"{\"a\":1}", map[string]interface{}{"a": interface{}(1), "b": interface{}(2)}, nil},
		{"{\"a\":1}", map[string]interface{}{"a": interface{}("1")}, nil},
		{"{\"a\":1,\"b\":2}", map[string]interface{}{"a": interface{}(1)}, nil},
	}

	for testID, test := range tests {
		got, _, err := Parse([]byte(test.input), nil)

		if Equal(got, test.expect) {
			t.Errorf("[test=%d] Expected [%v], got [%v]\n", testID, test.expect, got)
		}
		if err != test.expectErr {
			t.Errorf("[test=%d] Expected [%v], got [%v]\n", testID, test.expectErr, err)
		}
	}
}
