package ast

import "github.com/agatan/gray/token"

// Stmt is an interface of statements
type Stmt interface {
	token.Pos
	stmt()
}

// StmtImpl provide default implementation for Stmt.
type StmtImpl struct {
	token.PosImpl
}

func (*StmtImpl) stmt() {}

type (
	// ExprStmt represent expression statements.
	ExprStmt struct {
		StmtImpl
		X Expr
	}

	// LetStmt represent let statements.
	LetStmt struct {
		StmtImpl
		Ident *Ident
		Type  Type // if nil, should be infered.
		Value Expr
	}
)
