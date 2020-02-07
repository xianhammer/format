package pdf

// var characters [256]tokenizer.Mask
type Mask = uint16

var characters [256]Mask

const (
	c_literalNameTerm Mask = 0x0001
	c_eol                  = 0x0002
	c_digit                = 0x0004
	// c_literalName           = 0x0008
	// c_delimiter       = 0x0010
	// c_open            = 0x0020
)

func init() {
	characters[0x00] = c_literalNameTerm
	characters[0x09] = c_literalNameTerm
	characters[0x0A] = c_literalNameTerm | c_eol
	characters[0x0C] = c_literalNameTerm | c_eol
	characters[0x0D] = c_literalNameTerm | c_eol
	characters[0x20] = c_literalNameTerm
	characters['/'] = c_literalNameTerm
	// characters[0x25] = c_delimiter          // %
	// characters[0x2F] = c_delimiter          // /
	// characters[0x28] = c_delimiter | c_open // (
	// characters[0x29] = c_delimiter          // )
	// characters[0x3C] = c_delimiter | c_open // <
	// characters[0x3E] = c_delimiter          // >
	// characters[0x5B] = c_delimiter | c_open // [
	// characters[0x5D] = c_delimiter          // ]
	// characters[0x7B] = c_delimiter | c_open // {
	// characters[0x7D] = c_delimiter          // }

	for digit := '0'; digit <= '9'; digit++ {
		characters[digit] |= c_digit
	}
}
