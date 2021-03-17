package olk15

type ContactManager struct {
	contacts []*Contacts
}

func NewContactManager() (m *ContactManager) {
	m = new(ContactManager)
	return
}

func (m *ContactManager) Len() (l int) {
	// return len(c.mails)
	for _, c := range m.contacts {
		l += c.Len()
	}
	return
}

func (m *ContactManager) Name(i int) (name string) {
	for _, c := range m.contacts {
		l := c.Len()
		if i < l {
			return c.Name(i)
		}
		i -= l
	}
	return
}

func (m *ContactManager) Mail(i int) (mail string) {
	for _, c := range m.contacts {
		l := c.Len()
		if i < l {
			return c.Mail(i)
		}
		i -= l
	}
	return
}

func (m *ContactManager) IndexMail(s string) int {
	for _, c := range m.contacts {
		if i := c.IndexMail(s); i >= 0 {
			return i
		}
	}
	return -1
}

func (m *ContactManager) Add(c *Contacts) {
	m.contacts = append(m.contacts, c)
}

func (m *ContactManager) AddFromFile(filename string) (err error) {
	// c, err := ContactsFromFile(file.Path())
	c, err := ContactsFromFile(filename)
	if err == nil {
		m.Add(c)
	}
	return
}
