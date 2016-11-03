package types

import (
	"fmt"

	"github.com/agatan/gray/token"
)

// Error describes a type-checking error.
type Error struct {
	Message  string
	Pos      token.Position
	Filename string
}

// Error returns the error message.
func (e *Error) Error() string {
	return fmt.Sprintf("%s:%d:%d: %s", e.Filename, e.Pos.Line, e.Pos.Column, e.Message)
}
