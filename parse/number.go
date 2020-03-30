package parse

import (
	"math"
)

// Hex parses all hexadecimal digits from start of the given byte slice.
// If a non-hexdigit is met, parsing stops.
// Return n - number of bytes read, and i - the value found.
func Hex(b []byte) (i uint64, n int) {
	for l := len(b); n < l; n++ {
		if (b[n] - '0') < 10 {
			i = (i * 16) + uint64(b[n]&0x0F)
		} else if (b[n]&0xDF)-'A' < 6 {
			i = (i * 16) + uint64(b[n]&0x0F) + 9
		} else {
			break
		}
	}
	return
}

// Decimal parses all decimal digits from start of the given byte slice.
// If a non-digit is met, parsing stops.
// Return n - number of bytes read, and i - the value found.
func Decimal(b []byte) (i uint64, n int) {
	l := len(b)
	for ; n < l && (b[n]-'0') < 10; n++ {
		i = (i * 10) + uint64(b[n]&0x0F)
	}
	return
}

// Float parses a floating point value from start of the given byte slice.
// If a non-compliant byte is met, parsing stops.
// Return n - number of bytes read, and f - the value found.
func Float(b []byte) (f float64, n int) {
	var ip, fp, e int64
	l := len(b)
	for ; n < l && (b[n]-'0') < 10; n++ {
		ip = 10*ip + int64(b[n]&0x0f)
	}

	e, f = 1, float64(ip)
	if n < l && b[n] == '.' {
		for n++; n < l && (b[n]-'0') < 10; n++ {
			fp = 10*fp + int64(b[n]&0x0f)
			e *= 10
		}
		f += float64(fp) / float64(e)
	}

	if n >= l || (b[n]&0xDF) != 'E' {
		return
	}

	n++
	exp, sign := 0, 1
	if n < l && b[n] == '-' {
		sign = -1
		n++
	}

	for ; n < l && (b[n]-'0') < 10; n++ {
		exp = 10*exp + int(b[n]&0x0f)
	}

	if exp != 0 {
		f *= math.Pow10(sign * exp)
	}

	return
}

// Number parses "digits" (0-9 and A-..., depending on base) from start of the given byte slice.
// Parsing stop if a non-"digit" is met.
// Use parse.Hex to parse known hexadecimals (is slightly faster)
// Return n - number of bytes read, and i - the value found.
func Number(b []byte, base uint64) (i uint64, n int) {
	for bs, l := byte(base-10), len(b); n < l; n++ {
		if (b[n] - '0') < 10 {
			i = (i * base) + uint64(b[n]&0x0F)
		} else if (b[n]&0xDF)-'A' < bs {
			i = (i * base) + uint64(b[n]&0x0F) + 9
		} else {
			break
		}
	}
	return
}
