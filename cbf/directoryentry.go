package cbf

import (
	"encoding/binary"
)

type DirectoryEntry struct {
	Directory

	id       uint32
	level    uint32
	doc      *Document
	parent   *DirectoryEntry
	children []*DirectoryEntry
}

func NewDirectoryEntry(d *Document, s *Stream) (dir *DirectoryEntry, err error) {
	dir = new(DirectoryEntry)
	dir.doc = d
	if err = binary.Read(s, d.byteOrder, &dir.Directory); err == nil {
		err = dir.Validate()
	}
	return
}

func (d *DirectoryEntry) ID() uint32 {
	return d.id
}

func (d *DirectoryEntry) Level() uint32 {
	return d.level
}

func (d *DirectoryEntry) FullName() string {
	if d.Parent() == nil || d.Parent().Parent() == nil { // Prevent "Root Entry"
		return d.Name()
	}

	return d.Parent().Name() + "/" + d.Name()
}

func (d *DirectoryEntry) Parent() *DirectoryEntry {
	return d.parent
}

func (d *DirectoryEntry) Children() []*DirectoryEntry {
	return d.children
}

func (d *DirectoryEntry) Walk(f func(d *DirectoryEntry)) {
	f(d)

	for _, child := range d.children {
		child.Walk(f)
	}
}

func (d *DirectoryEntry) Stream() (s *Stream, err error) {
	if d.Type == STGTY_STREAM {
		s, err = d.doc.stream(d.Start, d.Size)
	} else {
		err = ErrNotStream
	}
	return
}
