package types

import "github.com/agatan/gray/ast"

func (c *Checker) checkExpr(s *Scope, e ast.Expr) (Type, error) {
	switch e := e.(type) {
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
