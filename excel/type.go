package excel

type Type byte

const (
	Boolean Type = 'b'
	Date         = 'd'
	Error        = 'e'
	Inline       = 'i'
	Number       = 'n'
	String       = 's'
	Formula      = 'f'
)
