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

		{`"<a href=\"http://www.apache.org/\">"`, `<a href="http://www.apache.org/">`, 2 + 35},
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
		{"[ ]", []interface{}{}, 3},
		{"[1]", []interface{}{1.0}, 3},
		{"[1,2]", []interface{}{1.0, 2.0}, 5},
		{"[null,null]", []interface{}{nil, nil}, 11},
		{"[true,1 ,  \"string\"]", []interface{}{true, 1.0, "string"}, 20},
		{"[true,[1 ,  \"string\"]]", []interface{}{true, []interface{}{1.0, "string"}}, 22},

		{"{}", map[string]interface{}{}, 2},
		{"{ }", map[string]interface{}{}, 3},
		{"{\t}", map[string]interface{}{}, 3},
		{"{\n}", map[string]interface{}{}, 3},
		{`{"a":1}`, map[string]interface{}{"a": 1.0}, 7},
		{`{"a":1,"b":2}`, map[string]interface{}{"a": 1.0, "b": 2.0}, 13},
		{`{"key1":"a text\nsecond line","other key":[1,2,3]}`,
			map[string]interface{}{"key1": "a text\nsecond line", "other key": []interface{}{1.0, 2.0, 3.0}},
			50},

		{"[{}]", []interface{}{map[string]interface{}{}}, 4},
		{"[{ }]", []interface{}{map[string]interface{}{}}, 5},
		{"[{\n}]", []interface{}{map[string]interface{}{}}, 5},

		{"[{},[],{}]", []interface{}{map[string]interface{}{}, []interface{}{}, map[string]interface{}{}}, 10},
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

const large = `{
  "assignedLabels" : [
    {
      
    }
  ],
  "mode" : "EXCLUSIVE",
  "nodeDescription" : "the master Jenkins node",
  "nodeName" : "",
  "numExecutors" : 0,
  "description" : "<a href=\"http://www.apache.org/\"><img src=\"https://www.apache.org/images/asf_logo_wide.gif\"></img></a>\r\n<p>\r\nThis is a public build and test server for <a href=\"http://projects.apache.org/\">projects</a> of the\r\n<a href=\"http://www.apache.org/\">Apache Software Foundation</a>. All times on this server are UTC.\r\n</p>\r\n<p>\r\nSee the <a href=\"http://wiki.apache.org/general/Hudson\">Jenkins wiki page</a> for more information\r\nabout this service.\r\n</p>"

  "overallLoad" : {
    
  },
  "primaryView" : {
    "name" : "All",
    "url" : "https://builds.apache.org/"
  },
  "quietingDown" : false,
  "slaveAgentPort" : 0,
  "unlabeledLoad" : {
    
  },
  "useCrumbs" : true,
  "useSecurity" : true,
  "views" : [
    {
      "name" : "All",
      "url" : "https://builds.apache.org/"
    },
    {
      "name" : "CloudStack",
      "url" : "https://builds.apache.org/view/CloudStack/"
    },
    {
      "name" : "Hadoop",
      "url" : "https://builds.apache.org/view/Hadoop/"
    },
    {
      "name" : "Onami",
      "url" : "https://builds.apache.org/view/Onami/"
    }
  ]
}`

func TestJSON_large(t *testing.T) {
	// Test if the input is parsed and check length...
	got, gotSize := JSON([]byte(large), nil)
	if got == nil {
		t.Errorf("[test=%d] Expected non-nil, got [%v]\n", 0, got)
	}
	if gotSize != len(large) {
		t.Errorf("[test=%d] Expected [%v], got [%v]\n", 0, len(large), gotSize)
	}
}
