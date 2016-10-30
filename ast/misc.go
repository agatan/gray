package ast

// Param represents function parameters.
type Param struct {
	Ident *Ident
	Type  Type
}
