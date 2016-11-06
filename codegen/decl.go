package codegen

import (
	"fmt"

	"github.com/agatan/gray/ast"
	"llvm.org/llvm/bindings/go/llvm"
)

func (c *Context) genDecl(d ast.Decl) error {
	switch d := d.(type) {
	case *ast.FuncDecl:
		f := c.llmodule.NamedFunction(c.mangle(d.Ident.Name))
		c.enterNewScope()
		defer c.exitScope()
		// define parameters
		for i, param := range f.Params() {
			c.valuemap.Insert(d.Type.Params[i].Ident.Name, param)
		}
		bb := llvm.AddBasicBlock(f, "entry")
		c.llbuilder.SetInsertPointAtEnd(bb)
		v, err := c.genExpr(d.Body)
		if err != nil {
			return err
		}
		c.llbuilder.CreateRet(v)
		return nil
	default:
		panic(fmt.Sprintf("unreachable %T", d))
	}
}
