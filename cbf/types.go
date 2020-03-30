package cbf

import (
	"bytes"
	"fmt"
	"io"
)

type GUID struct {
	DataA uint32
	DataB uint16
	DataC uint16
	DataD [8]byte
}

func (g GUID) String() string {
	return fmt.Sprintf("{%08X-%04X-%04X-%X}", g.DataA, g.DataB, g.DataC, g.DataD)
}

type FILETIME struct {
	LowDateTime  uint32 // Windows FILETIME structure
	HighDateTime uint32 // Windows FILETIME structure
}

type Sector []byte

func (s Sector) Reader() (r io.Reader) {
	return bytes.NewReader(s)
}
