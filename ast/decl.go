package ast

import "github.com/agatan/gray/token"

// Decl is an interface of declarations
type Decl interface {
	token.Pos
	Node
	decl()
}

// DeclImpl provide default implementation for Decl.
type DeclImpl struct {
	token.PosImpl
	NodeImpl
}

func (*DeclImpl) decl() {}

type (
	FuncDecl struct {
		DeclImpl
		Ident *Ident
		Type  *FuncType
		Body  *BlockExpr
	}
)
