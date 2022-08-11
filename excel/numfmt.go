package excel

import (
	"fmt"
	"regexp"
	"time"

	"github.com/xianhammer/format/parse"
	"github.com/xianhammer/format/xml"
)

const customNumFmtID = 164

func FormatDefault(f *numFmt, data []byte) (out string) {
	return string(data)
}

func FormatStandard(f *numFmt, data []byte) (out string) {
	return fmt.Sprintf(f.format, data)
}

func FormatInteger(f *numFmt, data []byte) (out string) {
	val, _ := parse.Decimal(data)
	return fmt.Sprintf(f.format, val)
}

func FormatFloat(f *numFmt, data []byte) (out string) {
	val, _ := parse.Float(data)
	return fmt.Sprintf(f.format, val)
}

func FormatDatetime(f *numFmt, data []byte) (out string) {
	val, _ := parse.Float(data)
	v := time.Unix(int64(float64(excel1900Epoc)+86400.0*val), 0).UTC()
	return v.Format(f.format)
}

type numFmt struct {
	numFmtId  string
	Code      string
	format    string
	builtin   bool
	formatter func(f *numFmt, data []byte) (out string)
}

func NewNumFmt(numFmtId, code, goFormat string, formatter func(f *numFmt, data []byte) (out string)) (f *numFmt) {
	f = new(numFmt)
	f.numFmtId = numFmtId
	f.format = goFormat
	f.formatter = formatter

	if formatter == nil {
		f.SetDatetime(code)
	} else {
		f.SetCode(code)
	}
	return
}

func (f *numFmt) IsCustom() (custom bool) {
	return f.builtin == false
}

func (f *numFmt) SetCode(code string) {
	// TODO parse the given code...
	// f.SetDatetime(code)
	// f.Code = strings.ReplaceAll(code, "\\", "")
	f.Code = code
}

func (f *numFmt) SetDatetime(code string) {
	repl := formatreplacer()
	f.SetCode(code)
	f.format = rFormatCode.ReplaceAllStringFunc(f.Code, repl)
	f.formatter = FormatDatetime
}

func (f *numFmt) toXMLBuilder(b *xml.Builder) {
	b.Tag([]byte("numFmt"))
	b.Attr([]byte("numFmtId"), []byte(f.numFmtId))
	b.Attr([]byte("formatCode"), []byte(f.Code))
	b.EndTag() // End <numFmt>
}

var rFormatCode = regexp.MustCompile(`TZName|Z:Z|ZZ|dd|MM|yyyy|yy|hh|mm|ss|[\[\]ZdMyhms]`) //+[0#?.,*_@]
var formatreplacer = func() (replacer func(in string) (out string)) {
	// So, stupid M$ has at least TWO distinct formats (for eg. dates): "dd\-mm\-yyyy\ hh:mm:ss" and "dd-MM-yyyy hh:mm:ss"
	// This, of course, makes parsing so much harder - grrrrrrr....
	previous := ""
	return func(in string) (out string) {
		defer func() { previous = in }()
		if len(in) > 0 {
			switch in[0] {
			case 'T':
				if in == "TZName" {
					return "MST"
				}
			case 'd':
				return "02"[2-len(in):]
			case 'M':
				return "01"[2-len(in):]
			case 'y':
				return "2006"[4-len(in):]
			case 'h':
				return "15"
			case 'm':
				if previous == "dd" && in == "mm" { // May need to handle mm-dd m-d mm-d m-dd d-mm dd-m, .... SIGH!
					return "01"
				}
				return "04"
			case 's':
				return "05"
			case 'Z':
				if len(in) == 1 {
					return "Z07"
				}
				if len(in) == 2 {
					return "Z0700"
				}
				if len(in) == 3 && in[1] == ':' {
					return "Z07:00"
				}
			case '[', ']':
				return ""
			}
		}
		return in
	}
}
