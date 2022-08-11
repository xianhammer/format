package excel

import (
	"bufio"
	"io"
	"strconv"

	"github.com/xianhammer/format/xml"
)

var defaultNumFmts map[string]*numFmt

// TODO Proper implement all cases below
// https://exceljet.net/custom-number-formats
// Character	Purpose
// 0	Display insignificant zeros
// #	Display significant digits
// ?	Display aligned decimals
// .	Decimal point
// ,	Thousands separator
// *	Repeat digit
// _	Add space
// @	Placeholder for text

func init() {
	defaultNumFmts = make(map[string]*numFmt)
	defaultNumFmts["0"] = NewNumFmt("0", "", "", FormatDefault)
	defaultNumFmts["1"] = NewNumFmt("1", "0", "%d", FormatInteger)
	defaultNumFmts["2"] = NewNumFmt("2", "0.00", "%f", FormatFloat)
	defaultNumFmts["3"] = NewNumFmt("3", "#,##0", "%d", FormatInteger)
	defaultNumFmts["4"] = NewNumFmt("4", "#,##0.00", "%f", FormatFloat)

	defaultNumFmts["9"] = NewNumFmt("9", "0%", "%d%%", FormatInteger)
	defaultNumFmts["10"] = NewNumFmt("10", "0.00%", "%f%%", FormatFloat)
	defaultNumFmts["11"] = NewNumFmt("11", "0.00E+00", "%e", FormatFloat)
	defaultNumFmts["12"] = NewNumFmt("12", "# ?/?", "%d", FormatInteger)
	defaultNumFmts["13"] = NewNumFmt("13", "# ??/??", "%d", FormatInteger)
	defaultNumFmts["14"] = NewNumFmt("14", "mm-dd-yy", "", nil)
	defaultNumFmts["15"] = NewNumFmt("15", "d-mmm-yy", "", nil)
	defaultNumFmts["16"] = NewNumFmt("16", "d-mmm", "", nil)
	defaultNumFmts["17"] = NewNumFmt("17", "mmm-yy", "", nil)
	defaultNumFmts["18"] = NewNumFmt("18", "h:mm", "", nil)    // "h:mm AM/PM"
	defaultNumFmts["19"] = NewNumFmt("19", "h:mm:ss", "", nil) // "h:mm:ss AM/PM"
	defaultNumFmts["20"] = NewNumFmt("20", "h:mm", "", nil)
	defaultNumFmts["21"] = NewNumFmt("21", "h:mm:ss", "", nil)
	defaultNumFmts["22"] = NewNumFmt("22", "m/d/yy h:mm", "", nil)

	defaultNumFmts["37"] = NewNumFmt("37", "#,##0 ;(#,##0)", "%d", FormatInteger)
	defaultNumFmts["38"] = NewNumFmt("38", "#,##0 ;[Red](#,##0)", "%d", FormatInteger)
	defaultNumFmts["39"] = NewNumFmt("39", "#,##0.00;(#,##0.00)", "%f", FormatFloat)
	defaultNumFmts["40"] = NewNumFmt("40", "#,##0.00;[Red](#,##0.00)", "%f", FormatFloat)

	defaultNumFmts["45"] = NewNumFmt("45", "mm:ss", "", nil)     // "mm:ss"
	defaultNumFmts["46"] = NewNumFmt("46", "[h]:mm:ss", "", nil) // "[h]:mm:ss" What does [] mean?
	defaultNumFmts["47"] = NewNumFmt("47", "mmss.0", "%d", FormatInteger)
	defaultNumFmts["48"] = NewNumFmt("48", "##0.0E+0", "%e", FormatFloat)
	defaultNumFmts["49"] = NewNumFmt("49", "@", "", FormatDefault)

	// Special stuff
	defaultNumFmts["18"].format += " PM"
	defaultNumFmts["19"].format += " PM"

	for _, value := range defaultNumFmts {
		value.builtin = true
	}
}

// Styles represent Excel styles(.xml)
type Styles struct {
	cellXfs      []*cellXf
	cellStyleXfs []*cellStyleXf
	numFmts      map[string]*numFmt
}

func newStyles() (s *Styles) {
	s = new(Styles)

	s.numFmts = make(map[string]*numFmt)
	for _, nf := range defaultNumFmts {
		s.numFmts[nf.numFmtId] = nf
	}
	return s
}

func (s *Styles) AddCellStyleXf(xf *cellStyleXf) {
	// cellStyleXfs is 0-offset
	xf.index = len(s.cellStyleXfs)
	s.cellStyleXfs = append(s.cellStyleXfs, xf)
	// xf.index = len(s.cellXStylefs)
	xf.setUniqueID()
}

func (s *Styles) AddCellXf(xf *cellXf) {
	// cellXfs is 0-offset
	xf.index = len(s.cellXfs)
	s.cellXfs = append(s.cellXfs, xf)
	// xf.index = len(s.cellXfs)
	xf.setUniqueID()
}

func (s *Styles) AddNumFmt(nf *numFmt) {
	if _, found := s.numFmts[nf.numFmtId]; !found {
		s.numFmts[nf.numFmtId] = nf
	}
}

func (s *Styles) GetCellStyleXf(idx int) (xf *cellStyleXf) {
	return s.cellStyleXfs[idx]
}

func (s *Styles) GetCellXf(idx int) (xf *cellXf) {
	return s.cellXfs[idx]
}

func (s *Styles) GetNumFmt(id string) (nf *numFmt) {
	return s.numFmts[id]
}

func (s *Styles) GetCellStyleXfByFormat(format string) (xf *cellStyleXf) {
	for _, xf = range s.cellStyleXfs {
		if xf.nf.Code == format {
			return
		}
	}
	return nil
}

func (s *Styles) GetCellXfByFormat(format string) (xf *cellXf) {
	for _, xf = range s.cellXfs {
		if xf.nf.Code == format {
			return
		}
	}
	return nil
}

func (s *Styles) GetNumFmtByFormat(format string) (nf *numFmt) {
	for _, nf = range s.numFmts {
		if nf.Code == format {
			return
		}
	}
	return nil
}

func (s *Styles) importCellStyleXf(xf *cellStyleXf, nf *numFmt) (newXf *cellStyleXf) {
	newXf = new(cellStyleXf)
	newXf.FontId = xf.FontId
	newXf.FillId = xf.FillId
	newXf.BorderId = xf.BorderId
	newXf.XfId = xf.XfId
	newXf.ApplyNumberFormat = xf.ApplyNumberFormat
	newXf.ApplyFont = xf.ApplyFont
	newXf.ApplyFill = xf.ApplyFill
	newXf.ApplyBorder = xf.ApplyBorder
	newXf.ApplyProtection = xf.ApplyProtection
	newXf.QuotePrefix = xf.QuotePrefix
	newXf.nf = nf
	newXf.numFmtId = nf.numFmtId
	s.AddCellStyleXf(newXf)
	return
}

func (s *Styles) importCellXf(xf *cellXf, nf *numFmt) (newXf *cellXf) {
	newXf = new(cellXf)
	newXf.FontId = xf.FontId
	newXf.FillId = xf.FillId
	newXf.BorderId = xf.BorderId
	newXf.XfId = xf.XfId
	newXf.ApplyNumberFormat = xf.ApplyNumberFormat
	newXf.ApplyFont = xf.ApplyFont
	newXf.ApplyFill = xf.ApplyFill
	newXf.ApplyBorder = xf.ApplyBorder
	newXf.QuotePrefix = xf.QuotePrefix
	newXf.nf = nf
	newXf.numFmtId = nf.numFmtId
	s.AddCellXf(newXf)
	return
}

func (s *Styles) merge(src *Styles) {
	customId := customNumFmtID

	xfRemap := make(map[string]*cellXf)
	for _, xf := range s.cellXfs {
		xfRemap[xf.uniqueID] = xf
		if xf.nf.IsCustom() {
			customId++
		}
	}

	for _, xf := range src.cellXfs {
		_, found := xfRemap[xf.uniqueID]
		if found {
			continue
		}

		var numFmtId string
		if xf.nf.IsCustom() {
			numFmtId = strconv.Itoa(customId)
			customId++

			newNF := NewNumFmt(numFmtId, xf.nf.Code, xf.nf.format, xf.nf.formatter)
			s.numFmts[numFmtId] = newNF
		} else {
			numFmtId = xf.numFmtId
		}

		xfRemap[xf.uniqueID] = s.importCellXf(xf, s.numFmts[numFmtId])
	}

	styleXfRemap := make(map[string]*cellStyleXf)
	for _, xf := range s.cellStyleXfs {
		styleXfRemap[xf.uniqueID] = xf
		if xf.nf.IsCustom() {
			customId++
		}
	}

	for _, xf := range src.cellStyleXfs {
		_, found := styleXfRemap[xf.uniqueID]
		if found {
			continue
		}

		var numFmtId string
		if xf.nf.IsCustom() {
			numFmtId = strconv.Itoa(customId)
			customId++

			newNF := NewNumFmt(numFmtId, xf.nf.Code, xf.nf.format, xf.nf.formatter)
			s.numFmts[numFmtId] = newNF
		} else {
			numFmtId = xf.numFmtId
		}

		styleXfRemap[xf.uniqueID] = s.importCellStyleXf(xf, s.numFmts[numFmtId])
	}
}

func (s *Styles) open(file *File) (err error) {
	r, err := file.Open()
	if err != nil {
		return
	}

	saxer := new(saxStyles)
	saxer.target = s

	t := xml.NewTokenizer(saxer)
	if _, err = t.ReadFrom(bufio.NewReader(r)); err == io.EOF {
		err = nil
	}

	for _, xf := range s.cellStyleXfs {
		if xf.nf == nil {
			xf.nf = s.numFmts[xf.numFmtId]
		}
		xf.setUniqueID()
	}

	for _, xf := range s.cellXfs {
		if xf.nf == nil {
			xf.nf = s.numFmts[xf.numFmtId]
		}
		xf.setUniqueID()
	}

	return
}
