package parse

import (
	"io"
	"math"
)

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

func JSON(b []byte, buffer []byte) (out interface{}, n int) {
	l := len(b)

	// Whitespaces
	for ; n < l && b[n] <= 32 && (b[n] == 0x09 || b[n] == 0x0A || b[n] == 0x0D || b[n] == 0x20); n++ {
		// out = byte(0x20) // If n>=l condition below is true, a proper return values is needed. Here char(32) is returned.
	}

	if n >= l {
		return io.EOF, n
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

	if n+2 > l {
		return io.EOF, n + 2 // Both string, array and objects require at least two characters...
	}

	// String
	// Returned values is UN-INTERPRETED - that is, any escaped characters are still the backslash followed by the char escaped.
	// TODO Proper handling of UTF8
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
		for n++; n < l && b[n] != ']'; {
			v, n0 := JSON(b[n:], buffer)
			if n0 != 0 && v != nil {
				array = append(array, v)
			}

			n += n0
			if b[n] == ',' {
				n++
			}
		}
		return array, n + 1
	}

	if b[n] == '{' {
		obj := make(map[string]interface{})
		for n++; n < l && b[n] != '}'; {
			vKey, n0 := JSON(b[n:], buffer)
			n += n0

			key, ok := vKey.(string)
			if !ok && b[n] == ',' {
				n++
				vKey, n0 = JSON(b[n:], buffer)
				n += n0
				key, ok = vKey.(string)
			}

			if !ok {
				return obj, n
			}

			_, n0 = JSON(b[n:], buffer)
			n += n0
			if b[n] != ':' {
				return obj, n
			}

			n++
			obj[key], n0 = JSON(b[n:], buffer)
			n += n0
		}
		return obj, n + 1
	}
	// Object
	return
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
