package types

import (
	"fmt"

	"github.com/agatan/gray/ast"
)

func (c *Checker) checkDecl(s *Scope, d ast.Decl) error {
	switch d := d.(type) {
	case *ast.FuncDecl:
		fty, err := c.checkType(s, d.Type)
		if err != nil {
			return err
		}
		sig := fty.(*Signature)
		fscope := s.newChild(d.Ident.Name)
		for _, v := range sig.Params().vars {
			fscope.Insert(v)
		}
		sig.scope = fscope
		f := NewFunc(d.Ident.Name, fty)
		f.scope = fscope
		s.Insert(f)
		ty, err := c.checkExpr(fscope, d.Body)
		if err != nil {
			return err
		}
		if !c.isSameType(sig.Result(), ty) {
			return &Error{
				Message: fmt.Sprintf("type mismatch: expected %s, but got %s", sig.Result(), ty),
				Pos:     d.Position(),
			}
		}
		return nil
	default:
		panic("internal error: unreachable")
	}
}
