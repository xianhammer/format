package cfb

import "errors"

const (
	// REGSECT 0x00000000 - 0xFFFFFFF9 Regular sector number.
	MAXREGSECT    uint32 = 0xFFFFFFFA // Maximum regular sector number.
	NotApplicable        = 0xFFFFFFFB // Reserved for future use.
	DIFSECT              = 0xFFFFFFFC // Specifies a DIFAT sector in the FAT.
	FATSECT              = 0xFFFFFFFD // Specifies a FAT sector in the FAT.
	ENDOFCHAIN           = 0xFFFFFFFE // End of a linked chain of sectors.
	FREESECT             = 0xFFFFFFFF // Specifies an unallocated sector in the FAT, Mini FAT, or DIFAT.
)

// Object types in storage (from AAF specifications)
const (
	STGTY_INVALID   byte = 0 // empty directory entry
	STGTY_STORAGE        = 1 // element is a storage object
	STGTY_STREAM         = 2 // element is a stream object
	STGTY_LOCKBYTES      = 3 // element is an ILockBytes object
	STGTY_PROPERTY       = 4 // element is an IPropertyStorage object
	STGTY_ROOT           = 5 // element is a root storage
)

const (
	DE_RED   byte = 0
	DE_BLACK      = 1
)

// Directory Entry IDs (from AAF specifications)
const (
	MAXREGSID uint32 = 0xFFFFFFFA // (-6) maximum directory entry ID
	NOSTREAM  uint32 = 0xFFFFFFFF // (-1) unallocated directory entry
)

var Signature = [8]byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}
var CLSID_NULL GUID

var (
	ErrDoNotUse         = errors.New("Do not use")
	ErrSignature        = errors.New("Bad signature")
	ErrCLSID            = errors.New("Bad header CLSID")
	ErrVersion          = errors.New("Bad version")
	ErrSectorShift      = errors.New("Invalid sector shift, expected 9 or 12")
	ErrMiniSectorShift  = errors.New("Invalid mini sector shift, expected 6")
	ErrMiniSectorCutoff = errors.New("Incorrect MiniStreamCutoffSize")
	ErrStorageType      = errors.New("Unknown storage type")
	ErrStorageSize      = errors.New("Unknown storage size")
	ErrNameLength       = errors.New("Invalid name property length")
	ErrSeekIndex        = errors.New("Bad seek index")
	ErrSectorSize       = errors.New("Bad sector read")
	ErrRootSibling      = errors.New("Root directory has illegal siblings")
	ErrNotStream        = errors.New("Entry is not a stream")
	// ErrDirectorSectors            = errors.New("Incorrect number of directory sectors")
	// ErrTransactionSignatureNumber = errors.New("TransactionSignatureNumber is not zero")
)
