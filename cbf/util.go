package cbf

import (
	"io"
	"sort"
	"unicode/utf16"
)

func UTF16String(b []uint16) (d string) {
	return string(utf16.Decode(b))
}

func FilePos(r io.Reader) (pos int64) {
	s := r.(io.Seeker)
	pos, _ = s.Seek(0, io.SeekCurrent)
	return
}

func SortDirectories(dir []*DirectoryEntry) {
	sort.Slice(dir, func(i, j int) bool {
		return dir[i].Name() < dir[j].Name()
	})
}
