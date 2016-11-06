package codegen

import (
	"fmt"

	"github.com/agatan/gray/ast"
)

func (c *Context) genStmt(s ast.Stmt) {
	switch s := s.(type) {
	case *ast.ExprStmt:
		c.genExpr(s.X)
	case *ast.LetStmt:
		v := c.genExpr(s.Value)
		c.valuemap.Insert(s.Ident.Name, v)
	case *ast.ReturnStmt:
		if s.X == nil {
			c.llbuilder.CreateRetVoid()
		} else {
			c.llbuilder.CreateRet(c.genExpr(s.X))
		}
	default:
		panic(fmt.Sprintf("unimplemented yet: %T", s))
	}
}
