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
			return BasicTypes[Unit], nil
		case token.BOOL:
			return BasicTypes[Bool], nil
		case token.INT:
			return BasicTypes[Int], nil
		case token.STRING:
			return BasicTypes[String], nil
		default:
			panic("internal error: unreachable")
		}
	case *ast.RefExpr:
		ty, err := c.checkExpr(s, e.Value)
		if err != nil {
			return nil, err
		}
		ref := builtinGenericTypes[refType]
		t := ref.Instantiate([]Type{ty})
		return t, nil
	case *ast.BlockExpr:
		if len(e.Stmts) == 0 {
			return BasicTypes[Unit], nil
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
		return BasicTypes[Unit], nil
	default:
		panic("unimplemented yet")
	}
}
