package olk15

import "fmt"

const offsetMCXE = 0x0C
const sizeMCXE = 0x0C

type MCXE struct {
	// Ignore01 [8]byte
	// Ignore02 [4]byte // uint16
}

// Block seems to be 55h bytes long
// 4D 63 78 45 38 A4 01 00 00 00 00 00 03 00 00 00 C3 02 01 00 00 00 00 00 00 00 00 00 9C 31 01 00 00 00 00 00 00 00 00 00 7B 00 00 00 01 01 D5 68 C3 FA 51 98 34 6B B6 A5 2C 4C B7 BF 33 E5 93 42 F6 F0 A7 26 D9 4E 80 80 00 02 0E 80 80 00 01 65 00 01 00 00 01

func (m *MCXE) parse(b []byte) (n int, err error) {
	// copy(h.Ignore01[:], b[0:8])
	// copy(h.Ignore02[:], b[8:12])
	fmt.Printf("\tMCXE %x\n", b[0:16])
	return sizeMCXE, nil // TODO Reflect realtity
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
