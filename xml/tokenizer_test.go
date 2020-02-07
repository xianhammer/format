package xml

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/xianhammer/format/parse"
)

const DEBUG = false

type token int

const (
	Tag token = iota
	TagEnd
	Text
	Comment
	Section
	Attribute
)

func (t token) String() string {
	switch t {
	case Tag:
		return "Tag"
	case TagEnd:
		return "TagEnd"
	case Text:
		return "Text"
	case Comment:
		return "Comment"
	case Section:
		return "Section"
	case Attribute:
		return "Attribute"
	default:
		return "<error>"
	}
}

type element struct {
	input     string
	expectErr error
	expect    []token
	expectStr []string
	t         *testing.T
	idx       int
	idxStr    int
	testIdx   int
}

func (e *element) Tag(name []byte) {
	if DEBUG {
		fmt.Println("Tag")
	}
	e.check(Tag)
	e.checkString(name)
}
func (e *element) TagEnd(autoclose bool) {
	if DEBUG {
		fmt.Println("TagEnd")
	}
	e.check(TagEnd)
}
func (e *element) Text(value []byte) {
	if DEBUG {
		fmt.Println("Text")
	}
	e.check(Text)
	e.checkString(value)
}
func (e *element) Comment(value []byte) {
	if DEBUG {
		fmt.Println("Comment")
	}
	e.check(Comment)
	e.checkString(value)
}
func (e *element) Section(name []byte) {
	if DEBUG {
		fmt.Println("Section")
	}
	e.check(Section)
	e.checkString(name)
}
func (e *element) Attribute(tag, name, value []byte) {
	if DEBUG {
		fmt.Println("Attribute")
	}
	e.check(Attribute)
	e.checkString(name)
	e.checkString(value)
}

func (e *element) check(got token) {
	if e.expect[e.idx] != got {
		e.t.Errorf("Test %d: Token error, expected [%s], got [%s]\n", e.testIdx, e.expect[e.idx], got)
	}
	e.idx++
}

func (e *element) checkString(got []byte) {
	if e.expectStr[e.idxStr] != string(got) {
		e.t.Errorf("Test %d: Identifier/Value error, expected [%s], got [%s]\n", e.testIdx, e.expectStr[e.idxStr], got)
	}
	e.idxStr++
}

func (e *element) run(testIdx int) {
	e.testIdx = testIdx

	e.t.Logf("Test %d: %s\n", testIdx, e.input)

	br := bytes.NewBufferString(e.input)
	t := NewTokenizer(e)
	n, err := t.ReadFrom(br)

	if err != e.expectErr {
		e.t.Errorf("Test %d @%d: Token error, expected [%v], got [%v]\n", e.testIdx, n, e.expectErr, err)
	}
	if err == io.EOF && n != int64(len(e.input)) {
		e.t.Errorf("Test %d @%d: Input read error, expected [%d] bytes, read [%v]\n", e.testIdx, n, len(e.input), n)
	}
	return
}

func TestTreeStructure(t *testing.T) {
	tests := []*element{
		// Basic tags
		{"<tag>", io.EOF, []token{Tag, TagEnd}, []string{"tag"}, t, 0, 0, 0},
		{"</tag>", io.EOF, []token{Tag, TagEnd}, []string{"/tag"}, t, 0, 0, 0},
		{"<tag/>", io.EOF, []token{Tag, TagEnd}, []string{"tag"}, t, 0, 0, 0},
		{"<tag >", io.EOF, []token{Tag, TagEnd}, []string{"tag"}, t, 0, 0, 0},
		{"</tag >", io.EOF, []token{Tag, TagEnd}, []string{"/tag"}, t, 0, 0, 0},
		{"<tag />", io.EOF, []token{Tag, TagEnd}, []string{"tag"}, t, 0, 0, 0},
		{"<tag></tag>", io.EOF, []token{Tag, TagEnd, Tag, TagEnd}, []string{"tag", "/tag"}, t, 0, 0, 0},
		{"<tag></other>", io.EOF, []token{Tag, TagEnd, Tag, TagEnd}, []string{"tag", "/other"}, t, 0, 0, 0},

		// Tags and text
		{"<tag>abc</tag>", io.EOF, []token{Tag, TagEnd, Text, Tag, TagEnd}, []string{"tag", "abc", "/tag"}, t, 0, 0, 0},
		{"12<tag>abc</tag>56", io.EOF, []token{Text, Tag, TagEnd, Text, Tag, TagEnd, Text}, []string{"12", "tag", "abc", "/tag", "56"}, t, 0, 0, 0},

		// Open and close tags, namespace
		{"<ns:tag>", io.EOF, []token{Tag, TagEnd}, []string{"ns:tag"}, t, 0, 0, 0},
		{"</ns:tag>", io.EOF, []token{Tag, TagEnd}, []string{"/ns:tag"}, t, 0, 0, 0},
		{"<ns:tag/>", io.EOF, []token{Tag, TagEnd}, []string{"ns:tag"}, t, 0, 0, 0},
		{"<ns:tag></ns:tag>", io.EOF, []token{Tag, TagEnd, Tag, TagEnd}, []string{"ns:tag", "/ns:tag"}, t, 0, 0, 0},

		// Attributes
		{"<tag a=\"1\">", io.EOF, []token{Tag, Attribute, TagEnd}, []string{"tag", "a", "1"}, t, 0, 0, 0},
		{"<tag a=\"2\" />", io.EOF, []token{Tag, Attribute, TagEnd}, []string{"tag", "a", "2"}, t, 0, 0, 0},
		{"<tag a=\"2\" success=\"true\"/>", io.EOF, []token{Tag, Attribute, Attribute, TagEnd}, []string{"tag", "a", "2", "success", "true"}, t, 0, 0, 0},

		// Special tags, see:
		// https://xmlwriter.net/xml_guide/cdata_section.shtml
		// https://xmlwriter.net/xml_guide/conditional_section.shtml
		// https://xmlwriter.net/xml_guide/comment.shtml
		{`<?tag?>`, io.EOF, []token{Tag, TagEnd}, []string{"?tag"}, t, 0, 0, 0},
		{`<?xml version="1.0" encoding="utf-8"?>`, io.EOF, []token{Tag, Attribute, Attribute, TagEnd}, []string{"?xml", "version", "1.0", "encoding", "utf-8"}, t, 0, 0, 0},
		{`<!--comment content-->`, io.EOF, []token{Comment, TagEnd}, []string{"comment content"}, t, 0, 0, 0},
		{`<![INCLUDE[cdata-block]]>`, io.EOF, []token{Section, Text, TagEnd}, []string{"INCLUDE", "cdata-block"}, t, 0, 0, 0},
		{`<![IGNORE[cdata-block]]>`, io.EOF, []token{Section, Text, TagEnd}, []string{"IGNORE", "cdata-block"}, t, 0, 0, 0},
		{`<![CDATA[cdata-block]]>`, io.EOF, []token{Section, Text, TagEnd}, []string{"CDATA", "cdata-block"}, t, 0, 0, 0},
		{`<![CDATA[cdata-block]]><tag>`, io.EOF, []token{Section, Text, TagEnd, Tag, TagEnd}, []string{"CDATA", "cdata-block", "tag"}, t, 0, 0, 0},
		{`<![CDATA[cdata-block]]><tag/>`, io.EOF, []token{Section, Text, TagEnd, Tag, TagEnd}, []string{"CDATA", "cdata-block", "tag"}, t, 0, 0, 0},
		{`<!DOCTYPE cafProductFeed SYSTEM "http://www.affiliatewindow.com/DTD/affiliate/datafeed.1.5.dtd">`,
			io.EOF,
			[]token{Section, Text, TagEnd},
			[]string{"DOCTYPE", `cafProductFeed SYSTEM "http://www.affiliatewindow.com/DTD/affiliate/datafeed.1.5.dtd"`},
			t, 0, 0, 0,
		},

		// More realistic test
		{`<tag attr1="value1" attr2="2" hidden="1">'MixedChars'!$A$1:$AE$10754</tag>`, io.EOF,
			[]token{Tag, Attribute, Attribute, Attribute, TagEnd, Text, Tag, TagEnd},
			[]string{"tag", "attr1", "value1", "attr2", "2", "hidden", "1", "'MixedChars'!$A$1:$AE$10754", "/tag"}, t, 0, 0, 0},

		// Syntax errors
		{"<ns:123 />", ErrBadIdentifier, nil, nil, t, 0, 0, 0},
		{"<ns::tag />", ErrBadIdentifier, nil, nil, t, 0, 0, 0},
		{"</ns::tag>", ErrBadIdentifier, nil, nil, t, 0, 0, 0},
		{"<tag a=1>", ErrIllegalCharacter, []token{Tag}, []string{"tag"}, t, 0, 0, 0},
		{"<tag =\"1\">", ErrIllegalCharacter, []token{Tag}, []string{"tag"}, t, 0, 0, 0},
		{"<tag =1>", ErrIllegalCharacter, []token{Tag}, []string{"tag"}, t, 0, 0, 0},
		{"<tag \"1\">", ErrIllegalCharacter, []token{Tag}, []string{"tag"}, t, 0, 0, 0},
		{"<tag 1>", ErrIllegalCharacter, []token{Tag}, []string{"tag"}, t, 0, 0, 0},
		{"<tag", ErrUnterminatedTag, nil, nil, t, 0, 0, 0},
		{"</tag", ErrUnterminatedTag, nil, nil, t, 0, 0, 0},
		{"<tag/", ErrUnterminatedTag, []token{Tag}, []string{"tag"}, t, 0, 0, 0},
		{"<tag /", ErrUnterminatedTag, []token{Tag}, []string{"tag"}, t, 0, 0, 0},
		{"<ns:>", ErrUnterminatedIdentifier, []token{Tag, TagEnd}, []string{"ns:"}, t, 0, 0, 0},
		{"</ns:>", ErrUnterminatedIdentifier, []token{Tag, TagEnd}, []string{"/ns:"}, t, 0, 0, 0},
		{"<ns:/>", ErrUnterminatedIdentifier, []token{Tag, TagEnd}, []string{"ns:"}, t, 0, 0, 0},
		{"<tag a=\">", ErrUnterminatedString, []token{Tag}, []string{"tag"}, t, 0, 0, 0},
		{"<tag a=\"value>", ErrUnterminatedString, []token{Tag}, []string{"tag"}, t, 0, 0, 0},
		{"<tag a=\"value/>", ErrUnterminatedString, []token{Tag}, []string{"tag"}, t, 0, 0, 0},
		{"<tag a=\"value  />", ErrUnterminatedString, []token{Tag}, []string{"tag"}, t, 0, 0, 0},

		{"<tag <tag", ErrUnterminatedTag, []token{Tag}, []string{"tag"}, t, 0, 0, 0},
		{"<tag a=<tag", ErrUnterminatedTag, []token{Tag}, []string{"tag"}, t, 0, 0, 0},
		{"<tag a=\"1\"<tag", ErrUnterminatedTag, []token{Tag, Attribute}, []string{"tag", "a", "1"}, t, 0, 0, 0},
		{`< --comment content-->`, ErrIllegalCharacter, nil, nil, t, 0, 0, 0},
		{`< -- comment content-->`, ErrIllegalCharacter, nil, nil, t, 0, 0, 0},
		{`<- - comment content-->`, ErrIllegalCharacter, nil, nil, t, 0, 0, 0},
		{`<!- - comment content-->`, ErrIllegalCharacter, nil, nil, t, 0, 0, 0},
		{`<! -- comment content-->`, ErrIllegalCharacter, nil, nil, t, 0, 0, 0},
		{`< ![CDATA[cdata-block]]>`, ErrIllegalCharacter, nil, nil, t, 0, 0, 0},
	}

	for i, e := range tests {
		e.run(i)
	}
}

func TestBufferBounds(t *testing.T) {
	e := &element{"<tagname>", io.EOF, []token{Tag, TagEnd}, []string{"tagname"}, t, 0, 0, 0}

	br := bytes.NewBufferString(e.input)

	tok := NewTokenizer(e)
	tok.TagName = parse.NewBuffer(3)

	streamSize, err := tok.ReadFrom(br)

	if err != parse.ErrOutOfBounds {
		t.Errorf("[pos=%d] Expected error [%v], got [%v]\n", streamSize, parse.ErrOutOfBounds, err)
	}
}

// func BenchmarkFib10(b *testing.B) {
//         // run the Fib function b.N times
//         for n := 0; n < b.N; n++ {
//                 Fib(10)
//         }
// }
