package csv

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"
)

const DEBUG = false

type token int

type callback struct {
	input     string
	expectErr error
	expectStr []string
	t         *testing.T
	idxStr    int
	testIdx   int
}

func (e *callback) Field(column, line uint, value []byte) {
	if DEBUG {
		fmt.Println("Field")
	}
	e.checkString(value)
}

func (e *callback) checkString(got []byte) {
	if e.expectStr[e.idxStr] != string(got) {
		e.t.Errorf("Test %d: Identifier/Value error, expected [%s], got [%s]\n", e.testIdx, e.expectStr[e.idxStr], got)
	}
	e.idxStr++
}

func (e *callback) run(testIdx, inputLen int, r io.Reader) {
	e.testIdx = testIdx

	e.t.Logf("Test %d: %s\n", testIdx, e.input)

	// br := bytes.NewBufferString(e.input)
	t := NewTokenizer(e)
	t.LineComment = '#'

	n, err := t.ReadFrom(r)

	if err != e.expectErr {
		e.t.Errorf("Test %d @%d: Token error, expected [%v], got [%v]\n", e.testIdx, n, e.expectErr, err)
	}
	if err == io.EOF && n != int64(inputLen) {
		e.t.Errorf("Test %d @%d: Input read error, expected [%d] bytes, read [%v]\n", e.testIdx, n, inputLen, n)
	}
	return
}

type limitedReader struct {
	r     io.Reader
	limit int
	err   error
}

func (r *limitedReader) Read(b []byte) (n int, err error) {
	if r.limit <= 0 {
		return 0, r.err
	}

	if len(b) > r.limit {
		n, err = r.r.Read(b[:r.limit])
	} else {
		n, err = r.r.Read(b)
	}

	r.limit -= n
	return
}

func Test1(t *testing.T) {
	tests := []*callback{
		// Basic tags
		{`1,2,3`, io.EOF, []string{"1", "2", "3"}, t, 0, 0},
		{`"1",2,3`, io.EOF, []string{"1", "2", "3"}, t, 0, 0},
		{`1,"2",3`, io.EOF, []string{"1", "2", "3"}, t, 0, 0},
		{`1,2,"3"`, io.EOF, []string{"1", "2", "3"}, t, 0, 0},
		{`"1",2,"3"`, io.EOF, []string{"1", "2", "3"}, t, 0, 0},
		{"1,2,3\n4,5,6", io.EOF, []string{"1", "2", "3", "4", "5", "6"}, t, 0, 0},
		{"a1,b2,3c\na4,b5,6c", io.EOF, []string{"a1", "b2", "3c", "a4", "b5", "6c"}, t, 0, 0},
		{"1,22,333,4444,55555", io.EOF, []string{"1", "22", "333", "4444", "55555"}, t, 0, 0},
		{"#1,22,333,4444,55555", io.EOF, []string{}, t, 0, 0},
		{"#1,22,333,4444,55555\nabc,def,ghi", io.EOF, []string{"abc", "def", "ghi"}, t, 0, 0},
		{`unquoted,"simple quote",unquoted with ""quotes"","ab""cd"""`, io.EOF, []string{"unquoted", "simple quote", "unquoted with \"quotes\"", "ab\"cd\""}, t, 0, 0},
	}

	for i, e := range tests {
		e.run(i, len(e.input), bytes.NewBufferString(e.input))
	}
}

func TestError(t *testing.T) {
	expectErr := errors.New("dummy")
	tests := []*callback{
		// Basic tags
		{`1,2,3`, expectErr, []string{"1", "2", "3"}, t, 0, 0},
		{`"1",2,3`, expectErr, []string{"1", "2", "3"}, t, 0, 0},
		{`1,"2",3`, expectErr, []string{"1", "2", "3"}, t, 0, 0},
		{`1,2,"3"`, expectErr, []string{"1", "2", "3"}, t, 0, 0},
		{`"1",2,"3"`, expectErr, []string{"1", "2", "3"}, t, 0, 0},
		{"1,2,3\n4,5,6", expectErr, []string{"1", "2", "3", "4", "5", "6"}, t, 0, 0},
		{"a1,b2,3c\na4,b5,6c", expectErr, []string{"a1", "b2", "3c", "a4", "b5", "6c"}, t, 0, 0},
		{"1,22,333,4444,55555", expectErr, []string{"1", "22", "333", "4444", "55555"}, t, 0, 0},
		{"#1,22,333,4444,55555", expectErr, []string{}, t, 0, 0},
		{"#1,22,333,4444,55555\nabc,def,ghi", expectErr, []string{"abc", "def", "ghi"}, t, 0, 0},
	}

	for i, e := range tests {
		lr := &limitedReader{
			bytes.NewBufferString(e.input),
			len(e.input) - 3,
			expectErr,
		}
		e.run(i, len(e.input), lr)
	}
}
