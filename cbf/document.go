package cbf

import (
	"encoding/binary"
	"io"
	"os"
)

type Document struct {
	Header

	byteOrder      binary.ByteOrder
	sectorSize     uint32
	miniSectorSize uint32
	sectors        []Sector
	fat            []uint32
	minifat        []uint32
	ministream     *Stream
	root           *DirectoryEntry
}

func New() (d *Document, err error) {
	d = new(Document)

	d.MajorVersion = 3
	d.MinorVersion = 0x003E

	d.SectorShift = 9
	if d.MajorVersion == 4 {
		d.SectorShift = 12
	}

	d.MiniSectorShift = 6
	d.MiniSectorCutoff = 0x1000

	d.initDocument()

	return
}

func NewFromFile(filename string) (d *Document, err error) {
	if d, err = New(); err != nil {
		return
	}

	f, err := os.Open(filename)
	if err != nil {
		return
	}

	_, err = d.ReadFrom(f)
	f.Close()
	return
}

func (d *Document) ReadFrom(r io.Reader) (n int64, err error) {
	if err = binary.Read(r, binary.LittleEndian, &d.Header); err != nil {
		return
	}

	d.initDocument()
	if n, err = d.readSectors(r); err != nil {
		return
	}

	if err = d.buildFAT(d.Fat[:]); err != nil {
		return
	}

	return
}

func (d *Document) Root() (root *DirectoryEntry, err error) {
	if d.root != nil {
		return d.root, nil
	}

	s, err := d.stream(d.SectDirStart, 0)
	if err != nil {
		return
	}

	if root, err = NewDirectoryEntry(d, s); err != nil {
		return
	}
	root.id = 0

	if err = d.buildMiniFAT(root.Start, root.Size); err != nil {
		return
	}

	if root.LeftSibling != NOSTREAM || root.RightSibling != NOSTREAM {
		err = ErrRootSibling
		return
	}

	d.root = root
	return root, d.entry(s, root, root.Child)
}

func (d *Document) Walk(f func(d *DirectoryEntry)) {
	root, err := d.Root()
	if err != nil {
		return
	}
	root.Walk(f)
}

func (d *Document) buildFAT(src []uint32) (err error) {
	var sectorIDs [128]uint32
	for _, sID := range src {
		if sID == ENDOFCHAIN {
			break
		}
		if sID == FREESECT {
			continue
		}
		if err = d.readBinary(sID, &sectorIDs); err != nil {
			return
		}

		d.fat = append(d.fat, sectorIDs[:]...)
	}

	if len(d.fat) > len(d.sectors) {
		d.fat = d.fat[:len(d.sectors)]
	}

	return
}

func (d *Document) buildMiniFAT(streamStart, size uint32) (err error) {
	if d.MiniFat <= 0 {
		return
	}

	var sectorIDs [128]uint32
	sID := d.MiniFatStart
	for ; sID <= MAXREGSECT; sID = d.fat[sID] {
		if err = d.readBinary(sID, &sectorIDs); err != nil {
			return
		}
		d.minifat = append(d.minifat, sectorIDs[:]...)
	}

	// The minifat is shortened, but has really no need
	minifatSize := size / d.miniSectorSize
	if uint32(len(d.minifat)) > minifatSize {
		d.minifat = d.minifat[:minifatSize]
	}

	// Prepare the ministream
	d.ministream = NewStream(d.sectorSize, size)
	for sID = streamStart; sID <= MAXREGSECT; sID = d.fat[sID] {
		d.ministream.add(d.sectors[sID], false)
	}

	return
}

func (d *Document) stream(sID, size uint32) (s *Stream, err error) {
	s = NewStream(d.sectorSize, size)
	addSize := s.size == 0
	if 0 < size && size < d.MiniSectorCutoff {
		s.sectorSize = d.miniSectorSize

		var n int
		for ; sID <= MAXREGSECT; sID = d.minifat[sID] {
			if _, err = d.ministream.Seek(int64(sID*d.miniSectorSize), 0); err != nil {
				break
			}

			data := make(Sector, d.miniSectorSize)
			s.add(data, addSize)
			if n, err = d.ministream.Read(data[:]); err != nil {
				break
			} else if uint32(n) != d.miniSectorSize {
				err = ErrSectorSize
				break
			}
		}
		return
	}

	for ; sID <= MAXREGSECT; sID = d.fat[sID] {
		s.add(d.sectors[sID], addSize)
	}

	return
}

func (d *Document) entry(s *Stream, parent *DirectoryEntry, dirIndex uint32) (err error) {
	if dirIndex == NOSTREAM {
		return
	}

	if _, err = s.Seek(int64(128*dirIndex), 0); err != nil {
		return
	}

	var dir *DirectoryEntry
	if dir, err = NewDirectoryEntry(d, s); err != nil {
		return
	}

	dir.id = dirIndex
	dir.parent = parent
	if parent != nil {
		dir.level = parent.level + 1
	}
	parent.children = append(parent.children, dir)

	if err = d.entry(s, parent, dir.LeftSibling); err != nil {
		return
	}

	if err = d.entry(s, parent, dir.RightSibling); err != nil {
		return
	}

	if dir.Child == NOSTREAM || dir.Type != STGTY_STORAGE {
		return
	}

	return d.entry(s, dir, dir.Child)
}

func (d *Document) initDocument() {
	d.byteOrder = binary.BigEndian
	if d.ByteOrder == 0xFFFE {
		d.byteOrder = binary.LittleEndian
	}

	d.sectorSize = uint32(1) << d.SectorShift
	d.miniSectorSize = uint32(1) << d.MiniSectorShift
}

func (d *Document) doValidate() (err error) {
	if d.Signature != Signature {
		return ErrSignature
	}
	if d.CLSID != CLSID_NULL {
		return ErrCLSID
	}

	if (d.MajorVersion != 3 && d.MajorVersion != 4) || d.MinorVersion != 0x003E {
		return ErrVersion
	}

	if (d.MajorVersion == 3 && d.SectorShift != 9) || (d.MajorVersion == 4 && d.SectorShift != 12) {
		return ErrSectorShift
	}

	if d.MiniSectorShift != 6 {
		return ErrMiniSectorShift
	}

	if d.MiniSectorCutoff != 0x1000 {
		return ErrMiniSectorCutoff
	}

	return
}

func (d *Document) readSectors(r io.Reader) (n int64, err error) {
	var read int
	for err == nil {
		s := Sector(make([]byte, d.sectorSize))
		if read, err = r.Read(s); err == nil {
			if uint32(read) != d.sectorSize {
				err = ErrSectorSize
				return
			}
			n += int64(read)
			d.sectors = append(d.sectors, s)
		}
	}

	if err == io.EOF {
		err = nil
	}

	return
}

func (d *Document) readBinary(sID uint32, target interface{}) (err error) {
	if err = binary.Read(d.sectors[sID].Reader(), d.byteOrder, target); err == io.EOF {
		err = nil
	}
	return
}
