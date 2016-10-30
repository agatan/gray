package ast

import "github.com/agatan/gray/token"

// Type is an interface of types
type Type interface {
	token.Pos
	types()
}

// TypeImpl provides default implementations for Type.
type TypeImpl struct {
	token.PosImpl
}

func (*TypeImpl) types() {}

type (
	// TypeIdent represent identifier types.
	TypeIdent struct {
		TypeImpl
		Name string
	}

	// FuncType represent function types.
	FuncType struct {
		TypeImpl
		Params []*Param
		Result Type
	}
)
