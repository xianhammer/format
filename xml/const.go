package xml

const (
	csWhitespace int = 0x0001
	csLetter         = 0x0002
	csDigit          = 0x0004

	csIdentifierStart  = csLetter
	csIdentifierFollow = csIdentifierStart | csDigit
)

var charset [256]int

func init() {
	charset[' '] = csWhitespace
	charset['\r'] = csWhitespace
	charset['\n'] = csWhitespace
	charset['\t'] = csWhitespace

	charset['_'] = csLetter

	for letter := 'A'; letter <= 'Z'; letter++ {
		charset[letter] = csLetter
	}
	for letter := 'a'; letter <= 'z'; letter++ {
		charset[letter] = csLetter
	}
	for digit := '0'; digit <= '9'; digit++ {
		charset[digit] = csDigit
	}
}
