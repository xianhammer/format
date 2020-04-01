package parse

import (
	"bytes"
	"errors"
	"io"
	"math"
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

/*
var (
	win16be  = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	utf16bom = unicode.BOMOverride(win16be.NewDecoder())
)

// return transform.NewReader(s, utf16bom)*/

// TODO Proper handling of UTF8
// TODO Compare to:
// - https://github.com/buger/jsonparser
// - github.com/francoispqt/gojay
// - encoding/json
// ...
func JSON(b []byte, buffer []byte) (out interface{}, n int, err error) {
	array := make([]interface{}, 32)
	if buffer == nil {
		buffer = b[:]
	}
	return internalJSON(b, buffer, array)
}

func internalJSON(b []byte, buffer []byte, array []interface{}) (out interface{}, n int, err error) {
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

	// Number - 0 prefix is accepted in this parser. Using '1' in condiftion below ban the 0-prefix.
	if isdigit := (b[n] - '0') < 10; isdigit || (b[n]&0xF9) == 0x29 { // Second test is actually a bit too inclusive, weed out in body.
		n0, negNum := n, !isdigit && b[n] == '-'
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
				var n0 int
				for n++; n < l && (b[n]-'0') < 10; n++ {
					fp = (fp * 10.0) + float64(b[n]&0x0F)
					n0++
				}
				number += fp * math.Pow10(-n0)
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
		} else if n0 == n {
			return nil, n, ErrIllegalOperator
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
				v, n0 := Hex(b[n : n+4]) // TODO Faster convert?
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
		for n++; n < l && b[n] != ']'; {
			array[idx], n0, err = JSON(b[n:], buffer)
			n += n0
			if err != nil {
				if err == Terminal {
					break
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
		if idx > 0 {
			output = append(output, array[:idx]...)
		}
		return output, n + 1, nil
	}

	// Object
	if b[n] == '{' {
		obj := make(map[string]interface{})
		for n++; n < l && b[n] != '}'; {
			v, n0, err := JSON(b[n:], buffer)
			n += n0
			if err != nil {
				if err == Terminal {
					break
				}
				return nil, n, err
			}

			key, ok := v.(string)
			if !ok {
				return nil, n, ErrMissingString
			}

			for ; n < l && b[n] <= 32 && (b[n] == 0x09 || b[n] == 0x0A || b[n] == 0x0D || b[n] == 0x20); n++ {
			}
			if n < l && b[n] != ':' {
				return nil, l, ErrMissingColon
			}
			n++

			obj[key], n0, err = JSON(b[n:], buffer)
			n += n0
			if err != nil {
				delete(obj, key)
				return nil, n, err
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
		return obj, n + 1, nil
	}

	return nil, l, io.EOF
}

func JSONEqual(a, b interface{}) bool {
	switch v := a.(type) {
	case float64:
		return v == b.(float64)
	case string:
		return v == b.(string)
	case bool:
		return v == b.(bool)
	case []interface{}:
		w := b.([]interface{})
		if len(w) != len(v) {
			return false
		}
		for i, v0 := range v {
			if !JSONEqual(v0, w[i]) {
				return false
			}
		}
	case map[string]interface{}:
		w := b.(map[string]interface{})
		for key, v0 := range v {
			if !JSONEqual(v0, w[key]) {
				return false
			}
		}
	default:
		return b == nil
	}
	return true
}

/*
func JSONUnmarshal(source, target interface{}) (err error) {
	switch v := source.(type) {
	case float64:
		// return v == b.(float64)
	case string:
		// return v == b.(string)
	case bool:
		// return v == b.(bool)
	case []interface{}:
		// w := b.([]interface{})
		// if len(w) != len(v) {
		// 	return false
		// }
		// for i, v0 := range v {
		// 	if !JSONEqual(v0, w[i]) {
		// 		return false
		// 	}
		// }
	case map[string]interface{}:
		// w := b.(map[string]interface{})
		// for key, v0 := range v {
		// 	if !JSONEqual(v0, w[key]) {
		// 		return false
		// 	}
		// }
	default:
		// return b == nil
	}
	// return true
}
*/
