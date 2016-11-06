package codegen

import (
	"fmt"

	"github.com/agatan/gray/ast"
	"github.com/agatan/gray/types"
	"llvm.org/llvm/bindings/go/llvm"
)

func (c *Context) forwardDecl(s *types.Scope, d ast.Decl) {
	switch d := d.(type) {
	case *ast.FuncDecl:
		f := s.Lookup(d.Ident.Name).(*types.Func)
		sig := f.Type().(*types.Signature)
		fty := c.sigType(sig)
		llvm.AddFunction(c.llmodule, c.mangle(f.Name()), fty)
	default:
		panic(fmt.Sprintf("unreachable %T", d))
	}
}

func (c *Context) forwardDecls(s *types.Scope, ds []ast.Decl) {
	for _, d := range ds {
		c.forwardDecl(s, d)
	}
}
