package codegen

import (
	"fmt"

	"llvm.org/llvm/bindings/go/llvm"

	"github.com/agatan/gray/ast"
)

func (c *Context) genStmt(s ast.Stmt) {
	switch s := s.(type) {
	case *ast.ExprStmt:
		c.genExpr(s.X)
	case *ast.LetStmt:
		v := c.genExpr(s.Value)
		v.SetName(s.Ident.Name)
		c.valuemap.Insert(s.Ident.Name, v)
	case *ast.ReturnStmt:
		if s.X == nil {
			c.llbuilder.CreateRetVoid()
		} else {
			c.llbuilder.CreateRet(c.genExpr(s.X))
		}
	case *ast.WhileStmt:
		parentb := c.llbuilder.GetInsertBlock().Parent()
		condb := llvm.AddBasicBlock(parentb, "cond")
		bodyb := llvm.AddBasicBlock(parentb, "body")
		exitb := llvm.AddBasicBlock(parentb, "exit")
		// generate cond block
		c.llbuilder.CreateBr(condb)
		c.llbuilder.SetInsertPointAtEnd(condb)
		condv := c.genExpr(s.Cond)
		cmp := c.llbuilder.CreateICmp(llvm.IntNE, condv, llvm.ConstInt(c.boolType(), 0, false), "whilecond")
		c.llbuilder.CreateCondBr(cmp, bodyb, exitb)
		// generate body block
		c.llbuilder.SetInsertPointAtEnd(bodyb)
		c.genExpr(s.Body)
		c.llbuilder.CreateBr(condb)
		// generate exit block
		c.llbuilder.SetInsertPointAtEnd(exitb)
	default:
		panic(fmt.Sprintf("unimplemented yet: %T", s))
	}
}
