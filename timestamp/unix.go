package timestamp

import (
	"syscall"
	"time"
)

type TimeConvert func(b []byte) (t time.Time)

// FromFiletime convert a Windows Filetime to Go time.
// Filetime is expected to be in little endian and so is the input.
func FromFiletime(value int64) (t time.Time) {
	t := &syscall.Filetime{
		LowDateTime:  b >> 16,
		HighDateTime: b & 0xFFFF,
	}
	return time.Unix(0, t.Nanoseconds())
}

func FromUnix(value int64) (t time.Time) {
	return time.Unix(value, 0)
}

func FromUnixNsec(value, nsec int64) (t time.Time) {
	return time.Unix(value, nsec)
}
