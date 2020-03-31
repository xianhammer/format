package parse

import (
	"errors"
	"fmt"
	"io"
	"math"
)

// Spec: https://www.json.org/json-en.html

var escaped [256]byte

func init() {
	escaped['"'] = '"'
	escaped['\\'] = '\\'
	escaped['b'] = '\b'
	escaped['f'] = '\f'
	escaped['n'] = '\n'
	escaped['r'] = '\r'
	escaped['t'] = '\t'
}

const Comma = int(',')      // (golang) ints cannot be returned from JSON so use this as a marker.
const Colon = int(':')      // (golang) ints cannot be returned from JSON so use this as a marker.
const TermArray = int(']')  // (golang) ints cannot be returned from JSON so use this as a marker.
const TermObject = int('}') // (golang) ints cannot be returned from JSON so use this as a marker.

var ErrMissingComma = errors.New("Expected comma")
var ErrMissingColon = errors.New("Expected colon")
var ErrMissingString = errors.New("Expected string (key)")

/*
var (
	win16be  = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	utf16bom = unicode.BOMOverride(win16be.NewDecoder())
)

// return transform.NewReader(s, utf16bom)*/

// TODO Proper handling of UTF8
// TODO Compare to:
// - https://github.com/buger/jsonparser
// - github.com/francoispqt/gojay
// - encoding/json
// ...
func JSON(b []byte, buffer []byte) (out interface{}, n int) {
	l := len(b)

	// Whitespaces
	for ; n < l && b[n] <= 32 && (b[n] == 0x09 || b[n] == 0x0A || b[n] == 0x0D || b[n] == 0x20); n++ {
		// out = byte(0x20) // If n>=l condition below is true, a proper return values is needed. Here char(32) is returned.
	}

	if n >= l {
		return io.EOF, n
	}

	switch b[n] {
	case ',':
		return Comma, n + 1
	case ':':
		return Colon, n + 1
	case '}':
		return TermObject, n + 1
	case ']':
		return TermArray, n + 1
	}

	// Number
	if (b[n] - '0') < 10 {
		v, n0 := Decimal(b[n:])
		n += n0

		var number = float64(v)
		if n >= l {
			return number, n
		}

		if n+1 < l && b[n] == '.' {
			fp, n0 := Decimal(b[n+1:])
			n += n0 + 1
			number += float64(fp) * math.Pow10(-n0)
		}

		if n < l && (b[n]&0xDF) == 'E' {
			n++
			negative := b[n] == '-'
			if negative || b[n] == '+' {
				n++
			}

			exp0, n0 := Decimal(b[n:])
			n += n0

			if negative {
				number *= math.Pow10(-int(exp0))
			} else {
				number *= math.Pow10(int(exp0))
			}
		}

		return number, n
	}

	// Keywords: true, false, null
	if n+4 <= l {
		if b[n] == 't' && b[n+1] == 'r' && b[n+2] == 'u' && b[n+3] == 'e' { // true
			return true, n + 4
		}
		if b[n] == 'n' && b[n+1] == 'u' && b[n+2] == 'l' && b[n+3] == 'l' { // null
			return nil, n + 4
		}
		if n+5 <= l && b[n] == 'f' && b[n+1] == 'a' && b[n+2] == 'l' && b[n+3] == 's' && b[n+4] == 'e' { // false
			return false, n + 5
		}
	}

	if n+1 > l { // string, array and objects require at least two characters...
		return io.EOF, n
	}

	// String
	if b[n] == '"' {
		n++
		if buffer == nil {
			buffer = b[n:]
		}

		i := 0
		for ; n < l && b[n] != '"'; n++ {
			if b[n] == '\\' {
				if n++; b[n] == 'u' {
					v, n0 := Hex(b[n : n+4])
					n += n0
					buffer[i] = byte(v)
				} else {
					buffer[i] = escaped[b[n]]
				}
			} else {
				buffer[i] = b[n]
			}
			i++
		}
		return string(buffer[:i]), n + 1
	}

	// Array
	if b[n] == '[' {
		array := make([]interface{}, 0)
		for n++; n < l; {
			v, n0 := JSON(b[n:], buffer)
			n += n0
			if v == TermArray {
				break
			}
			if v != Comma {
				array = append(array, v)
			}
		}
		return array, n
	}

	// Object
	if b[n] == '{' {
		obj := make(map[string]interface{})
		for n++; n < l; {
			v, n0 := JSON(b[n:], buffer)
			n += n0
			if v == TermObject {
				break
			}
			if v == Comma {
				v, n0 = JSON(b[n:], buffer)
				n += n0
			}

			key, ok := v.(string)
			if !ok {
				fmt.Printf("v = %v\n", v)
				return ErrMissingString, n
			}

			v, n0 = JSON(b[n:], buffer)
			n += n0
			if v != Colon {
				return ErrMissingColon, n // Error state, actually
			}

			obj[key], n0 = JSON(b[n:], buffer)
			n += n0
		}
		return obj, n
	}

	return io.EOF, n
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
