package types

import (
	"fmt"

	"github.com/agatan/gray/ast"
	"github.com/agatan/gray/token"
)

func (c *Checker) checkExpr(s *Scope, e ast.Expr) (Type, error) {
	switch e := e.(type) {
	case *ast.Ident:
		obj := s.LookupParent(e.Name)
		if obj == nil {
			return nil, &Error{
				Message: fmt.Sprintf("undefined variable %s", e.Name),
				Pos:     e.Position(),
			}
		}
		return obj.Type(), nil
	case *ast.BasicLit:
		switch e.Kind {
		case token.UNIT:
			return BuiltinTypes[Unit], nil
		case token.BOOL:
			return BuiltinTypes[Bool], nil
		case token.INT:
			return BuiltinTypes[Int], nil
		case token.STRING:
			return BuiltinTypes[String], nil
		default:
			panic("internal error: unreachable")
		}
	case *ast.BlockExpr:
		if len(e.Stmts) == 0 {
			return BuiltinTypes[Unit], nil
		}
		blockScope := s.newChild("")
		for _, stmt := range e.Stmts[:len(e.Stmts)-1] {
			if err := c.checkStmt(blockScope, stmt); err != nil {
				return nil, err
			}
		}
		last := e.Stmts[len(e.Stmts)-1]
		if le, ok := last.(*ast.ExprStmt); ok {
			t, err := c.checkExpr(blockScope, le.X)
			return t, err
		}
		if err := c.checkStmt(blockScope, last); err != nil {
			return nil, err
		}
		return BuiltinTypes[Unit], nil
	default:
		panic("unimplemented yet")
	}
}
