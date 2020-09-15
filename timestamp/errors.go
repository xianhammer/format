package timestamp

import "errors"

var (
	// ErrUnknownTimestampFormat returned when the format is not recognised.
	ErrUnknownTimestampFormat = errors.New("Unknown timestamp format")
)
