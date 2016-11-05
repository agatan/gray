package ast

import "github.com/agatan/gray/token"

// Type is an interface of types
type Type interface {
	token.Pos
	Node
	types()
}

// TypeImpl provides default implementations for Type.
type TypeImpl struct {
	token.PosImpl
	NodeImpl
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

	// InstType represent instantiation types.
	InstType struct {
		TypeImpl
		Base Type
		Args []Type
	}
)
