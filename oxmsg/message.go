package oxmsg

import (
	"github.com/xianhammer/format/cbf"
)

type Message struct {
	*cbf.Document

	isUnicode  bool
	properties map[string][]*Entry
}

func New() *Message {
	return &Message{}
}

func NewFromFile(filename string) (m *Message, err error) {
	doc, err := cbf.NewFromFile(filename)
	if err == nil {
		m = new(Message)
		m.Document = doc
	}

	return
}

func (m *Message) Properties() map[string][]*Entry {
	if m.properties == nil {
		m.properties = make(map[string][]*Entry)
		m.Document.Walk(func(de *cbf.DirectoryEntry) {
			e := newEntry(de)
			n := e.Name()
			m.properties[n] = append(m.properties[n], e)
		})

		if eStore := m.Get(PidTagStoreSupportMask); eStore != nil {
			v, _ := eStore[0].Uint32()
			m.isUnicode = (v & 0x00040000) != 0
		}
	}
	return m.properties
}

func (m *Message) Walk(f func(d *cbf.DirectoryEntry)) {
	m.Document.Walk(f)
}

func (m *Message) Get(propertyID string) []*Entry {
	if m.properties == nil {
		m.Properties()
	}
	return m.properties[propertyID]
}

func (m *Message) Header() (h *Header, err error) {
	entries := m.Get(PidTagTransportMessageHeaders)
	if entries == nil {
		return nil, ErrPropertyNotFound
	}
	if len(entries) > 1 {
		return nil, ErrPropertyIllegalInstances
	}

	h = &Header{entries[0]}
	return
}
