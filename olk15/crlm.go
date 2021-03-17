package olk15

import (
	"encoding/binary"
)

const offsetCRLM = 0x20

type CRLM struct { // Only entries identified are used.
	Marker      [4]byte // CRLM
	Unknown01   uint32
	RecordCount uint32 // Number of "records" in a section not yet decoded.
	// Size        uint32   // CRLM data size
	// References []uint32 // Various reeferences, i think.
}

// func (c *CRLM) validate() bool { return true }

func (c *CRLM) parse(b []byte) (n int, err error) {
	copy(c.Marker[:], b[0:4])

	bo := binary.LittleEndian
	c.RecordCount = bo.Uint32(b[8:12])

	size := bo.Uint32(b[12:16])

	// data := b[16:]
	// for i, l := 0, int(size>>2); i < l; i++ {
	// 	c.References = append(c.References, bo.Uint32(data[i:i+4]))
	// }

	return int(size) + 16, nil
}

// func (h *Header) Received() (blocks []*Block, err error) {
// 	var r io.Reader
// 	if r, err = h.TypedReader(); err == nil {
// 		blocks = h.readBlocks(r)

// 		// Re-order to first "received" first
// 		l := len(blocks) - 1
// 		for i := 0; i < l/2; i++ {
// 			blocks[i], blocks[l-i] = blocks[l-i], blocks[i]
// 		}
// 	}
// 	return
// }

// func (h *Header) readBlocks(r io.Reader) (blocks []*Block) {
// 	scanner := bufio.NewScanner(r)
// 	split := bufio.ScanLines
// 	scanner.Split(split)

// 	var block *Block
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		if strings.HasPrefix(line, sPrefix) {
// 			block = new(Block)
// 			block.add(line, false)
// 			blocks = append(blocks, block)
// 		} else {
// 			block.add(line, rContinuedLine.MatchString(line))
// 		}
// 	}
// 	return
// }
