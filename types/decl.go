package types

import (
	"fmt"

	"github.com/agatan/gray/ast"
)

func (c *Checker) checkDecl(s *Scope, d ast.Decl) error {
	c.setID(d)
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
		c.resetReturnInfos()
		ty, err := c.checkExpr(fscope, d.Body)
		if err != nil {
			return err
		}
		if _, ok := c.compatibleType(sig.Result(), ty); !ok {
			return &Error{
				Message: fmt.Sprintf("type mismatch: expected %s, but got %s", sig.Result(), ty),
				Pos:     d.Position(),
			}
		}
		for _, info := range c.currentReturnInfos() {
			if _, ok := c.compatibleType(sig.Result(), info.typ); !ok {
				return &Error{
					Message: fmt.Sprintf("type mismatch: expected %s, but got %s", sig.Result(), info.typ),
					Pos:     info.pos,
				}
			}
		}
		return nil
	default:
		panic("internal error: unreachable")
	}
}
