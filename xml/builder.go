package xml

import (
	"errors"
	"io"
)

var (
	attrassign   = []byte("=\"")
	tagautoclose = []byte("/>")
	tagclose     = []byte("</")
	commentstart = []byte("<!-- ")
	commentend   = []byte(" -->")
	sectionstart = []byte("<![")
	sectionend   = []byte("]]>")
	cdata        = []byte("CDATA")

	separator          = []byte(" ")
	attrend            = []byte("\"")
	tagopen            = []byte("<")
	tagend             = []byte(">")
	sectionstartfinish = []byte("[")

	// DefaultHeader set the default header line for XML output.
	DefaultHeader = []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n")

	// ErrUnmatchedCloseTag is returned if a node tag isn't properly matched by a closing node.
	ErrUnmatchedCloseTag = errors.New("Unmatched close tag")
)

type Builder struct {
	w            io.Writer
	written      int
	err          error
	nodepath     [][]byte
	closespecial []byte
	autoclose    bool
	parentopen   bool
}

// NewBuilder creates a new XML builder.
func NewBuilder(target io.Writer) (b *Builder) {
	b = new(Builder)
	b.w = target
	b.write(DefaultHeader)
	return
}

// Close end all currently open tags.
func (b *Builder) Close() (err error) {
	for len(b.nodepath) > 0 && b.closetag() == nil {
	}
	return b.err
}

// Written, number of bytes written to target writer.
func (b *Builder) Written() int {
	return b.written
}

// Error return the current error met.
func (b *Builder) Error() error {
	return b.err
}

func (b *Builder) write(value []byte) (err error) {
	var n int
	n, b.err = b.w.Write(value)
	b.written += n
	if n < len(value) && b.err == nil {
		b.err = io.ErrShortWrite
	}
	return b.err
}

func (b *Builder) closetag() (err error) {
	b.parentopen = false

	if b.closespecial != nil {
		b.write(b.closespecial)
		b.closespecial = nil
		return b.err
	}

	l := len(b.nodepath)
	if l == 0 {
		b.err = ErrUnmatchedCloseTag
		return b.err
	}

	nodename := b.nodepath[l-1]
	if l > 1 {
		b.nodepath = b.nodepath[:l-1]
	} else {
		b.nodepath = nil
	}

	if b.autoclose {
		b.autoclose = false
		return b.write(tagautoclose)
	}

	if b.write(tagclose) != nil {
		return b.err
	}

	if b.write(nodename) != nil {
		return b.err
	}

	return b.write(tagend)
}

// AddAttribute is a wrapper for Attribute(nil, name, value).
func (b *Builder) Attr(name, value []byte) {
	b.Attribute(nil, name, value)
}

// CDATA wraps a Section write and a Text write (of content).
func (b *Builder) CDATA(content []byte) {
	if b.Section(cdata); b.err != nil {
		return
	}
	b.Text(content)
}

// EndTag is a simple wrapper TagEnd(false).
func (b *Builder) EndTag() {
	b.TagEnd(false)
}

// Tag corresponds to the Tokenizer.Tag method, writing starting a tag (named node).
func (b *Builder) Tag(name []byte) {
	if b.parentopen {
		if b.write(tagend) != nil {
			return
		}
		b.parentopen = false
	}

	if b.write(tagopen) != nil {
		return
	}

	n := make([]byte, len(name))
	copy(n, name)
	b.nodepath = append(b.nodepath, n)

	if b.write(name) == nil {
		b.autoclose = true
		b.parentopen = true
	}
}

// TagEnd corresponds to the Tokenizer.TagEnd method, writing a tag end.
func (b *Builder) TagEnd(autoclose bool) {
	b.closetag()
}

// Attribute corresponds to the Tokenizer.Attribute method, writing an attribute name and value.
func (b *Builder) Attribute(tag, name, value []byte) {
	if b.write(separator) != nil {
		return
	}

	if b.write(name) != nil {
		return
	}

	if b.write(attrassign) != nil {
		return
	}

	if b.write(value) != nil {
		return
	}

	b.write(attrend)
}

// Text corresponds to the Tokenizer.Text method, writing a text node.
// An open tag declaration will be closed.
func (b *Builder) Text(value []byte) {
	if b.parentopen {
		if b.write(tagend) != nil {
			return
		}
		b.parentopen = false
	}

	b.autoclose = false
	b.write(value)
}

// Comment corresponds to the Tokenizer.Comment method, writing a comment node.
// An open tag declaration will be closed.
func (b *Builder) Comment(value []byte) {
	if b.parentopen {
		if b.write(tagend) != nil {
			return
		}
		b.parentopen = false
	}

	b.autoclose = false
	if b.write(commentstart) != nil {
		return
	}

	if b.write(value) != nil {
		return
	}

	b.write(commentend)
	b.autoclose = false

}

// Section corresponds to the Tokenizer.Section method, writing a section (CDATA, INCLUDE, ...).
// An open tag declaration will be closed.
func (b *Builder) Section(name []byte) {
	if b.parentopen {
		if b.write(tagend) != nil {
			return
		}
		b.parentopen = false
	}

	b.autoclose = false
	if b.write(sectionstart) != nil {
		return
	}

	if b.write(name) != nil {
		return
	}

	if b.write(sectionstartfinish) != nil {
		return
	}

	b.closespecial = sectionend
}

// Write interface implementation for TEXT entry writing.
// Mainly intended to support eg. Fprintf()
// Note calling this may corrupt the output (if it contains XML key-characters).
func (b *Builder) Write(value []byte) (n int, err error) {
	if b.parentopen {
		if b.write(tagend) != nil {
			return
		}
		b.parentopen = false
	}

	b.autoclose = false
	n, b.err = b.w.Write(value)
	b.written += n
	if n < len(value) && b.err == nil {
		b.err = io.ErrShortWrite
	}
	err = b.err
	return
}
