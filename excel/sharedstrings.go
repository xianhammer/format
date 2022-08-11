package excel

import (
	"strconv"
)

type SharedStrings struct {
	strings []string
	indeces map[string]int
}

func newSharedStrings() (s *SharedStrings) {
	s = new(SharedStrings)
	return s
}

func (s *SharedStrings) open(file *File) (err error) {
	r, err := file.Open()
	if err == nil {
		s.strings, err = NodeContent(r, []byte("t"))
	}
	return
}

func (s *SharedStrings) merge(src *SharedStrings) (err error) {
	for _, str := range src.strings {
		s.add(str)
	}
	return
}

func (s *SharedStrings) addIdx(value string) (idx int) {
	if s.indeces == nil {
		s.indeces = make(map[string]int)
		for i, v := range s.strings {
			s.indeces[v] = i
		}
	}

	idx, ok := s.indeces[value]
	if !ok {
		idx = len(s.strings)
		s.indeces[value] = idx
		s.strings = append(s.strings, value)
	}

	return
}

func (s *SharedStrings) add(value string) (v string) {
	return strconv.Itoa(s.addIdx(value))
}

func (s *SharedStrings) Get(idx int) (v string) {
	return s.strings[idx]
}

// GetIdx return 1-offset shared strings index. If 0 is returned no string matched.
func (s *SharedStrings) GetIdx(v string) (idx int) {
	for i, l := 0, len(s.strings); i < l; i++ {
		if s.strings[i] == v {
			return i
		}
	}
	return
}

func (s *SharedStrings) GetFromValue(value []byte) (v string) {
	idx := 0
	for i, l := 0, len(value); i < l && '0' <= value[i] && value[i] <= '9'; i++ {
		idx = 10*idx + int(value[i]&0x0f)
	}
	return s.strings[idx]
}
