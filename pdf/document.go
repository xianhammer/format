package pdf

type Document struct {
	Version   Version
	StartXRef *XRef
	// XRefs   []*XRef
}

func New(major, minor int) (d *Document, err error) {
	d = new(Document)
	err = d.Version.Set(major, minor)
	return
}

/*
func (d *Document) GetXRef(offset int64) (xref *XRef) {
	for _, xref = range d.XRefs {
		if xref.offset == offset {
			return
		}
	}
	return nil
}

func (d *Document) ReadXRef(r *bufio.Reader, offset int64) (err error) {
	xref := new(XRef)
	xref.offset = offset

	d.XRefs = append(d.XRefs, xref)
	return xref.Read(r)
}
*/
