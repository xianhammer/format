package json

import (
	"io"
	"math"
)

type Kind byte

const (
	True  Kind = 't'
	False      = 'f'
	Null       = 'n'
)

func Parse(input io.Reader, emit SAX, buffer []byte) (err error) {
	if buffer == nil {
		buffer = make([]byte, 10*1024) // Output buffer for strings and keys
	}

	p := 0
	level := 0
	isNegative := false

	var terminal [256]bool // Terminals for tokens: true, false and null.
	terminal[','] = true
	terminal[']'] = true
	terminal['}'] = true
	terminal[0x20] = true
	terminal[0x0A] = true
	terminal[0x0D] = true
	terminal[0x09] = true

	truePostfix := []byte("rue")
	falsePostfix := []byte("alse")
	nullPostfix := []byte("ull")

	pushchar := false
	next := func() {
		if pushchar {
			pushchar = false
		} else if err == nil {
			_, err = input.Read(buffer[p : p+1])
		}
	}

	expect := func(e []byte) {
		t := buffer[p]
		if next(); err != nil || buffer[p] != e[0] {
			err = ErrUnexpectedInput
			return
		}
		if next(); err != nil || buffer[p] != e[1] {
			err = ErrUnexpectedInput
			return
		}
		if next(); err != nil || buffer[p] != e[2] {
			err = ErrUnexpectedInput
			return
		}

		if len(e) > 3 {
			if next(); err != nil || buffer[p] != e[3] {
				err = ErrUnexpectedInput
				return
			}
		}

		// fmt.Printf("buffer[%d]=%c\n", p, buffer[p])
		if next(); err == io.EOF || (err == nil && terminal[buffer[p]]) {
			// fmt.Printf("--> buffer[%d]=%c\n", p, buffer[p])
			pushchar = true
			emit.Literal(Kind(t))
		} else {
			err = ErrUnexpectedInput
		}
	}

	for next(); err == nil; next() {
		switch buffer[p] {
		case ':':
		case ',':
		case '[':
			level++
			emit.Array()

		case ']':
			level--
			emit.ArrayEnd()

		case '{':
			level++
			emit.Object()

		case '}':
			level--
			emit.ObjectEnd()

		case '"', '\'':
			end, escaped := buffer[p], false
			p := 0 // Use full extent of buffer
			for next(); err == nil && !(escaped && buffer[p] == end); next() {
				if buffer[p] == '\\' {
					escaped = !escaped
				}

				p++
				if p >= len(buffer) {
					emit.String(buffer[:p])
					p = 0
				}
			}

			if p > 0 {
				emit.String(buffer[:p])
			}

			emit.StringEnd()
			p = 0 // Reset p

		case 't':
			expect(truePostfix)

		case 'f':
			expect(falsePostfix)

		case 'n':
			expect(nullPostfix)

		case 0x20, 0x0A, 0x0D, 0x09: // Whitespace - ignore

		case '-':
			if next(); err != nil || buffer[p]-'0' >= 10 {
				return ErrUnexpectedInput
			}

			isNegative = true
			fallthrough

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			ip := int64(buffer[p] - '0')
			for next(); err == nil && buffer[p]-'0' < 10; next() {
				ip = 10*ip + int64(buffer[p]&0x0f)
			}

			isFloat := false
			fv := float64(ip)
			if buffer[p] == '.' {
				fp, exp := int64(0), 1
				for next(); err == nil && buffer[p]-'0' < 10; next() {
					fp = 10*fp + int64(buffer[p]&0x0f)
					exp *= 10
				}
				isFloat = fp != 0
				fv += float64(fp) / float64(exp)
			}

			exp := 0
			if (buffer[p] & 0xdf) == 'E' {
				sign := 1
				next()
				if err == nil && (buffer[p] == '-' || buffer[p] == '+') {
					if buffer[p] == '-' {
						sign = -1
					}
					next()
				}

				for ; err == nil && buffer[p]-'0' < 10; next() {
					exp = 10*exp + int(buffer[p]&0x0f)
				}

				isFloat = sign < 0 && exp != 0
				exp *= sign
			}

			// fmt.Printf("[fv=%v] [float=%v] [exp=%v] [buffer=%v] [err=%v]\n", fv, isFloat, exp, buffer[p], err)
			if fv == 0.0 {
				emit.Integer(0)
			} else if isFloat {
				if exp != 0 {
					fv *= math.Pow10(exp)
				}

				if isNegative {
					emit.Float(-fv)
				} else {
					emit.Float(fv)
				}
			} else {
				if exp != 0 {
					ip *= int64(math.Pow10(exp))
				}
				if isNegative {
					emit.Integer(-ip)
				} else {
					emit.Integer(ip)
				}
			}
			isNegative = false

		default:
			// fmt.Printf("default: p=%d, b=%v\n", p, buffer[:64])
			return ErrUnexpectedInput
		}
		// fmt.Printf("error = %v\n", err)
	}

	if err == io.EOF {
		err = nil
	}

	return
}
