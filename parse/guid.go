package parse

import (
	"bytes"
	"errors"
)

// ErrGUIDParse is returned if input couldn't be parsed
var ErrGUIDParse = errors.New("guid: Cannot parse")

var emptyPartGUID = []byte{}
var dashPartGUID = []byte{'-'}

// Guid parses a byte slice and turn it into a GUID (16 byte array). Any dashes are ignored.
func Guid(in []byte) (guid [16]byte, err error) {
	core := bytes.Replace(in, dashPartGUID, emptyPartGUID, -1)
	if len(core) != 32 {
		err = ErrGUIDParse
		return
	}

	for i := 0; i < 31; i += 2 {
		part, _ := Hex([]byte(core[i : i+2]))
		guid[i>>1] = byte(part)
	}
	return
}
