package types

import "github.com/agatan/gray/ast"

func (c *Checker) checkDecl(s *Scope, d ast.Decl) error {
	switch d := d.(type) {
	case *ast.FuncDecl:
		fty, err := c.checkType(s, d.Type)
		if err != nil {
			return err
		}
		s.Insert(NewFunc(d.Ident.Name, fty))
		return nil
	default:
		panic("internal error: unreachable")
	}
}
