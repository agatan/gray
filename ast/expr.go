package ast

import "github.com/agatan/gray/token"

// Expr is an interface of expressions
type Expr interface {
	token.Pos
	expr()
}

// ExprImpl provides default implementations for Expr.
type ExprImpl struct {
	token.PosImpl
}

func (*ExprImpl) expr() {}

type (
	// Ident represent identifier expressions.
	Ident struct {
		ExprImpl
		Name string
	}

	// BasicLit represent literal nodes.
	BasicLit struct {
		ExprImpl
		Kind token.Token // Kind should be token.UNIT, token.BOOL or token.INT
		Lit  string
	}

	// BlockExpr represent block expression.
	BlockExpr struct {
		ExprImpl
		Stmts []Stmt // If empty, the value is Unit.
	}
)
