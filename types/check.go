package types

// Checker contains a type checking state.
type Checker struct {
	Filename string
	scope    *Scope
}

// NewChecker creates a Checker with given file name.
func NewChecker(filename string) *Checker {
	return &Checker{Filename: filename, scope: NewScope(filename)}
}
