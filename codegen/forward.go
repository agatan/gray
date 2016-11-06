package codegen

import (
	"fmt"

	"github.com/agatan/gray/ast"
	"github.com/agatan/gray/types"
	"llvm.org/llvm/bindings/go/llvm"
)

func (c *Context) forwardDecl(s *types.Scope, d ast.Decl) error {
	switch d := d.(type) {
	case *ast.FuncDecl:
		f := s.Lookup(d.Ident.Name).(*types.Func)
		sig := f.Type().(*types.Signature)
		fty, err := c.sigType(sig)
		if err != nil {
			return err
		}
		llvm.AddFunction(c.llmodule, c.mangle(f.Name()), fty)
		return nil
	default:
		panic(fmt.Sprintf("unreachable %T", d))
	}
}

func (c *Context) forwardDecls(s *types.Scope, ds []ast.Decl) error {
	for _, d := range ds {
		if err := c.forwardDecl(s, d); err != nil {
			return err
		}
	}
	return nil
}
