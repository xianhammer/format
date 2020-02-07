package excel

import (
	"io"
	"path"
)

// Workbook represent an Excel Workbook(.xml).
type Workbook struct {
	root          string
	files         Files
	types         Files
	sheets        []*Sheet
	styles        *Styles
	sharedstrings *SharedStrings
}

func newWorkbook() (wb *Workbook) {
	wb = new(Workbook)
	wb.styles = newStyles()
	wb.sharedstrings = newSharedStrings()
	return wb
}

// Sheet return a named sheet, the sheet file size or an error.
func (w *Workbook) Sheet(name string) (s *Sheet, n int64, err error) {
	if s = w.get(name); s == nil {
		err = ErrUnknownSheet
	} else if n, err = s.open(); err == io.EOF {
		err = nil
	}
	return
}

// Sheets return all sheets
func (w *Workbook) Sheets() (s []*Sheet) {
	return w.sheets
}

// AddSheet return a new named sheet.
func (w *Workbook) AddSheet(name string) (s *Sheet, err error) {
	if s = w.get(name); s != nil {
		err = ErrDuplicateSheet
	} else {
		s = newSheet(w, name)
		w.sheets = append(w.sheets, s)
	}
	return
}

// Import sheet (deep clone).
func (w *Workbook) ImportSheet(name string, s *Sheet, addCols []ColumnUpdate) (t *Sheet, err error) {
	if t, err = w.AddSheet(name); err != nil {
		return
	}

	t.appendSheet(s, addCols)
	return
}

// SharedStrings return a reference to the sharedstrings(.xml)
func (w *Workbook) SharedStrings() (ss *SharedStrings) {
	return w.sharedstrings
}

// Styles return a reference to the style(.xml)
func (w *Workbook) Styles() (ss *Styles) {
	return w.styles
}

func (w *Workbook) get(name string) (s *Sheet) {
	for _, sheet := range w.sheets {
		if sheet.Attributes["name"] == name {
			s = sheet
			return
		}
	}
	return
}

func (w *Workbook) open(target string, files Files) (err error) {
	var workbookfile string
	w.root, workbookfile = path.Split(target)

	rels := path.Join(w.root, "_rels", workbookfile+".rels")
	relationships, err := files[rels].QuerySelectorAll("Relationship")
	if err != nil {
		return
	}

	w.files = make(map[string]*File)
	w.types = make(map[string]*File)
	for _, element := range relationships {
		f := files[path.Join(w.root, element["Target"])]
		w.files[element["Id"]] = f
		w.types[element["Type"]] = f
	}

	// Get sheets
	sheets, err := files[target].QuerySelectorAll("sheet")
	if err != nil {
		return
	}

	for _, element := range sheets {
		sheet := newSheet(w, element["name"])
		sheet.file = w.files[element["r:id"]]

		for key, value := range element {
			sheet.Attributes[key] = value
		}
		sheet.Attributes = element
		w.sheets = append(w.sheets, sheet)
	}

	styles, exist := w.types[RelationshipStyles]
	if !exist {
		err = ErrMissingStyles
		return
	}

	sharedstrings, exist := w.types[RelationshipSharedstrings]
	if !exist {
		err = ErrMissingSharedstrings
		return
	}

	// w.styles = newStyles()
	// w.sharedstrings = new(SharedStrings)

	if err = w.styles.open(styles); err == nil {
		err = w.sharedstrings.open(sharedstrings)
	}
	return
}
