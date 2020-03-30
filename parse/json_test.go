package parse

import (
	"testing"
)

func TestJSON_number(t *testing.T) {
	// func JSON(b []byte) (i uint64, n int)
	tests := []struct {
		input      string
		expect     interface{}
		expectSize int
	}{
		{"0", float64(0), 1}, // JSON Numbers are float64
		{"1", float64(1), 1},
		{"123", float64(123), 3},
		{"123a", float64(123), 3},

		{"123.0", float64(123.0), 5},
		{"123.33", float64(123.33), 6},
		{"123.66", float64(123.66), 6},

		{"123e-2", float64(1.23), 6},
		{"123e+2", float64(12300), 6},
		{"123e2", float64(12300), 5},
	}

	for testID, test := range tests {
		got, gotSize := JSON([]byte(test.input), nil)
		if !JSONEqual(got, test.expect) {
			t.Errorf("[test=%d] Expected [%v], got [%v]\n", testID, test.expect, got)
		}
		if gotSize != test.expectSize {
			t.Errorf("[test=%d] Expected [%v], got [%v]\n", testID, test.expectSize, gotSize)
		}
	}
}

func TestJSON_simple(t *testing.T) {
	// func JSON(b []byte) (i uint64, n int)
	tests := []struct {
		input      string
		expect     interface{}
		expectSize int
	}{
		{"true", true, 4},
		{"false", false, 5},
		{"null", nil, 4},
		{`"some text"`, "some text", 2 + 9},
		{`"some\ntext"`, "some\ntext", 2 + 9 + 1},
		{`"some\ttext"`, "some\ttext", 2 + 9 + 1},
		{`"some\"text"`, "some\"text", 2 + 9 + 1},

		{"trueZ", true, 4},
		{"falseZ", false, 5},
		{"nullZ", nil, 4},
		{`"some text"z`, "some text", 2 + 9},
	}

	for testID, test := range tests {
		got, gotSize := JSON([]byte(test.input), nil)
		if !JSONEqual(got, test.expect) {
			t.Errorf("[test=%d] Expected [%v], got [%v]\n", testID, test.expect, got)
		}
		if gotSize != test.expectSize {
			t.Errorf("[test=%d] Expected [%v], got [%v]\n", testID, test.expectSize, gotSize)
		}
	}
}

func TestJSON_compound(t *testing.T) {
	// func JSON(b []byte) (i uint64, n int)
	tests := []struct {
		input      string
		expect     interface{}
		expectSize int
	}{
		{"[]", []interface{}{}, 2},
		{"[1]", []interface{}{1.0}, 3},
		{"[1,2]", []interface{}{1.0, 2.0}, 5},
		{"[true,1 ,  \"string\"]", []interface{}{true, 1.0, "string"}, 20},
		{"[true,[1 ,  \"string\"]]", []interface{}{true, []interface{}{1.0, "string"}}, 22},

		{"{}", map[string]interface{}{}, 2},
		{`{"a":1}`, map[string]interface{}{"a": 1.0}, 7},
		{`{"a":1,"b":2}`, map[string]interface{}{"a": 1.0, "b": 2.0}, 13},
		{`{"key1":"a text\nsecond line","other key":[1,2,3]}`,
			map[string]interface{}{"key1": "a text\nsecond line", "other key": []interface{}{1.0, 2.0, 3.0}},
			50},
	}

	for testID, test := range tests {
		got, gotSize := JSON([]byte(test.input), nil)
		if !JSONEqual(got, test.expect) {
			t.Errorf("[test=%d] Expected [%v], got [%v]\n", testID, test.expect, got)
		}
		if gotSize != test.expectSize {
			t.Errorf("[test=%d] Expected [%v], got [%v]\n", testID, test.expectSize, gotSize)
		}
	}
}
