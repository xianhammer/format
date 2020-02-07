package xml

import (
	"io"
	"strings"
	"testing"
)

func newbuilder(t *testing.T, sb *strings.Builder) (b *Builder) {
	b = NewBuilder(sb)
	if b == nil {
		t.Fatal("NewBuilder returned nil")
	}

	if sb.String() != string(DefaultHeader) {
		t.Fatalf("NewBuilder, expected '%s' got '%s'", DefaultHeader, sb.String())
	}
	return
}

func expectbody(t *testing.T, b *Builder, sb *strings.Builder, body string) {
	expect := string(DefaultHeader) + body
	if sb.String() != expect {
		t.Fatalf("Expected '%s' got '%s'", expect, sb.String())
	}
	if len(expect) != b.Written() {
		t.Fatalf("Expected bytes written %v got %v", len(expect), b.Written())
	}
	if sb.Len() != b.Written() {
		t.Fatalf("Expected string size %v got %v", sb.Len(), b.Written())
	}
}

func TestBuilderSimpleNode(t *testing.T) {
	var buffer strings.Builder
	b := newbuilder(t, &buffer)

	b.Tag([]byte("node"))
	if b.Error() != nil {
		t.Fatalf("Tag(), expected error '%v' got '%v'", nil, b.Error())
	}

	b.Close()
	if b.Error() != nil {
		t.Fatalf("Close(), expected error '%v' got '%v'", nil, b.Error())
	}

	expectbody(t, b, &buffer, "<node/>")
}

func TestBuilderSimpleNode2(t *testing.T) {
	var buffer strings.Builder
	b := newbuilder(t, &buffer)

	b.Tag([]byte("node"))
	if b.Error() != nil {
		t.Fatalf("Tag(), expected error '%v' got '%v'", nil, b.Error())
	}

	b.EndTag()
	if b.Error() != nil {
		t.Fatalf("Close(), expected error '%v' got '%v'", nil, b.Error())
	}

	expectbody(t, b, &buffer, "<node/>")
}

func TestBuilderSimpleTree(t *testing.T) {
	var buffer strings.Builder
	b := newbuilder(t, &buffer)

	b.Tag([]byte("parent"))
	if b.Error() != nil {
		t.Fatalf("Tag(), expected error '%v' got '%v'", nil, b.Error())
	}

	b.Tag([]byte("child"))
	if b.Error() != nil {
		t.Fatalf("Tag(), expected error '%v' got '%v'", nil, b.Error())
	}

	b.Close()
	if b.Error() != nil {
		t.Fatalf("Close(), expected error '%v' got '%v'", nil, b.Error())
	}

	expectbody(t, b, &buffer, "<parent><child/></parent>")
}

func TestBuilderSimpleTree2(t *testing.T) {
	var buffer strings.Builder
	b := newbuilder(t, &buffer)

	b.Tag([]byte("parent"))
	if b.Error() != nil {
		t.Fatalf("Tag(), expected error '%v' got '%v'", nil, b.Error())
	}

	b.Tag([]byte("child"))
	if b.Error() != nil {
		t.Fatalf("Tag(), expected error '%v' got '%v'", nil, b.Error())
	}

	b.TagEnd(false)
	if b.Error() != nil {
		t.Fatalf("TagEnd(), expected error '%v' got '%v'", nil, b.Error())
	}

	b.TagEnd(false)
	if b.Error() != nil {
		t.Fatalf("TagEnd(), expected error '%v' got '%v'", nil, b.Error())
	}

	expectbody(t, b, &buffer, "<parent><child/></parent>")
}

func TestBuilderNodeAttr(t *testing.T) {
	var buffer strings.Builder
	b := newbuilder(t, &buffer)

	b.Tag([]byte("node"))
	if b.Error() != nil {
		t.Fatalf("Tag(), expected error '%v' got '%v'", nil, b.Error())
	}

	b.Attr([]byte("attr0"), []byte("123"))

	b.Close()
	if b.Error() != nil {
		t.Fatalf("Close(), expected error '%v' got '%v'", nil, b.Error())
	}

	expectbody(t, b, &buffer, "<node attr0=\"123\"/>")
}

func TestBuilderText(t *testing.T) {
	var buffer strings.Builder
	b := newbuilder(t, &buffer)

	b.Tag([]byte("node"))
	if b.Error() != nil {
		t.Fatalf("Tag(), expected error '%v' got '%v'", nil, b.Error())
	}

	b.Text([]byte("this is a text"))

	b.Close()
	if b.Error() != nil {
		t.Fatalf("Close(), expected error '%v' got '%v'", nil, b.Error())
	}

	expectbody(t, b, &buffer, "<node>this is a text</node>")
}

func TestBuilderWrite(t *testing.T) {
	var buffer strings.Builder
	b := newbuilder(t, &buffer)

	b.Tag([]byte("node"))
	if b.Error() != nil {
		t.Fatalf("Tag(), expected error '%v' got '%v'", nil, b.Error())
	}

	text := "this is a text"
	n, err := b.Write([]byte(text))

	b.Close()
	if b.Error() != nil {
		t.Fatalf("Close(), expected error '%v' got '%v'", nil, b.Error())
	}

	if b.Error() != err {
		t.Fatalf("Write(), expected error '%v' got '%v'", err, b.Error())
	}

	if n != len(text) {
		t.Fatalf("Write(), expected write length %v got %v", len(text), n)
	}

	expectbody(t, b, &buffer, "<node>this is a text</node>")
}

func TestBuilderComment(t *testing.T) {
	var buffer strings.Builder
	b := newbuilder(t, &buffer)

	b.Tag([]byte("node"))
	if b.Error() != nil {
		t.Fatalf("Tag(), expected error '%v' got '%v'", nil, b.Error())
	}

	b.Comment([]byte("this is a comment"))

	b.Close()
	if b.Error() != nil {
		t.Fatalf("Close(), expected error '%v' got '%v'", nil, b.Error())
	}

	expectbody(t, b, &buffer, "<node><!-- this is a comment --></node>")
}

func TestBuilderCDATA(t *testing.T) {
	var buffer strings.Builder
	b := newbuilder(t, &buffer)

	b.Tag([]byte("node"))
	if b.Error() != nil {
		t.Fatalf("Tag(), expected error '%v' got '%v'", nil, b.Error())
	}

	b.CDATA([]byte("cdata content"))

	b.Close()
	if b.Error() != nil {
		t.Fatalf("Close(), expected error '%v' got '%v'", nil, b.Error())
	}

	expectbody(t, b, &buffer, "<node><![CDATA[cdata content]]></node>")
}

// Testing error branches.
type writer_maxcount struct {
	free int
	w    io.Writer
	err  error
}

func (w *writer_maxcount) Write(p []byte) (n int, err error) {
	if len(p) > w.free {
		n, err = w.w.Write(p[:w.free])
	} else {
		n, err = w.w.Write(p)
	}

	w.free -= n
	return
}

func newbuilder_writer(t *testing.T, buffer *strings.Builder, maxwrite int, err error) (b *Builder, w *writer_maxcount) {
	w = &writer_maxcount{maxwrite, buffer, err}

	if b = NewBuilder(w); b == nil {
		t.Fatal("NewBuilder returned nil")
	}

	return
}

func TestBuilderCloseTagError1(t *testing.T) {
	var buffer strings.Builder
	b := newbuilder(t, &buffer)

	b.EndTag()
	if b.Error() != ErrUnmatchedCloseTag {
		t.Fatalf("Close(), expected error '%v' got '%v'", ErrUnmatchedCloseTag, b.Error())
	}
}

func TestBuilderCloseTagError2(t *testing.T) {
	var buffer strings.Builder
	b := newbuilder(t, &buffer)

	b.Tag([]byte("node"))
	if b.Error() != nil {
		t.Fatalf("Tag(), expected error '%v' got '%v'", nil, b.Error())
	}

	b.EndTag()
	if b.Error() != nil {
		t.Fatalf("Close(), expected error '%v' got '%v'", nil, b.Error())
	}

	b.EndTag()
	if b.Error() != ErrUnmatchedCloseTag {
		t.Fatalf("Close(), expected error '%v' got '%v'", ErrUnmatchedCloseTag, b.Error())
	}
}

func TestBuilder_ErrShortWrite1(t *testing.T) {
	var buffer strings.Builder

	node := []byte("node")
	expectErr := io.ErrShortWrite

	b, w := newbuilder_writer(t, &buffer, 1000, expectErr)

	w.free = 0 // Too short for "<""
	b.Tag(node)

	if b.Error() != expectErr {
		t.Fatalf("Tag(), expected error '%v' got '%v'", expectErr, b.Error())
	}
}

func TestBuilder_ErrShortWrite2(t *testing.T) {
	var buffer strings.Builder

	node := []byte("node")
	expectErr := io.ErrShortWrite

	b, w := newbuilder_writer(t, &buffer, 1000, expectErr)

	w.free = len(node) - 2 // Too short for "<node/>"
	b.Tag(node)

	if b.Error() != expectErr {
		t.Fatalf("Tag(), expected error '%v' got '%v'", expectErr, b.Error())
	}
}

func TestBuilder_ErrShortWrite3(t *testing.T) {
	var buffer strings.Builder

	node := []byte("node")
	expectErr := io.ErrShortWrite

	b, w := newbuilder_writer(t, &buffer, 1000, expectErr)

	b.Tag(node)
	w.free = 1 // Stop after / in "<node/>"
	b.EndTag()

	if b.Error() != expectErr {
		t.Fatalf("EndTag(), expected error '%v' got '%v'", expectErr, b.Error())
	}
}

func TestBuilder_ErrShortWrite4(t *testing.T) {
	var buffer strings.Builder

	node := []byte("node")
	expectErr := io.ErrShortWrite

	b, w := newbuilder_writer(t, &buffer, 1000, expectErr)

	b.Tag(node)
	b.Text([]byte("text")) // Force proper close node (</node>)
	w.free = 1             // Stop after < in "</node>"
	b.EndTag()

	if b.Error() != expectErr {
		t.Fatalf("EndTag(), expected error '%v' got '%v'", expectErr, b.Error())
	}
}

func TestBuilder_ErrShortWrite5(t *testing.T) {
	var buffer strings.Builder

	node := []byte("node")
	expectErr := io.ErrShortWrite

	b, w := newbuilder_writer(t, &buffer, 1000, expectErr)

	b.Tag(node)
	b.Text([]byte("text")) // Force proper close node (</node>)
	w.free = 3             // Stop after </n in "</node>"
	b.EndTag()

	if b.Error() != expectErr {
		t.Fatalf("EndTag(), expected error '%v' got '%v'", expectErr, b.Error())
	}
}

func TestBuilder_ErrShortWrite6a(t *testing.T) {
	var buffer strings.Builder

	node := []byte("node")
	expectErr := io.ErrShortWrite

	b, w := newbuilder_writer(t, &buffer, 1000, expectErr)

	b.Tag(node)
	w.free = 0             // Stop prior to > in "<node>"
	b.Text([]byte("text")) // Force proper close node (</node>)

	b.Close()
	if b.Error() != expectErr {
		t.Fatalf("EndTag(), expected error '%v' got '%v'", expectErr, b.Error())
	}
}

func TestBuilder_ErrShortWrite6b(t *testing.T) {
	var buffer strings.Builder

	node := []byte("node")
	expectErr := io.ErrShortWrite

	b, w := newbuilder_writer(t, &buffer, 1000, expectErr)

	b.Tag(node)
	w.free = 0  // Stop prior to > in "<node>"
	b.Tag(node) // Force proper close node (</node>)

	b.Close()
	if b.Error() != expectErr {
		t.Fatalf("EndTag(), expected error '%v' got '%v'", expectErr, b.Error())
	}
}

func TestBuilder_ErrShortWrite6c(t *testing.T) {
	var buffer strings.Builder

	node := []byte("node")
	expectErr := io.ErrShortWrite

	b, w := newbuilder_writer(t, &buffer, 1000, expectErr)

	b.Tag(node)
	w.free = 0      // Stop prior to > in "<node>"
	b.Comment(node) // Force proper close node (</node>)

	b.Close()
	if b.Error() != expectErr {
		t.Fatalf("EndTag(), expected error '%v' got '%v'", expectErr, b.Error())
	}
}

func TestBuilder_ErrShortWrite6d(t *testing.T) {
	var buffer strings.Builder

	node := []byte("node")
	expectErr := io.ErrShortWrite

	b, w := newbuilder_writer(t, &buffer, 1000, expectErr)

	b.Tag(node)
	w.free = 0    // Stop prior to > in "<node>"
	b.CDATA(node) // Force proper close node (</node>)

	b.Close()
	if b.Error() != expectErr {
		t.Fatalf("EndTag(), expected error '%v' got '%v'", expectErr, b.Error())
	}
}

func TestBuilder_ErrShortWrite7a(t *testing.T) {
	var buffer strings.Builder

	node := []byte("node")
	expectErr := io.ErrShortWrite

	b, w := newbuilder_writer(t, &buffer, 1000, expectErr)

	b.Tag(node)
	w.free = 2              // Stop prior in CDATA section
	b.CDATA([]byte("text")) // Force proper close node (</node>)
	b.EndTag()

	if b.Error() != expectErr {
		t.Fatalf("EndTag(), expected error '%v' got '%v'", expectErr, b.Error())
	}
}

func TestBuilder_ErrShortWrite7b(t *testing.T) {
	var buffer strings.Builder

	node := []byte("node")
	expectErr := io.ErrShortWrite

	b, w := newbuilder_writer(t, &buffer, 1000, expectErr)

	b.Tag(node)
	w.free = 5              // Stop in middle of CDATA name
	b.CDATA([]byte("text")) // Force proper close node (</node>)
	b.EndTag()

	if b.Error() != expectErr {
		t.Fatalf("EndTag(), expected error '%v' got '%v'", expectErr, b.Error())
	}
}

func TestBuilder_ErrShortWrite7c(t *testing.T) {
	var buffer strings.Builder

	node := []byte("node")
	expectErr := io.ErrShortWrite

	b, w := newbuilder_writer(t, &buffer, 1000, expectErr)

	b.Tag(node)
	w.free = 9              // Stop in middle of CDATA name
	b.CDATA([]byte("text")) // Force proper close node (</node>)
	b.EndTag()

	if b.Error() != expectErr {
		t.Fatalf("EndTag(), expected error '%v' got '%v'", expectErr, b.Error())
	}
}

func TestBuilder_ErrShortWrite8a(t *testing.T) {
	var buffer strings.Builder

	node := []byte("node")
	expectErr := io.ErrShortWrite

	b, w := newbuilder_writer(t, &buffer, 1000, expectErr)

	b.Tag(node)
	w.free = 2                // Stop prior in comment section
	b.Comment([]byte("text")) // Force proper close node (</node>)
	b.EndTag()

	if b.Error() != expectErr {
		t.Fatalf("EndTag(), expected error '%v' got '%v'", expectErr, b.Error())
	}
}

func TestBuilder_ErrShortWrite8b(t *testing.T) {
	var buffer strings.Builder

	node := []byte("node")
	expectErr := io.ErrShortWrite

	b, w := newbuilder_writer(t, &buffer, 1000, expectErr)

	b.Tag(node)
	w.free = 5 + 2            // Stop in middle of comment value
	b.Comment([]byte("text")) // Force proper close node (</node>)
	b.EndTag()

	if b.Error() != expectErr {
		t.Fatalf("EndTag(), expected error '%v' got '%v'", expectErr, b.Error())
	}
}

func TestBuilder_ErrShortWrite8c(t *testing.T) {
	var buffer strings.Builder

	node := []byte("node")
	expectErr := io.ErrShortWrite

	b, w := newbuilder_writer(t, &buffer, 1000, expectErr)

	b.Tag(node)
	w.free = 5 + 4 + 1        // Stop in middle of comment end
	b.Comment([]byte("text")) // Force proper close node (</node>)
	b.EndTag()

	if b.Error() != expectErr {
		t.Fatalf("EndTag(), expected error '%v' got '%v'", expectErr, b.Error())
	}
}

func TestBuilder_ErrShortWrite9a(t *testing.T) {
	var buffer strings.Builder

	node := []byte("node")
	expectErr := io.ErrShortWrite

	b, w := newbuilder_writer(t, &buffer, 1000, expectErr)

	b.Tag(node)
	w.free = 0 // In attribute separator
	b.Attr(node, node)
	b.EndTag()

	if b.Error() != expectErr {
		t.Fatalf("EndTag(), expected error '%v' got '%v'", expectErr, b.Error())
	}
}

func TestBuilder_ErrShortWrite9b(t *testing.T) {
	var buffer strings.Builder

	node := []byte("node")
	expectErr := io.ErrShortWrite

	b, w := newbuilder_writer(t, &buffer, 1000, expectErr)

	b.Tag(node)
	w.free = 2 // In attribute name
	b.Attr(node, node)
	b.EndTag()

	if b.Error() != expectErr {
		t.Fatalf("EndTag(), expected error '%v' got '%v'", expectErr, b.Error())
	}
}

func TestBuilder_ErrShortWrite9c(t *testing.T) {
	var buffer strings.Builder

	node := []byte("node")
	expectErr := io.ErrShortWrite

	b, w := newbuilder_writer(t, &buffer, 1000, expectErr)

	b.Tag(node)
	w.free = 4 + 1 // In attribute assign (=")
	b.Attr(node, node)
	b.EndTag()

	if b.Error() != expectErr {
		t.Fatalf("EndTag(), expected error '%v' got '%v'", expectErr, b.Error())
	}
}

func TestBuilder_ErrShortWrite9d(t *testing.T) {
	var buffer strings.Builder

	node := []byte("node")
	expectErr := io.ErrShortWrite

	b, w := newbuilder_writer(t, &buffer, 1000, expectErr)

	b.Tag(node)
	w.free = 4 + 2 + 3 // In attribute value
	b.Attr(node, node)
	b.EndTag()

	if b.Error() != expectErr {
		t.Fatalf("EndTag(), expected error '%v' got '%v'", expectErr, b.Error())
	}
}

func TestBuilder_ErrShortWrite10a(t *testing.T) {
	var buffer strings.Builder

	node := []byte("node")
	expectErr := io.ErrShortWrite

	b, w := newbuilder_writer(t, &buffer, 1000, expectErr)

	text := "text"

	b.Tag(node)
	w.free = 0
	b.Write([]byte(text)) // Force proper close node (</node>)
	b.EndTag()

	if b.Error() != expectErr {
		t.Fatalf("EndTag(), expected error '%v' got '%v'", expectErr, b.Error())
	}
}

func TestBuilder_ErrShortWrite10b(t *testing.T) {
	var buffer strings.Builder

	node := []byte("node")
	expectErr := io.ErrShortWrite

	b, w := newbuilder_writer(t, &buffer, 1000, expectErr)

	text := "text"

	b.Tag(node)
	w.free = len(text)
	b.Write([]byte(text)) // Force proper close node (</node>)
	b.EndTag()

	if b.Error() != expectErr {
		t.Fatalf("EndTag(), expected error '%v' got '%v'", expectErr, b.Error())
	}
}
