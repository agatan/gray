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
			return nil
		}
		if !c.isSameType(BasicTypes[Unit], ty) {
			return &Error{
				Message: fmt.Sprintf("type mismatch: expected %s, but got %s", BasicTypes[Unit], ty),
				Pos:     stmt.Position(),
			}
		}
		return nil
	default:
		panic("unimplemented yet")
	}
}
