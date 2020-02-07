package excel

func ParseReference(value []byte) (column, row int, err error) {
	i, l := 0, len(value)
	for ; i < l && (value[i]-'A') < azBase; i++ {
		column = azBase*column + int(value[i]&0x1f)
	}

	for ; i < l && (value[i]-'0') < 10; i++ {
		row = 10*row + int(value[i]&0x0f)
	}

	if column <= 0 || row <= 0 {
		err = ErrInvalidReference
	}
	return
}
