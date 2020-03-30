package oxmsg

import (
	"bufio"
	"io"
	"strings"
)

type Header struct {
	*Entry
}

func (h *Header) Received() (blocks []*Block, err error) {
	var r io.Reader
	if r, err = h.TypedReader(); err == nil {
		blocks = h.readBlocks(r)

		// Re-order to first "received" first
		l := len(blocks) - 1
		for i := 0; i < l/2; i++ {
			blocks[i], blocks[l-i] = blocks[l-i], blocks[i]
		}
	}
	return
}

func (h *Header) readBlocks(r io.Reader) (blocks []*Block) {
	scanner := bufio.NewScanner(r)
	split := bufio.ScanLines
	scanner.Split(split)

	var block *Block
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, sPrefix) {
			block = new(Block)
			block.add(line, false)
			blocks = append(blocks, block)
		} else {
			block.add(line, rContinuedLine.MatchString(line))
		}
	}
	return
}
