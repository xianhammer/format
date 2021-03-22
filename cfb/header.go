package cfb

// Header the OLE file header.
type Header struct {
	Signature        [8]byte
	CLSID            GUID
	MinorVersion     uint16
	MajorVersion     uint16
	ByteOrder        uint16
	SectorShift      uint16
	MiniSectorShift  uint16
	_                uint16 // Reserved, should be 0
	_                uint32 // Reserved, should be 0
	_                uint32 // Reserved, should be 0
	SectFAT          uint32
	SectDirStart     uint32
	_                uint32 // Transaction signature, should be 0
	MiniSectorCutoff uint32
	MiniFatStart     uint32
	MiniFat          uint32
	DifStart         uint32
	Dif              uint32
	Fat              [109]uint32 // Double Indexed Fat
}
