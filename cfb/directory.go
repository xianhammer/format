package cfb

import (
	"fmt"
	"unicode"
)

// Directory is an OLE directory entry header. Names are not really kept here as they are not clearly descibing the field.
type Directory struct {
	NameUTF16    [32]uint16
	NameLength   uint16
	Type         byte
	Flags        byte // DECOLOR (Red/Black-tree color, red=0)
	LeftSibling  uint32
	RightSibling uint32
	Child        uint32
	CLSID        GUID   // If Type=STGTY_STORAGE
	UserFlags    uint32 // If Type=STGTY_STORAGE
	Time         [2]FILETIME
	Start        uint32
	Size         uint32
	PropType     uint16
}

func (d *Directory) Validate() (err error) {
	// STGTY_EMPTY     uint8 = 0 // empty directory entry
	// STGTY_STORAGE   uint8 = 1 // element is a storage object
	// STGTY_STREAM    uint8 = 2 // element is a stream object
	// STGTY_LOCKBYTES uint8 = 3 // element is an ILockBytes object
	// STGTY_PROPERTY  uint8 = 4 // element is an IPropertyStorage object
	// STGTY_ROOT      uint8 = 5 // element is a root storage

	if d.Type > STGTY_ROOT {
		fmt.Printf("Unexpected type: d.Type=%d\n", d.Type)
		return ErrStorageType
	}

	if d.NameLength > 64 {
		return ErrNameLength
	}

	if d.Type == STGTY_STORAGE && d.Size != 0 {
		return ErrStorageSize
	}

	return
}

func (d *Directory) Name() (name string) {
	if d.NameLength <= 0 {
		return ""
	}

	offset := 0
	if !unicode.IsPrint(rune(d.NameUTF16[0])) {
		offset = 1
	}
	return UTF16String(d.NameUTF16[offset : (d.NameLength>>1)-1])
}
