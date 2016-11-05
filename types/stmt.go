package types

import (
	"fmt"

	"github.com/agatan/gray/ast"
)

func (c *Checker) checkStmt(s *Scope, stmt ast.Stmt) error {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		ty, err := c.checkExpr(s, stmt.X)
		if err != nil {
			return err
		}
		if !c.isCompatibleType(BasicTypes[Unit], ty) {
			return &Error{
				Message: fmt.Sprintf("type mismatch: expected %s, but got %s", BasicTypes[Unit], ty),
				Pos:     stmt.Position(),
			}
		}
		return nil
	case *ast.LetStmt:
		ety, err := c.checkExpr(s, stmt.Value)
		if err != nil {
			return err
		}
		if stmt.Type != nil {
			ty, err := c.checkType(s, stmt.Type)
			if err != nil {
				return err
			}
			if !c.isCompatibleType(ty, ety) {
				return &Error{
					Message: fmt.Sprintf("type mismatch: expected %s, but got %s", ty, ety),
					Pos:     stmt.Value.Position(),
				}
			}
		}
		s.Insert(NewVar(stmt.Ident.Name, ety))
		return nil
	case *ast.ReturnStmt:
		if stmt.X == nil {
			c.addReturnInfo(BasicTypes[Unit], stmt.Position())
			return nil
		}
		ety, err := c.checkExpr(s, stmt.X)
		if err != nil {
			return nil
		}
		c.addReturnInfo(ety, stmt.Position())
		return nil
	default:
		panic("unimplemented yet")
	}
}
