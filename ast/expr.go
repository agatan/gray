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
		Kind token.Kind // Kind should be token.UNIT, token.BOOL or token.INT
		Lit  string
	}

	// RefExpr represent creating reference expression.
	RefExpr struct {
		ExprImpl
		Value Expr
	}

	// DerefExpr represent dereference expression.
	DerefExpr struct {
		ExprImpl
		Ref Expr
	}

	// InfixExpr represent binary operation expression.
	InfixExpr struct {
		ExprImpl
		LHS      Expr
		Operator string
		RHS      Expr
	}

	// ParenExpr represent parensed expressions .
	ParenExpr struct {
		ExprImpl
		X Expr
	}

	// BlockExpr represent block expression.
	BlockExpr struct {
		ExprImpl
		Stmts []Stmt // If empty, the value is Unit.
	}

	// CallExpr represent calling expression.
	CallExpr struct {
		ExprImpl
		Func Expr
		Args []Expr
	}

	// IfExpr represent if expression.
	IfExpr struct {
		ExprImpl
		Cond Expr
		Then *BlockExpr
		Else Expr // Else is always *IfExpr or *BlockExpr (else if or else)
	}
)
