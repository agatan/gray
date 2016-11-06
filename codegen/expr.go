package codegen

import (
	"fmt"

	"github.com/agatan/gray/ast"
	"github.com/agatan/gray/token"
	"github.com/agatan/gray/types"
	"llvm.org/llvm/bindings/go/llvm"
)

func (c *Context) genCallExpr(e *ast.CallExpr) llvm.Value {
	var fvalue llvm.Value
	switch f := e.Func.(type) {
	case *ast.Ident:
		fobj := c.currentScope().LookupParent(f.Name).(*types.Func)
		if fobj.IsBuiltin() {
			// is builtin function, do not mangle the name.
			fvalue = c.llmodule.NamedFunction(fobj.Name())
		} else {
			fvalue = c.llmodule.NamedFunction(c.mangle(fobj.Name()))
		}
	default:
		panic(fmt.Sprintf("unreachable: %T", e.Func))
	}
	args := make([]llvm.Value, len(e.Args))
	for i, arg := range e.Args {
		args[i] = c.genExpr(arg)
	}
	return c.llbuilder.CreateCall(fvalue, args, "calltmp")
}

func (c *Context) genExpr(e ast.Expr) llvm.Value {
	switch e := e.(type) {
	case *ast.BasicLit:
		switch e.Kind {
		case token.UNIT:
			return llvm.ConstInt(c.unitType(), 0, false)
		case token.BOOL:
			if e.Lit == "true" {
				return llvm.ConstInt(c.boolType(), 1, false)
			}
			return llvm.ConstInt(c.boolType(), 0, false)
		case token.INT:
			return llvm.ConstIntFromString(c.intType(), e.Lit, 10)
		case token.STRING:
			sptr := c.llbuilder.CreateGlobalStringPtr(e.Lit, fmt.Sprintf("str.%d", e.ID()))
			length := llvm.ConstInt(c.intType(), uint64(len(e.Lit)), false)
			return llvm.ConstStruct([]llvm.Value{length, sptr}, false)
		default:
			panic(fmt.Sprintf("unreachable %#v", e))
		}
	case *ast.Ident:
		return c.valuemap.LookupParent(e.Name)
	case *ast.CallExpr:
		return c.genCallExpr(e)
	case *ast.BlockExpr:
		if len(e.Stmts) == 0 {
			return c.unitValue()
		}
		last := e.Stmts[len(e.Stmts)-1]
		for _, s := range e.Stmts[:len(e.Stmts)-1] {
			c.genStmt(s)
		}
		if e.IsExpr {
			return c.genExpr(last.(*ast.ExprStmt).X)
		}
		c.genStmt(last)
		return c.unitValue()
	default:
		panic(fmt.Sprintf("unimplemented %T", e))
	}
}
