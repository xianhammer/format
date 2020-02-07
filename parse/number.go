package parse

import (
	"math"
)

// Integer parses all decimal digits from start of the given byte slice.
// If a non-digit is met, parsing stops.
// Return n - number of bytes read, and i - the value found.
func Integer(b []byte) (i int64, n int) {
	l := len(b)
	for ; n < l && (b[n]-'0') < 10; n++ {
		i = 10*i + int64(b[n]&0x0f)
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

	if n >= l || (b[n]&0xdf) != 'E' {
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
