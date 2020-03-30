package json

type testsax struct {
	intValue   int64
	floatValue float64

	boolValue   bool
	nullValue   bool
	stringValue []byte

	boolSet, nullSet bool
	intSet, floatSet bool
	stringSet        bool

	str, stringEnd    int
	array, arrayEnd   int
	object, objectEnd int
}

func (s *testsax) Array()     { s.array++ }
func (s *testsax) ArrayEnd()  { s.arrayEnd++ }
func (s *testsax) Object()    { s.object++ }
func (s *testsax) ObjectEnd() { s.objectEnd++ }

func (s *testsax) Literal(t Kind) {
	s.boolSet = t == 't' || t == 'f'
	s.boolValue = t == 't'
	s.nullSet = t == 'n'
	s.nullValue = t == 'n'
}
func (s *testsax) Integer(v int64) { s.intSet = true; s.intValue = v }
func (s *testsax) Float(v float64) { s.floatSet = true; s.floatValue = v }
func (s *testsax) StringEnd()      { s.stringEnd++ }

func (s *testsax) String(part []byte) {
	if !s.stringSet {
		s.str++
	}
	s.stringSet = true
	s.stringValue = append(s.stringValue, part...)
}

func accept(s *testsax) bool { return true }

func isTrue(s *testsax) bool    { return s.boolSet && s.boolValue == true }
func isFalse(s *testsax) bool   { return s.boolSet && s.boolValue == false }
func isNull(s *testsax) bool    { return s.nullSet && s.nullValue == true }
func isInteger(s *testsax) bool { return s.intSet }
func isFloat(s *testsax) bool   { return s.floatSet }
func isString(s *testsax) bool  { return s.stringSet }

func isNumberInteger(v int64) func(s *testsax) bool {
	return func(s *testsax) bool { return isInteger(s) && v == s.intValue }
}
func isNumberFloat(v float64) func(s *testsax) bool {
	return func(s *testsax) bool { return isFloat(s) && v == s.floatValue }
}
