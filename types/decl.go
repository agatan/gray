package types

import "github.com/agatan/gray/ast"

func (c *Checker) checkDecl(s *Scope, d ast.Decl) error {
	switch d := d.(type) {
	case *ast.FuncDecl:
		fty, err := c.checkType(s, d.Type)
		if err != nil {
			return err
		}
		fscope := s.newChild(d.Ident.Name)
		for _, v := range fty.(*Signature).Params().vars {
			fscope.Insert(v)
		}
		fty.(*Signature).scope = fscope
		f := NewFunc(d.Ident.Name, fty)
		f.scope = fscope
		s.Insert(f)
		return nil
	default:
		panic("internal error: unreachable")
	}
}
