package excel

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"regexp"
)

// For EXCEL limitations, see https://support.office.com/en-us/article/excel-specifications-and-limits-1672b34d-7043-467e-8e27-269d656771c3

// Theme         ContentType = "application/vnd.openxmlformats-officedocument.theme+xml"
// Styles                    = "application/vnd.openxmlformats-officedocument.spreadsheetml.styles+xml"
// Relationships             = "application/vnd.openxmlformats-package.relationships+xml"
// XML                       = "application/xml"
// Workbook                  = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"
// App                       = "application/vnd.openxmlformats-officedocument.extended-properties+xml"
// Sheet                     = "application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"
// SharedStrings             = "application/vnd.openxmlformats-officedocument.spreadsheetml.sharedStrings+xml"
// Core                      = "application/vnd.openxmlformats-package.core-properties+xml"

type ContentType string

var (
	// See https://docs.microsoft.com/en-us/office/open-xml/structure-of-a-spreadsheetml-document
	EmbeddedExcel = regexp.MustCompile(`\.xlsx$`) // Only process top levele EXCEL files.
	EmbeddedZIP   = regexp.MustCompile(`\.zip$`)  // Only process top levele EXCEL files.
)

type Document struct {
	zf       *zip.ReadCloser
	files    Files
	Workbook *Workbook
}

func NewDocument() (doc *Document) {
	doc = new(Document)
	doc.Workbook = newWorkbook()
	return
}

// OpenFile open a zip container
func OpenFile(path string) (doc *Document, err error) {
	zf, err := zip.OpenReader(path)
	if err != nil {
		return
	}

	doc = NewDocument()
	err = doc.openReader(&zf.Reader)
	doc.zf = zf

	return
}

func Open(r io.ReaderAt, size int64) (doc *Document, err error) {
	zf, err := zip.NewReader(r, size)
	if err != nil {
		return
	}

	doc = NewDocument()
	err = doc.openReader(zf)

	return
}

// Close the document
func (d *Document) Close() (err error) {
	if d.zf == nil {
		return
	}
	return d.zf.Close()
}

func (d *Document) openReader(zf *zip.Reader) (err error) {
	if d.files, err = d.openZIP(zf); err != nil {
		return
	}

	if err = d.initialise(); err == io.EOF {
		err = nil
	}
	return
}

func (d *Document) openEmbeddedFile(file *zip.File) (files map[string]*File, err error) {
	reader, err := file.Open()
	if err != nil {
		return
	}
	defer reader.Close()

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}

	zipfile, err := zip.NewReader(bytes.NewReader(b), int64(file.UncompressedSize64))
	return d.openZIP(zipfile)
}

func (d *Document) openZIP(zipfile *zip.Reader) (files map[string]*File, err error) {
	files = make(map[string]*File)

	for _, file := range zipfile.File {
		if EmbeddedExcel.MatchString(file.Name) {
			return d.openEmbeddedFile(file)
		} else {
			files[file.Name] = (*File)(file)
		}
	}

	return
}

func (d *Document) initialise() (err error) {
	relationship, err := d.files["_rels/.rels"].QuerySelectorAll("Relationship")
	if err != nil {
		return
	}

	workbookReference := Elements(relationship).Get("Type", OfficeDocument)
	if workbookReference == nil || workbookReference["Target"] == "" {
		err = ErrMissingWorkbook
		return
	}

	err = d.Workbook.open(workbookReference["Target"], d.files)
	if err != nil {
		return
	}

	return
}
