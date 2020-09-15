package olk15

type Message struct {
	Magic          [8]byte
	Unknown01      [24]byte
	CRLM           uint32
	Unknown2       uint32
	Unknown3       uint32
	SizeBlock1Size uint32
	SizeBlock2Size uint32 // SizeBlock2Size + SizeBlock1Size = msg size
}

type ContactRecord struct {
	Type        byte // 0x02 = Group, 0x00 = Entry
	RecordLen   uint16
	Unknown01   [28]byte
	Email, Name String // TODO Figure out why String has a 4 byte length while records (containing strings) has 2 bytes...
}

type ContactGroup struct { // Senders, Receivers, CCs, BCcs?
	Unknown     uint32 // 4 býtes: 00000000 (Recv+CC), 01000001 (Sender)
	RecordCount uint32 // 4 býtes

}

type Attachment struct {
	Unknown01 uint32   // 03 00 00 00
	Guid      [16]byte // 78 F7 FF F5 E2 D7 49 82 9B 54 F4 A0 5E 2E 52 CC
}

func ParseMessage(b []byte) (m *Message, err error) {
	m = new(Message)

	return
}

// Contact

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
