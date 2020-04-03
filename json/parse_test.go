package json

import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestParse_number(t *testing.T) {
	// func Parse(b []byte) (i uint64, n int)
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

func TestParse_simple(t *testing.T) {
	// func Parse(b []byte) (i uint64, n int)
	tests := []struct {
		input      string
		expect     interface{}
		expectSize int
		expectErr  error
	}{
		{"true", true, 4, nil},
		{"false", false, 5, nil},
		{"null", nil, 4, nil},

		{"trueZ", true, 4, nil},
		{"falseZ", false, 5, nil},
		{"nullZ", nil, 4, nil},
		{"tru", nil, 3, io.EOF},
		{"fals", nil, 4, io.EOF},
		{"nul", nil, 3, io.EOF},

		{`"some text"`, "some text", 2 + 9, nil},
		{`"some\ntext"`, "some\ntext", 2 + 9 + 1, nil},
		{`"some\ttext"`, "some\ttext", 2 + 9 + 1, nil},
		{`"some\"text"`, "some\"text", 2 + 9 + 1, nil},
		{`"some text"z`, "some text", 2 + 9, nil},
		{`"some text"z`, "some text", 2 + 9, nil},

		{`"<a href=\"http://www.apache.org/\">"`, `<a href="http://www.apache.org/">`, 2 + 35, nil},
		{`"<a href=\"http://www.apache.org/\"><img src=\"https://www.apache.org/images/asf_logo_wide.gif\">"`, `<a href="http://www.apache.org/"><img src="https://www.apache.org/images/asf_logo_wide.gif">`, 2 + 96, nil},
	}

	for testID, test := range tests {
		got, gotSize, err := Parse([]byte(test.input), nil)
		if err != test.expectErr {
			t.Errorf("[test=%d] Expected error [%v], got [%v]\n", testID, test.expectErr, err)
		}

		if !Equal(got, test.expect) {
			t.Errorf("[test=%d] Expected [%v], got [%v]\n", testID, test.expect, got)
		}
		if gotSize != test.expectSize {
			t.Errorf("[test=%d] Expected [%v], got [%v]\n", testID, test.expectSize, gotSize)
		}
	}
}

var complex = `{

  "nodeName" : "",
  "description" : "<a href=\"http://www.apache.org/\"><img src=\"https://www.apache.org/images/asf_logo_wide.gif\"></img></a>\r\n<p>\r\nThis is a public build and test server for <a href=\"http://projects.apache.org/\">projects</a> of the\r\n<a href=\"http://www.apache.org/\">Apache Software Foundation</a>. All times on this server are UTC.\r\n</p>\r\n<p>\r\nSee the <a href=\"http://wiki.apache.org/general/Hudson\">Jenkins wiki page</a> for more information\r\nabout this service.\r\n</p>"
}`
var complexExpect = map[string]interface{}{
	"nodeName":    "",
	"description": "<a href=\"http://www.apache.org/\"><img src=\"https://www.apache.org/images/asf_logo_wide.gif\"></img></a>\r\n<p>\r\nThis is a public build and test server for <a href=\"http://projects.apache.org/\">projects</a> of the\r\n<a href=\"http://www.apache.org/\">Apache Software Foundation</a>. All times on this server are UTC.\r\n</p>\r\n<p>\r\nSee the <a href=\"http://wiki.apache.org/general/Hudson\">Jenkins wiki page</a> for more information\r\nabout this service.\r\n</p>",
}

func TestParse_compound(t *testing.T) {
	// func JSON(b []byte) (i uint64, n int)
	tests := []struct {
		input      string
		expect     interface{}
		expectSize int
	}{
		{"[]", []interface{}{}, 2},
		{"[ ]", []interface{}{}, 3},
		{"[1]", []interface{}{1.0}, 3},
		{"[1,2]", []interface{}{1.0, 2.0}, 5},
		{"[null,null]", []interface{}{nil, nil}, 11},
		{"[true,1 ,  \"string\"]", []interface{}{true, 1.0, "string"}, 20},

		{"[true,[1 ,  \"string\"]]", []interface{}{true, []interface{}{1.0, "string"}}, 22},
		{"[1,2,3]\n", []interface{}{1.0, 2.0, 3.0}, 7},
		{"[1,2,3]\r", []interface{}{1.0, 2.0, 3.0}, 7},

		{"{}", map[string]interface{}{}, 2},
		{"{ }", map[string]interface{}{}, 3},
		{"{\t}", map[string]interface{}{}, 3},
		{"{\n}", map[string]interface{}{}, 3},
		{`{"a":1}`, map[string]interface{}{"a": 1.0}, 7},
		{`{"a":1,"b":2}`, map[string]interface{}{"a": 1.0, "b": 2.0}, 13},
		{`{"a":1,"b":2}` + "\n", map[string]interface{}{"a": 1.0, "b": 2.0}, 13},
		{`{"a":1,"b":2}` + "\r", map[string]interface{}{"a": 1.0, "b": 2.0}, 13},

		{`{"a":-1}`, map[string]interface{}{"a": -1.0}, 8},
		{`{"a":-1,"b":2}`, map[string]interface{}{"a": -1.0, "b": 2.0}, 14},
		{`{"a":1,"b":-2}`, map[string]interface{}{"a": 1.0, "b": -2.0}, 14},

		{`{"key1":"a text\nsecond line","other key":[1,2,3]}`,
			map[string]interface{}{"key1": "a text\nsecond line", "other key": []interface{}{1.0, 2.0, 3.0}},
			50},

		{"[{}]", []interface{}{map[string]interface{}{}}, 4},
		{"[{ }]", []interface{}{map[string]interface{}{}}, 5},
		{"[{\n}]", []interface{}{map[string]interface{}{}}, 5},

		{"[{},[],{}]", []interface{}{map[string]interface{}{}, []interface{}{}, map[string]interface{}{}}, 10},

		{"{\n\"d\"\n:\n[\n{\n}\n]\n}\n", map[string]interface{}{
			"d": []interface{}{map[string]interface{}{}},
		}, 17},

		{complex, complexExpect, 517},
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

func TestParse_large(t *testing.T) {
	tests := []struct {
		filename  string
		expectErr error
	}{
		{"apache_builds", nil},
		{"canada", nil},
		{"citm_catalog", nil},
		{"github_events", nil},
		{"gsoc-2018", nil},
		{"instruments", nil},
		{"marine_ik", nil},
		{"mesh", nil},
		{"mesh.pretty", nil},
		{"numbers", nil},
		{"random", nil},
		{"twitter", nil},
		{"twitterescaped", nil},
		{"update-center", nil},
	}

	for testID, test := range tests {
		file := filepath.Join("testdata", "json0", test.filename+"-utf8.json")
		src, err := ioutil.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}

		srcLen := len(bytes.TrimSpace(src))
		got, gotSize, err := Parse(src, nil)
		if err != nil && err != io.EOF {
			t.Errorf("[test=%d] Expected error [%v], got [%v]\n", testID, nil, err)
		}
		if got == nil {
			t.Errorf("[test=%d] Expected non-nil, got [%v]\n", testID, got)
		}
		if gotSize != srcLen {
			t.Errorf("[test=%d] Expected [%v], got [%v]\n", testID, srcLen, gotSize)
		}
	}
}
