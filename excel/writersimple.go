package excel

import (
	"io"
	"os"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/xianhammer/format/xml"
	"github.com/xianhammer/system/zip"
)

type simplewriter struct {
	doc *Document
	zw  *zip.Writer
}

func NewSimpleWriter(doc *Document) (w *simplewriter) {
	return &simplewriter{
		doc,
		nil,
	}
}

func WriteSimple(outputFile string, doc *Document) (err error) {
	output, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return
	}
	defer output.Close()

	w := NewSimpleWriter(doc)
	_, err = w.WriteTo(output)
	return
}

func (s *simplewriter) WriteTo(w io.Writer) (n int64, err error) {
	s.zw = zip.NewWriter(w)
	defer s.zw.Close()

	for _, sheetInfo := range s.doc.Workbook.Sheets() {
		name := sheetInfo.Attributes["name"]
		_, _, err = s.doc.Workbook.Sheet(name)
		if err == io.EOF {
			err = nil
		}
		if err != nil {
			return
		}
	}

	if err = s.writeRels(); err != nil {
		return
	}

	if err = s.writeCoreProperties(); err != nil {
		return
	}

	if err = s.writeExtendedProperties(); err != nil {
		return
	}

	if err = s.writeContentTypes(); err != nil {
		return
	}

	if err = s.writeRelationships(); err != nil {
		return
	}

	if err = s.writeTheme(); err != nil {
		return
	}

	if err = s.writeWorkbook(); err != nil {
		return
	}

	if err = s.writeStyles(); err != nil {
		return
	}

	if err = s.writeSheets(); err != nil {
		return
	}

	if err = s.writeSharedStrings(); err != nil {
		return
	}

	return
}

func (s *simplewriter) writeRels() (err error) {
	w, err := s.zw.Create("_rels/.rels")
	if err != nil {
		return
	}

	b := xml.NewBuilder(w)
	defer b.Close()

	b.Tag([]byte("Relationships"))
	b.Attr([]byte("xmlns"), []byte(Relationships))

	b.Tag([]byte("Relationship"))
	b.Attr([]byte("Id"), []byte("rId1"))
	b.Attr([]byte("Type"), []byte(OfficeDocument))
	b.Attr([]byte("Target"), []byte("xl/workbook.xml"))
	b.EndTag()

	b.Tag([]byte("Relationship"))
	b.Attr([]byte("Id"), []byte("rId2"))
	b.Attr([]byte("Type"), []byte(CoreProperties))
	b.Attr([]byte("Target"), []byte("docProps/core.xml"))
	b.EndTag()

	b.Tag([]byte("Relationship"))
	b.Attr([]byte("Id"), []byte("rId3"))
	b.Attr([]byte("Type"), []byte(ExtendedProperties))
	b.Attr([]byte("Target"), []byte("docProps/app.xml"))
	b.EndTag()

	return b.Error()
}

func (s *simplewriter) writeCoreProperties() (err error) {
	w, err := s.zw.Create("docProps/core.xml")
	if err != nil {
		return
	}

	b := xml.NewBuilder(w)
	defer b.Close()

	b.Tag([]byte("cp:coreProperties"))
	defer b.EndTag()

	b.Attr([]byte("xmlns:cp"), []byte("http://schemas.openxmlformats.org/package/2006/metadata/core-properties"))
	b.Attr([]byte("xmlns:dc"), []byte("http://purl.org/dc/elements/1.1/"))
	b.Attr([]byte("xmlns:dcterms"), []byte("http://purl.org/dc/terms/"))
	b.Attr([]byte("xmlns:dcmitype"), []byte("http://purl.org/dc/dcmitype/"))
	b.Attr([]byte("xmlns:xsi"), []byte("http://www.w3.org/2001/XMLSchema-instance"))

	b.Tag([]byte("dc:creator"))
	b.Text([]byte(Creator))
	b.EndTag()

	b.Tag([]byte("cp:lastModifiedBy"))
	b.Text([]byte(LastModifiedBy)) // TODO
	b.EndTag()

	created := time.Now().UTC().Format(time.RFC3339)
	b.Tag([]byte("dcterms:created"))
	b.Attr([]byte("xsi:type"), []byte("dcterms:W3CDTF"))
	b.Text([]byte(created))
	b.EndTag()

	b.Tag([]byte("dcterms:modified"))
	b.Attr([]byte("xsi:type"), []byte("dcterms:W3CDTF"))
	b.Text([]byte(created))
	b.EndTag()

	return
}

func (s *simplewriter) writeExtendedProperties() (err error) {
	w, err := s.zw.Create("docProps/app.xml")
	if err != nil {
		return
	}

	b := xml.NewBuilder(w)
	defer b.Close()

	b.Tag([]byte("Properties"))
	defer b.EndTag()
	b.Attr([]byte("xmlns"), []byte("http://schemas.openxmlformats.org/officeDocument/2006/extended-properties"))
	b.Attr([]byte("xmlns:vt"), []byte("http://schemas.openxmlformats.org/officeDocument/2006/docPropsVTypes"))

	taggedText(b, "Application", "Microsoft Excel")
	taggedText(b, "DocSecurity", "0")
	taggedText(b, "ScaleCrop", "false")

	b.Tag([]byte("HeadingPairs"))
	b.Tag([]byte("vt:vector"))
	b.Attr([]byte("size"), []byte("2"))
	b.Attr([]byte("baseType"), []byte("variant"))

	b.Tag([]byte("vt:variant"))
	taggedText(b, "vt:lpstr", "Spreadsheet")
	b.EndTag()
	b.Tag([]byte("vt:variant"))
	taggedText(b, "vt:i4", strconv.Itoa(len(s.doc.Workbook.sheets)))
	b.EndTag() // </vt:variant>
	b.EndTag() // </vt:vector>
	b.EndTag() // </HeadingPairs>

	b.Tag([]byte("TitlesOfParts"))
	b.Tag([]byte("vt:vector"))

	b.Attr([]byte("size"), []byte(strconv.Itoa(len(s.doc.Workbook.sheets))))
	b.Attr([]byte("baseType"), []byte("lpstr"))

	for _, sheet := range s.doc.Workbook.sheets {
		taggedText(b, "vt:lpstr", sheet.Attributes["name"])
	}
	b.EndTag() // </vt:vector>
	b.EndTag() // </TitlesOfParts>

	taggedText(b, "LinksUpToDate", "false")
	taggedText(b, "SharedDoc", "false")
	taggedText(b, "HyperlinksChanged", "false")
	taggedText(b, "AppVersion", AppVersion)

	return
}

func (s *simplewriter) writeContentTypes() (err error) {
	w, err := s.zw.Create("[Content_Types].xml")
	if err != nil {
		return
	}

	b := xml.NewBuilder(w)
	defer func() {
		err = b.Close()
	}()

	b.Tag([]byte("Types"))
	b.Attr([]byte("xmlns"), []byte(ContentTypes))

	b.Tag([]byte("Override"))
	b.Attr([]byte("PartName"), []byte("/xl/theme/theme1.xml"))
	b.Attr([]byte("ContentType"), []byte("application/vnd.openxmlformats-officedocument.theme+xml"))
	b.EndTag()

	b.Tag([]byte("Override"))
	b.Attr([]byte("PartName"), []byte("/xl/styles.xml"))
	b.Attr([]byte("ContentType"), []byte("application/vnd.openxmlformats-officedocument.spreadsheetml.styles+xml"))
	b.EndTag()

	b.Tag([]byte("Default"))
	b.Attr([]byte("Extension"), []byte("rels"))
	b.Attr([]byte("ContentType"), []byte("application/vnd.openxmlformats-package.relationships+xml"))
	b.EndTag()

	b.Tag([]byte("Default"))
	b.Attr([]byte("Extension"), []byte("xml"))
	b.Attr([]byte("ContentType"), []byte("application/xml"))
	b.EndTag()

	b.Tag([]byte("Override"))
	b.Attr([]byte("PartName"), []byte("/xl/workbook.xml"))
	b.Attr([]byte("ContentType"), []byte("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"))
	b.EndTag()

	b.Tag([]byte("Override"))
	b.Attr([]byte("PartName"), []byte("/docProps/app.xml"))
	b.Attr([]byte("ContentType"), []byte("application/vnd.openxmlformats-officedocument.extended-properties+xml"))
	b.EndTag()

	for id := range s.doc.Workbook.sheets {
		b.Tag([]byte("Override"))
		b.Attr([]byte("PartName"), []byte("/xl/worksheets/sheet"+strconv.Itoa(id+1)+".xml"))
		b.Attr([]byte("ContentType"), []byte("application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"))
		b.EndTag()
	}

	b.Tag([]byte("Override"))
	b.Attr([]byte("PartName"), []byte("/xl/sharedStrings.xml"))
	b.Attr([]byte("ContentType"), []byte("application/vnd.openxmlformats-officedocument.spreadsheetml.sharedStrings+xml"))
	b.EndTag()

	b.Tag([]byte("Override"))
	b.Attr([]byte("PartName"), []byte("/docProps/core.xml"))
	b.Attr([]byte("ContentType"), []byte("application/vnd.openxmlformats-package.core-properties+xml"))
	b.EndTag()

	return
}

func (s *simplewriter) writeRelationships() (err error) {
	w, err := s.zw.Create("xl/_rels/workbook.xml.rels")
	if err != nil {
		return
	}

	b := xml.NewBuilder(w)
	defer func() {
		err = b.Close()
	}()

	b.Tag([]byte("Relationships"))
	b.Attr([]byte("xmlns"), []byte(Relationships))

	for id := range s.doc.Workbook.sheets {
		b.Tag([]byte("Relationship"))
		b.Attr([]byte("Id"), []byte("rId"+strconv.Itoa(id+1)))
		b.Attr([]byte("Type"), []byte(RelationshipsWorksheet))
		b.Attr([]byte("Target"), []byte("worksheets/sheet"+strconv.Itoa(id+1)+".xml"))
		b.EndTag()
	}

	nextId := len(s.doc.Workbook.sheets) + 1
	b.Tag([]byte("Relationship"))
	b.Attr([]byte("Id"), []byte("rId"+strconv.Itoa(nextId)))
	b.Attr([]byte("Type"), []byte(RelationshipStyles))
	b.Attr([]byte("Target"), []byte("styles.xml"))
	b.EndTag()

	nextId++
	b.Tag([]byte("Relationship"))
	b.Attr([]byte("Id"), []byte("rId"+strconv.Itoa(nextId)))
	b.Attr([]byte("Type"), []byte(RelationshipSharedstrings))
	b.Attr([]byte("Target"), []byte("sharedStrings.xml"))
	b.EndTag()

	nextId++
	b.Tag([]byte("Relationship"))
	b.Attr([]byte("Id"), []byte("rId"+strconv.Itoa(nextId)))
	b.Attr([]byte("Type"), []byte(RelationshipsTheme))
	b.Attr([]byte("Target"), []byte("theme/theme1.xml"))
	b.EndTag()

	return
}

func (s *simplewriter) writeTheme() (err error) {
	w, err := s.zw.Create("xl/theme/theme1.xml")
	if err != nil {
		return
	}

	b := xml.NewBuilder(w)
	defer func() {
		err = b.Close()
	}()

	b.Write([]byte(contentTheme))

	return
}

func (s *simplewriter) writeWorkbook() (err error) {
	w, err := s.zw.Create("xl/workbook.xml")
	if err != nil {
		return
	}

	b := xml.NewBuilder(w)
	defer func() {
		err = b.Close()
	}()

	b.Tag([]byte("workbook"))
	defer b.EndTag() // End <workbook>

	b.Attr([]byte("xmlns"), []byte(Spreadsheet))
	b.Attr([]byte("xmlns:r"), []byte(RelationshipsDoc))
	b.Attr([]byte("xmlns:mc"), []byte(MarkupCompatibility))
	b.Attr([]byte("xmlns:x15"), []byte(X15))
	b.Attr([]byte("mc:Ignorable"), []byte("x15"))

	b.Tag([]byte("fileVersion"))
	b.Attr([]byte("appName"), []byte("xl"))
	b.Attr([]byte("lastEdited"), []byte("4"))
	b.Attr([]byte("lowestEdited"), []byte("4"))
	b.Attr([]byte("rupBuild"), []byte("4507"))
	b.EndTag() // End fileVersion

	b.Tag([]byte("workbookPr"))
	b.Attr([]byte("defaultThemeVersion"), []byte("124226"))
	b.EndTag() // End workbookPr

	b.Tag([]byte("bookViews"))
	b.Tag([]byte("workbookView"))
	b.Attr([]byte("xWindow"), []byte("360"))
	b.Attr([]byte("yWindow"), []byte("300"))
	b.Attr([]byte("windowWidth"), []byte("14895"))
	b.Attr([]byte("windowHeight"), []byte("7875"))
	b.EndTag() // End workbookView
	b.EndTag() // End bookViews

	b.Tag([]byte("sheets"))
	for id, sheet := range s.doc.Workbook.sheets {
		b.Tag([]byte("sheet"))
		name := sheet.Attributes["name"]
		b.Attr([]byte("name"), []byte(name))
		b.Attr([]byte("sheetId"), []byte(strconv.Itoa(id+1)))
		b.Attr([]byte("r:id"), []byte("rId"+strconv.Itoa(id+1)))
		b.EndTag() // End sheet
	}
	b.EndTag() // End sheets

	b.Tag([]byte("calcPr"))
	b.Attr([]byte("calcId"), []byte("125725"))
	b.EndTag() // End calcPr

	return
}

func (s *simplewriter) writeSheets() (err error) {
	base := "xl/worksheets/sheet"
	var w io.Writer
	for id, sheet := range s.doc.Workbook.sheets {
		w, err = s.zw.Create(base + strconv.Itoa(id+1) + ".xml")
		if err != nil {
			break
		}

		_, err = sheet.WriteTo(w)
		if err != nil {
			break
		}
	}
	return
}

func (s *simplewriter) writeSharedStrings() (err error) {
	w, err := s.zw.Create("xl/sharedStrings.xml")
	if err != nil {
		return
	}

	b := xml.NewBuilder(w)
	defer func() {
		err = b.Close()
	}()

	b.Tag([]byte("sst"))
	defer b.EndTag() // End <sst>

	ss := s.doc.Workbook.sharedstrings
	count := strconv.Itoa(len(ss.strings))
	b.Attr([]byte("xmlns"), []byte(Spreadsheet))
	b.Attr([]byte("count"), []byte(count)) // TODO count is NOT equal to uniqueCount, often it is higher... Don't know what the number stands for.
	b.Attr([]byte("uniqueCount"), []byte(count))

	buf := make([]byte, 1024)
	bufIdx := 0

	for _, value := range ss.strings {
		b.Tag([]byte("si"))
		b.Tag([]byte("t"))

		bufIdx = 0
		for _, c := range value {
			if bufIdx >= len(buf)-5 {
				temp := make([]byte, len(buf)+1024)
				copy(temp, buf[:bufIdx])
				buf = temp
			}

			switch c {
			case '&':
				buf[bufIdx] = '&'
				buf[bufIdx+1] = 'a'
				buf[bufIdx+2] = 'm'
				buf[bufIdx+3] = 'p'
				buf[bufIdx+4] = ';'
				bufIdx += 5
			case '<', '>':
				buf[bufIdx] = '&'

				if c == '<' {
					buf[bufIdx+1] = 'l'
				} else {
					buf[bufIdx+1] = 'g'
				}

				buf[bufIdx+2] = 't'
				buf[bufIdx+3] = ';'
				bufIdx += 4

			default:
				bufIdx += utf8.EncodeRune(buf[bufIdx:], c)
			}
		}

		b.Text(buf[:bufIdx])

		b.EndTag()
		b.EndTag()
	}

	return
}

func (s *simplewriter) writeStyles() (err error) {
	w, err := s.zw.Create("xl/styles.xml")
	if err != nil {
		return
	}

	b := xml.NewBuilder(w)
	defer func() {
		err = b.Close()
	}()

	b.Tag([]byte("styleSheet"))
	defer b.EndTag() // End <styleSheet>
	b.Attr([]byte("xmlns"), []byte(Spreadsheet))
	b.Attr([]byte("xmlns:mc"), []byte(MarkupCompatibility))
	b.Attr([]byte("xmlns:x14ac"), []byte(X14ac))
	b.Attr([]byte("xmlns:x16r2"), []byte(X16r2))
	b.Attr([]byte("mc:Ignorable"), []byte("x14ac x16r2"))

	styles := s.doc.Workbook.styles

	err = s.writeStylesNumFmts(b, styles)
	if err != nil {
		return
	}

	err = s.writeStylesFonts(b, styles)
	if err != nil {
		return
	}

	err = s.writeStylesFills(b, styles)
	if err != nil {
		return
	}

	err = s.writeStylesBorders(b, styles)
	if err != nil {
		return
	}

	err = s.writeStylesCellStyleXfs(b, styles)
	if err != nil {
		return
	}

	err = s.writeStylesCellXfs(b, styles)
	if err != nil {
		return
	}

	err = s.writeStylesCellStyles(b, styles)
	if err != nil {
		return
	}

	// b.Tag([]byte("dxfs"))
	// b.Attr([]byte("count"), []byte("0"))
	// b.EndTag() // End <dxfs>

	b.Tag([]byte("tableStyles"))
	b.Attr([]byte("count"), []byte("0"))
	b.Attr([]byte("defaultTableStyle"), []byte("TableStyleMedium9"))
	b.Attr([]byte("defaultPivotStyle"), []byte("PivotStyleLight16"))
	b.EndTag() // End <tableStyles>

	return
}

func (s *simplewriter) writeStylesNumFmts(b *xml.Builder, styles *Styles) (err error) {
	numFmts := make([]*numFmt, len(styles.numFmts))
	count := 0
	for _, nf := range styles.numFmts {
		if nf.IsCustom() {
			numFmts[count] = nf
			count++
		}
	}

	if count == 0 {
		return
	}

	// <numFmts count="5">...</numFmts>
	b.Tag([]byte("numFmts"))
	defer b.EndTag() // End <numFmts>

	b.Attr([]byte("count"), []byte(strconv.Itoa(count)))

	// <numFmt numFmtId="164" formatCode="dd-mm-yyyy hh:mm:ss"/>
	for _, nf := range numFmts[:count] {
		nf.toXMLBuilder(b)
	}

	return
}

func (s *simplewriter) writeStylesCellXfs(b *xml.Builder, styles *Styles) (err error) {
	b.Tag([]byte("cellXfs"))
	defer b.EndTag()

	count := len(styles.cellXfs)
	b.Attr([]byte("count"), []byte(strconv.Itoa(count)))

	for _, xf := range styles.cellXfs {
		// Reset to known fonts, etc..
		// TODO FIX.
		xf.FontId = 0
		xf.FillId = 0
		xf.BorderId = 0

		xf.toXMLBuilder(b)
	}

	return
}

// Write fixed set of fonts. That is, currently no font selection is available.
func (s *simplewriter) writeStylesFonts(b *xml.Builder, styles *Styles) (err error) {
	b.Tag([]byte("fonts"))
	defer b.EndTag() // End <fonts>
	b.Attr([]byte("count"), []byte("1"))

	b.Tag([]byte("font"))

	b.Tag([]byte("sz"))
	b.Attr([]byte("val"), []byte("11"))
	b.EndTag() // End <sz>

	b.Tag([]byte("name"))
	b.Attr([]byte("val"), []byte("Calibri"))
	b.EndTag() // End <name>

	b.EndTag() // End <font>

	return
}

// Write fixed set of fills. That is, currently no fill selection is available.
func (s *simplewriter) writeStylesFills(b *xml.Builder, styles *Styles) (err error) {
	b.Tag([]byte("fills"))
	defer b.EndTag() // End <fills>
	b.Attr([]byte("count"), []byte("2"))

	b.Tag([]byte("fill"))

	b.Tag([]byte("patternFill"))
	b.Attr([]byte("patternType"), []byte("none"))
	b.EndTag() // End <patternFill>

	b.EndTag() // End <fill>

	b.Tag([]byte("fill"))

	b.Tag([]byte("patternFill"))
	b.Attr([]byte("patternType"), []byte("gray125"))
	b.EndTag() // End <patternFill>

	b.EndTag() // End <fill>

	return
}

// Write fixed set of borders. That is, currently no border selection is available.
func (s *simplewriter) writeStylesBorders(b *xml.Builder, styles *Styles) (err error) {
	b.Tag([]byte("borders"))
	defer b.EndTag() // End <borders>
	b.Attr([]byte("count"), []byte("1"))

	b.Tag([]byte("border"))

	b.Tag([]byte("left"))
	b.EndTag() // End <left>
	b.Tag([]byte("right"))
	b.EndTag() // End <right>
	b.Tag([]byte("top"))
	b.EndTag() // End <top>
	b.Tag([]byte("bottom"))
	b.EndTag() // End <bottom>
	b.Tag([]byte("diagonal"))
	b.EndTag() // End <diagonal>

	b.EndTag() // End <border>

	return
}

// Write fixed set of cellStyleXfs. That is, currently no cellStyleXf selection is available.
func (s *simplewriter) writeStylesCellStyleXfs(b *xml.Builder, styles *Styles) (err error) {
	b.Tag([]byte("cellStyleXfs"))
	defer b.EndTag() // End <cellStyleXfs>

	count := len(styles.cellStyleXfs)
	b.Attr([]byte("count"), []byte(strconv.Itoa(count)))

	for _, xf := range styles.cellStyleXfs {
		// Reset to known fonts, etc...
		// TODO FIX!
		xf.FontId = 0
		xf.FillId = 0
		xf.BorderId = 0
		xf.toXMLBuilder(b)
	}

	// b.Attr([]byte("count"), []byte("1"))
	// b.Tag([]byte("xf"))
	// b.Attr([]byte("numFmtId"), []byte("0"))
	// b.Attr([]byte("fontId"), []byte("0"))
	// b.Attr([]byte("fillId"), []byte("0"))
	// b.Attr([]byte("borderId"), []byte("0"))
	// b.EndTag() // End <xf>

	return
}

// Write fixed set of cellStyles. That is, currently no cellStyle selection is available.
func (s *simplewriter) writeStylesCellStyles(b *xml.Builder, styles *Styles) (err error) {
	b.Tag([]byte("cellStyles"))
	defer b.EndTag() // End <cellStyles>
	b.Attr([]byte("count"), []byte("1"))

	b.Tag([]byte("cellStyle"))
	b.Attr([]byte("name"), []byte("Normal"))
	b.Attr([]byte("xfId"), []byte("0"))
	b.Attr([]byte("builtinId"), []byte("0"))
	b.EndTag() // End <cellStyle>

	return
}

func taggedText(b *xml.Builder, tag, text string) (err error) {
	b.Tag([]byte(tag))
	b.Text([]byte(text))
	b.EndTag()
	return b.Error()
}

var contentTheme = `<a:theme name="OfficeTheme" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
	<a:themeElements>
		<a:clrScheme name="Office">
			<a:dk1>
				<a:sysClr lastClr="000000" val="windowText"/>
			</a:dk1>
			<a:lt1>
				<a:sysClr lastClr="FFFFFF" val="window"/>
			</a:lt1>
			<a:dk2>
				<a:srgbClr val="1F497D"/>
			</a:dk2>
			<a:lt2>
				<a:srgbClr val="EEECE1"/>
			</a:lt2>
			<a:accent1>
				<a:srgbClr val="4F81BD"/>
			</a:accent1>
			<a:accent2>
				<a:srgbClr val="C0504D"/>
			</a:accent2>
			<a:accent3>
				<a:srgbClr val="9BBB59"/>
			</a:accent3>
			<a:accent4>
				<a:srgbClr val="8064A2"/>
			</a:accent4>
			<a:accent5>
				<a:srgbClr val="4BACC6"/>
			</a:accent5>
			<a:accent6>
				<a:srgbClr val="F79646"/>
			</a:accent6>
			<a:hlink>
				<a:srgbClr val="0000FF"/>
			</a:hlink>
			<a:folHlink>
				<a:srgbClr val="800080"/>
			</a:folHlink>
		</a:clrScheme>
		<a:fontScheme name="Office">
			<a:majorFont>
				<a:latin typeface="Cambria"/>
				<a:ea typeface=""/>
				<a:cs typeface=""/>
				<a:font script="Jpan" typeface="ＭＳ Ｐゴシック"/>
				<a:font script="Hang" typeface="맑은 고딕"/>
				<a:font script="Hans" typeface="宋体"/>
				<a:font script="Hant" typeface="新細明體"/>
				<a:font script="Arab" typeface="Times New Roman"/>
				<a:font script="Hebr" typeface="Times New Roman"/>
				<a:font script="Thai" typeface="Tahoma"/>
				<a:font script="Ethi" typeface="Nyala"/>
				<a:font script="Beng" typeface="Vrinda"/>
				<a:font script="Gujr" typeface="Shruti"/>
				<a:font script="Khmr" typeface="MoolBoran"/>
				<a:font script="Knda" typeface="Tunga"/>
				<a:font script="Guru" typeface="Raavi"/>
				<a:font script="Cans" typeface="Euphemia"/>
				<a:font script="Cher" typeface="Plantagenet Cherokee"/>
				<a:font script="Yiii" typeface="Microsoft Yi Baiti"/>
				<a:font script="Tibt" typeface="Microsoft Himalaya"/>
				<a:font script="Thaa" typeface="MV Boli"/>
				<a:font script="Deva" typeface="Mangal"/>
				<a:font script="Telu" typeface="Gautami"/>
				<a:font script="Taml" typeface="Latha"/>
				<a:font script="Syrc" typeface="Estrangelo Edessa"/>
				<a:font script="Orya" typeface="Kalinga"/>
				<a:font script="Mlym" typeface="Kartika"/>
				<a:font script="Laoo" typeface="DokChampa"/>
				<a:font script="Sinh" typeface="Iskoola Pota"/>
				<a:font script="Mong" typeface="Mongolian Baiti"/>
				<a:font script="Viet" typeface="Times New Roman"/>
				<a:font script="Uigh" typeface="Microsoft Uighur"/>
			</a:majorFont>
			<a:minorFont>
				<a:latin typeface="Calibri"/>
				<a:ea typeface=""/>
				<a:cs typeface=""/>
				<a:font script="Jpan" typeface="ＭＳ Ｐゴシック"/>
				<a:font script="Hang" typeface="맑은 고딕"/>
				<a:font script="Hans" typeface="宋体"/>
				<a:font script="Hant" typeface="新細明體"/>
				<a:font script="Arab" typeface="Arial"/>
				<a:font script="Hebr" typeface="Arial"/>
				<a:font script="Thai" typeface="Tahoma"/>
				<a:font script="Ethi" typeface="Nyala"/>
				<a:font script="Beng" typeface="Vrinda"/>
				<a:font script="Gujr" typeface="Shruti"/>
				<a:font script="Khmr" typeface="DaunPenh"/>
				<a:font script="Knda" typeface="Tunga"/>
				<a:font script="Guru" typeface="Raavi"/>
				<a:font script="Cans" typeface="Euphemia"/>
				<a:font script="Cher" typeface="Plantagenet Cherokee"/>
				<a:font script="Yiii" typeface="Microsoft Yi Baiti"/>
				<a:font script="Tibt" typeface="Microsoft Himalaya"/>
				<a:font script="Thaa" typeface="MV Boli"/>
				<a:font script="Deva" typeface="Mangal"/>
				<a:font script="Telu" typeface="Gautami"/>
				<a:font script="Taml" typeface="Latha"/>
				<a:font script="Syrc" typeface="Estrangelo Edessa"/>
				<a:font script="Orya" typeface="Kalinga"/>
				<a:font script="Mlym" typeface="Kartika"/>
				<a:font script="Laoo" typeface="DokChampa"/>
				<a:font script="Sinh" typeface="Iskoola Pota"/>
				<a:font script="Mong" typeface="Mongolian Baiti"/>
				<a:font script="Viet" typeface="Arial"/>
				<a:font script="Uigh" typeface="Microsoft Uighur"/>
			</a:minorFont>
		</a:fontScheme>
		<a:fmtScheme name="Office">
			<a:fillStyleLst>
				<a:solidFill>
					<a:schemeClr val="phClr"/>
				</a:solidFill>
				<a:gradFill rotWithShape="1">
					<a:gsLst>
						<a:gs pos="0">
							<a:schemeClr val="phClr">
								<a:tint val="50000"/>
								<a:satMod val="300000"/>
							</a:schemeClr>
						</a:gs>
						<a:gs pos="35000">
							<a:schemeClr val="phClr">
								<a:tint val="37000"/>
								<a:satMod val="300000"/>
							</a:schemeClr>
						</a:gs>
						<a:gs pos="100000">
							<a:schemeClr val="phClr">
								<a:tint val="15000"/>
								<a:satMod val="350000"/>
							</a:schemeClr>
						</a:gs>
					</a:gsLst>
					<a:lin ang="16200000" scaled="1"/>
				</a:gradFill>
				<a:gradFill rotWithShape="1">
					<a:gsLst>
						<a:gs pos="0">
							<a:schemeClr val="phClr">
								<a:shade val="51000"/>
								<a:satMod val="130000"/>
							</a:schemeClr>
						</a:gs>
						<a:gs pos="80000">
							<a:schemeClr val="phClr">
								<a:shade val="93000"/>
								<a:satMod val="130000"/>
							</a:schemeClr>
						</a:gs>
						<a:gs pos="100000">
							<a:schemeClr val="phClr">
								<a:shade val="94000"/>
								<a:satMod val="135000"/>
							</a:schemeClr>
						</a:gs>
					</a:gsLst>
					<a:lin ang="16200000" scaled="0"/>
				</a:gradFill>
			</a:fillStyleLst>
			<a:lnStyleLst>
				<a:ln algn="ctr" cap="flat" cmpd="sng" w="9525">
					<a:solidFill>
						<a:schemeClr val="phClr">
							<a:shade val="95000"/>
							<a:satMod val="105000"/>
						</a:schemeClr>
					</a:solidFill>
					<a:prstDash val="solid"/>
				</a:ln>
				<a:ln algn="ctr" cap="flat" cmpd="sng" w="25400">
					<a:solidFill>
						<a:schemeClr val="phClr"/>
					</a:solidFill>
					<a:prstDash val="solid"/>
				</a:ln>
				<a:ln algn="ctr" cap="flat" cmpd="sng" w="38100">
					<a:solidFill>
						<a:schemeClr val="phClr"/>
					</a:solidFill>
					<a:prstDash val="solid"/>
				</a:ln>
			</a:lnStyleLst>
			<a:effectStyleLst>
				<a:effectStyle>
					<a:effectLst>
						<a:outerShdw blurRad="40000" dir="5400000" dist="20000" rotWithShape="0">
							<a:srgbClr val="000000">
								<a:alpha val="38000"/>
							</a:srgbClr>
						</a:outerShdw>
					</a:effectLst>
				</a:effectStyle>
				<a:effectStyle>
					<a:effectLst>
						<a:outerShdw blurRad="40000" dir="5400000" dist="23000" rotWithShape="0">
							<a:srgbClr val="000000">
								<a:alpha val="35000"/>
							</a:srgbClr>
						</a:outerShdw>
					</a:effectLst>
				</a:effectStyle>
				<a:effectStyle>
					<a:effectLst>
						<a:outerShdw blurRad="40000" dir="5400000" dist="23000" rotWithShape="0">
							<a:srgbClr val="000000">
								<a:alpha val="35000"/>
							</a:srgbClr>
						</a:outerShdw>
					</a:effectLst>
					<a:scene3d>
						<a:camera prst="orthographicFront">
							<a:rot lat="0" lon="0" rev="0"/>
						</a:camera>
						<a:lightRig dir="t" rig="threePt">
							<a:rot lat="0" lon="0" rev="1200000"/>
						</a:lightRig>
					</a:scene3d>
					<a:sp3d>
						<a:bevelT h="25400" w="63500"/>
					</a:sp3d>
				</a:effectStyle>
			</a:effectStyleLst>
			<a:bgFillStyleLst>
				<a:solidFill>
					<a:schemeClr val="phClr"/>
				</a:solidFill>
				<a:gradFill rotWithShape="1">
					<a:gsLst>
						<a:gs pos="0">
							<a:schemeClr val="phClr">
								<a:tint val="40000"/>
								<a:satMod val="350000"/>
							</a:schemeClr>
						</a:gs>
						<a:gs pos="40000">
							<a:schemeClr val="phClr">
								<a:tint val="45000"/>
								<a:shade val="99000"/>
								<a:satMod val="350000"/>
							</a:schemeClr>
						</a:gs>
						<a:gs pos="100000">
							<a:schemeClr val="phClr">
								<a:shade val="20000"/>
								<a:satMod val="255000"/>
							</a:schemeClr>
						</a:gs>
					</a:gsLst>
					<a:path path="circle">
						<a:fillToRect b="180000" l="50000" r="50000" t="-80000"/>
					</a:path>
				</a:gradFill>
				<a:gradFill rotWithShape="1">
					<a:gsLst>
						<a:gs pos="0">
							<a:schemeClr val="phClr">
								<a:tint val="80000"/>
								<a:satMod val="300000"/>
							</a:schemeClr>
						</a:gs>
						<a:gs pos="100000">
							<a:schemeClr val="phClr">
								<a:shade val="30000"/>
								<a:satMod val="200000"/>
							</a:schemeClr>
						</a:gs>
					</a:gsLst>
					<a:path path="circle">
						<a:fillToRect b="50000" l="50000" r="50000" t="50000"/>
					</a:path>
				</a:gradFill>
			</a:bgFillStyleLst>
		</a:fmtScheme>
	</a:themeElements>
	<a:objectDefaults/>
	<a:extraClrSchemeLst/>
</a:theme>`
