package xml

import "errors"

var (
	// ErrIllegalCharacter returned when parser meets an illegal character.
	ErrIllegalCharacter = errors.New("Illegal character")
	// ErrUnterminatedTag returned when a '<' (tag start) is not balanced by a '>' (tag end).
	ErrUnterminatedTag = errors.New("Unterminated tag")
	// ErrUnterminatedString returned when a '"' (quote) is not balanced by a '""' (quote).
	ErrUnterminatedString = errors.New("Unterminated string")
	// ErrUnterminatedIdentifier returned when a namespace followed by a ':' (colon) is not followed by a proper identifier.
	ErrUnterminatedIdentifier = errors.New("Unterminated identifier")
	// ErrBadIdentifier returned when an invalid identifier is met.
	ErrBadIdentifier = errors.New("Bad identifier")
)
