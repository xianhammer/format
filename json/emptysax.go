package json

type EmptySAX struct{}

func (s *EmptySAX) Array()             {}
func (s *EmptySAX) ArrayEnd()          {}
func (s *EmptySAX) Object()            {}
func (s *EmptySAX) ObjectEnd()         {}
func (s *EmptySAX) Literal(t Kind)     {}
func (s *EmptySAX) Integer(v int64)    {}
func (s *EmptySAX) Float(v float64)    {}
func (s *EmptySAX) String(part []byte) {}
func (s *EmptySAX) StringEnd()         {}
