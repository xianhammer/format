package olk15

import (
	"io"
	"os"
	"strings"
)

type Contacts struct {
	mails, firstnames, lastnames []string
}

func NewContacts() (c *Contacts) {
	return new(Contacts)
}

func ContactsFromFile(filename string) (c *Contacts, err error) {
	r, err := os.Open(filename)
	if err != nil {
		return
	}
	defer r.Close()

	c = NewContacts()
	err = c.ReadFrom(r)
	return
}

func (c *Contacts) Len() int {
	return len(c.mails)
}

func (c *Contacts) Name(i int) (name string) {
	return strings.TrimSpace(c.firstnames[i] + " " + c.lastnames[i])
}

func (c *Contacts) Mail(i int) (mail string) {
	return strings.TrimSpace(c.mails[i])
}

func (c *Contacts) IndexMail(s string) int {
	for i, v := range c.mails {
		if v == s {
			return i
		}
	}
	return -1
}

/*func (c *Contacts) IndexName(s string) i {
	for i, v := range c.mails {
		if v == s {
			return i
		}
	}
	return -1
}*/

func (c *Contacts) ReadFrom(r io.Reader) (err error) {
	headerSize := 0x42
	header := make([]byte, headerSize)

	n, err := r.Read(header)
	if err != nil {
		return
	}
	if n != len(header) {
		return ErrHeaderSize
	}

	mails, err := c.readTable(r, header, 0x2E, nil)
	if err != nil {
		return
	}
	c.mails = append(c.mails, mails...)

	firstnames, err := c.readTable(r, header, 0x32, DecoderStrip(DecodeUTF8))
	if err != nil {
		return
	}
	c.firstnames = append(c.firstnames, firstnames...)

	lastnames, err := c.readTable(r, header, 0x036, DecoderStrip(DecodeUTF8))
	if err != nil {
		return
	}
	c.lastnames = append(c.lastnames, lastnames...)

	// TODO Check for data consistency - eg.
	return
}

func (c *Contacts) readTable(r io.Reader, header []byte, offset int, decoder func(to_decode []byte) string) (s []string, err error) {
	tableSize := int(header[offset]) + int(header[offset+1])<<8
	table, err := c.readBlock(r, tableSize)
	if err != nil {
		return
	}

	indexSize := int(header[offset+2]) + int(header[offset+3])<<8
	index, err := c.readBlock(r, indexSize)
	if err != nil {
		return
	}

	if decoder == nil {
		decoder = DecoderDefault
	}

	// Split table according to index (N x uint32)
	previous := 0
	for i := 4; i < indexSize; i += 4 {
		current := int(index[i]) + int(index[i+1])<<8 + int(index[i+2])<<16 + int(index[i+3])<<24
		s = append(s, decoder(table[previous:current]))
		previous = current
	}
	return
}

func (c *Contacts) readBlock(r io.Reader, size int) (b []byte, err error) {
	b = make([]byte, size)
	n, err := r.Read(b)
	if err == nil && n != size {
		err = ErrHeaderSize
	}
	return
}
