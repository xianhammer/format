package olk15

type Header struct {
	Magic    [4]byte
	Ignore01 [28]byte
}

func (h *Header) parse(b []byte) (n int, err error) {
	copy(h.Magic[:], b[0:4])
	copy(h.Ignore01[:], b[4:32])
	return 32, nil
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
