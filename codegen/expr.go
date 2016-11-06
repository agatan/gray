package codegen

import (
	"fmt"

	"github.com/agatan/gray/ast"
	"github.com/agatan/gray/token"
	"llvm.org/llvm/bindings/go/llvm"
)

func (c *Context) genExpr(e ast.Expr) (llvm.Value, error) {
	switch e := e.(type) {
	case *ast.BasicLit:
		switch e.Kind {
		case token.UNIT:
			return llvm.ConstInt(c.unitType(), 0, false), nil
		case token.BOOL:
			if e.Lit == "true" {
				return llvm.ConstInt(c.boolType(), 1, false), nil
			}
			return llvm.ConstInt(c.boolType(), 0, false), nil
		case token.INT:
			return llvm.ConstIntFromString(c.intType(), e.Lit, 10), nil
		case token.STRING:
			sptr := c.llbuilder.CreateGlobalStringPtr(e.Lit, fmt.Sprintf("str.%d", e.ID()))
			length := llvm.ConstInt(c.intType(), uint64(len(e.Lit)), false)
			return llvm.ConstStruct([]llvm.Value{length, sptr}, false), nil
		default:
			panic(fmt.Sprintf("unreachable %#v", e))
		}
	case *ast.BlockExpr:
		if len(e.Stmts) == 0 {
			return c.unitValue(), nil
		}
		if !e.IsExpr {
			panic("statement block is unimplemented yet")
		}
		// TODO: gen statements
		return c.genExpr(e.Stmts[len(e.Stmts)-1].(*ast.ExprStmt).X)
	default:
		panic(fmt.Sprintf("unimplemented %T", e))
	}
}
