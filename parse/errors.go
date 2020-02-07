package parse

import "errors"

var (
	// ErrOutOfBounds returned when buffer run out of space.
	ErrOutOfBounds = errors.New("Buffer index out of bounds")
)
