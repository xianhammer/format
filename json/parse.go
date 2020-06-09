package json

import (
	"bytes"
	"errors"
	"io"
	"math"

	"github.com/xianhammer/format/parse"
)

// Spec: https://www.json.org/json-en.html

var (
	escaped [256]byte

	keywordTrue  = []byte("true")
	keywordFalse = []byte("false")
	keywordNull  = []byte("null")

	Terminal = errors.New("Array or object terminal")

	// ErrMissingComma    = errors.New("Expected comma")
	ErrMissingColon    = errors.New("Expected colon")
	ErrMissingString   = errors.New("Expected string (key)")
	ErrIllegalOperator = errors.New("Illegal operator")
)

func init() {
	escaped['"'] = '"'
	escaped['\\'] = '\\'
	escaped['b'] = '\b'
	escaped['f'] = '\f'
	escaped['n'] = '\n'
	escaped['r'] = '\r'
	escaped['t'] = '\t'
}

// TODO Proper handling of UTF8
/*
var (
	win16be  = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	utf16bom = unicode.BOMOverride(win16be.NewDecoder())
)
... transform.NewReader(s, utf16bom)
*/

func Parse(b []byte, buffer []byte) (out interface{}, n int, err error) {
	if buffer == nil {
		buffer = b[:]
	}
	return internalParse(b, buffer)
}

func internalParse(b []byte, buffer []byte) (out interface{}, n int, err error) {
	l := len(b)

	// Whitespaces
	switch b[n] {
	case 0x09, 0x0A, 0x0D, 0x20:
		for n++; n < l; n++ {
			if b[n] > 0x20 || (b[n] != 0x09 && b[n] != 0x0A && b[n] != 0x0D && b[n] != 0x20) {
				break
			}
		}

		if n >= l { // Fast "fail"
			return nil, n, nil
		}
	}

	switch b[n] {
	case '}', ']':
		return nil, n, Terminal
	case 't':
		if bytes.Equal(b[n:n+4], keywordTrue) {
			return true, n + 4, nil
		}
	case 'n':
		if bytes.Equal(b[n:n+4], keywordNull) {
			return nil, n + 4, nil
		}
	case 'f':
		if bytes.Equal(b[n:n+5], keywordFalse) {
			return false, n + 5, nil
		}
	}

	// Numbers prefixed with 0 is accepted in this parser. Using '1' in condidtion below prevent this.
	if isdigit := (b[n] - '0') < 10; isdigit || (b[n]&0xF9) == 0x29 { // Second test is actually a bit too inclusive, weed out in body.
		negNum := !isdigit && b[n] == '-'
		if negNum || b[n] == '+' {
			n++
		}

		var number float64
		for ; n < l && (b[n]-'0') < 10; n++ {
			number = (number * 10.0) + float64(b[n]&0x0F)
		}

		if n < l {
			if b[n] == '.' {
				var fp float64
				var pwr int
				for n++; n < l && (b[n]-'0') < 10; n++ {
					fp = (fp * 10.0) + float64(b[n]&0x0F)
					pwr++
				}
				number += fp * math.Pow10(-pwr)
			}

			if n < l && (b[n]&0xDF) == 'E' {
				n++
				negExp := b[n] == '-'
				if negExp || b[n] == '+' {
					n++
				}

				var exp0 int
				for ; n < l && (b[n]-'0') < 10; n++ {
					exp0 = (exp0 * 10) + int(b[n]&0x0F)
				}

				if negExp {
					number *= math.Pow10(-exp0)
				} else {
					number *= math.Pow10(exp0)
				}
			}
		}

		if negNum {
			number = -number
		}
		return number, n, nil
	}

	// String
	if b[n] == '"' {
		i := 0
		for n++; n < l && b[n] != '"'; n++ {
			if b[n] != '\\' {
				buffer[i] = b[n]
			} else if n++; b[n] == 'u' {
				v, n0 := parse.Hex(b[n : n+4]) // TODO Faster convert?
				n += n0
				buffer[i] = byte(v) // TODO rune(v)
			} else {
				buffer[i] = escaped[b[n]]
			}
			i++
		}
		return string(buffer[:i]), n + 1, nil
	}

	// Array
	if b[n] == '[' {
		var output []interface{}
		var n0 int
		idx := 0
		array := make([]interface{}, 32)
		for n++; n < l && b[n] != ']'; {
			array[idx], n0, err = internalParse(b[n:], buffer)

			n += n0
			if err != nil {
				if err == Terminal {
					break
				}
				if err == io.EOF {
					err = Terminal
				}
				return nil, n, err
			}

			if idx++; idx == len(array) {
				output = append(output, array...)
				idx = 0
			}

			for ; n < l && b[n] <= 32 && (b[n] == 0x09 || b[n] == 0x0A || b[n] == 0x0D || b[n] == 0x20); n++ {
			}

			if n < l && b[n] != ']' {
				n++
			}
		}
		if n >= l || b[n] != ']' {
			return nil, n, Terminal
		}
		if idx > 0 {
			output = append(output, array[:idx]...)
		}
		return output, n + 1, nil
	}

	// Object
	if b[n] == '{' {
		obj := make(map[string]interface{})
		for n++; n < l && b[n] != '}'; {
			v, n0, err := internalParse(b[n:], buffer)
			n += n0
			if err != nil {
				if err == Terminal {
					break
				}
				if err == io.EOF {
					err = ErrMissingString
				}
				return nil, n, err
			}

			key, ok := v.(string)
			if !ok || n >= l {
				return nil, n, ErrMissingString
			}

			for ; n < l && b[n] <= 32 && (b[n] == 0x09 || b[n] == 0x0A || b[n] == 0x0D || b[n] == 0x20); n++ {
			}
			if n+1 >= l || b[n] != ':' {
				return nil, l, ErrMissingColon
			}
			n++

			obj[key], n0, err = internalParse(b[n:], buffer)
			n += n0
			if err != nil {
				if err == io.EOF {
					err = Terminal
				}
				delete(obj, key)
				return nil, n, err
			}

			if n >= l {
				delete(obj, key)
				return nil, l, Terminal
			}

			for ; n < l && b[n] <= 32 && (b[n] == 0x09 || b[n] == 0x0A || b[n] == 0x0D || b[n] == 0x20); n++ {
			}

			if n < l && b[n] != '}' {
				if b[n] != ',' {
					return nil, l, Terminal
				}
				n++
			}
		}

		if n >= l || b[n] != '}' {
			return nil, n, Terminal
		}
		return obj, n + 1, nil
	}

	return nil, l, io.EOF
}
