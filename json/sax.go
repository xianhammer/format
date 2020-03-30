package json

import (
	"fmt"
)

type SAX interface {
	Array()
	ArrayEnd()
	Object()
	ObjectEnd()
	String(part []byte)
	StringEnd()

	Literal(t Kind)
	Integer(v int64)
	Float(v float64)
}

var ErrUnexpectedInput = fmt.Errorf("Unexpected input")
