package excel

import "fmt"

type Row []Cell

func (r Row) Cell(idx int) (c *Cell) {
	return &r[idx]
}

func (r Row) Value(ss *SharedStrings, applyStyle bool) (row []string) {
	row = make([]string, len(r))
	for i := 0; i < len(r); i++ {
		row[i] = r[i].Value(ss, applyStyle)
	}
	return
}

func (r Row) Pick(pick []int, ss *SharedStrings, applyStyle bool) (row []string) {
	row = make([]string, len(pick))
	for i := 0; i < len(pick); i++ {
		row[i] = r[pick[i]].Value(ss, applyStyle)
	}
	return
}

func (r Row) Indeces(pick []string, ss *SharedStrings, applyStyle bool, translate map[string]string) (indeces []int, err error) {
	p := make(map[string]int)
	for i := 0; i < len(pick); i++ {
		p[pick[i]] = i
	}

	if len(pick) != len(p) {
		err = fmt.Errorf("Multiple defined picks")
		return
	}

	added := 0
	indeces = make([]int, len(pick))
	for i := 0; i < len(r); i++ {
		k := r[i].Value(ss, applyStyle)
		if translate != nil {
			if newkey, found := translate[k]; found {
				k = newkey
			}
		}
		if idx, found := p[k]; found {
			indeces[idx] = i
			added++
		}
	}

	if added != len(pick) {
		err = fmt.Errorf("Unknown index pick")
	}

	return
}
