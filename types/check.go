package types

import (
	"fmt"

	"github.com/agatan/gray/ast"
)

// Checker contains a type checking state.
type Checker struct {
	Filename string
	scope    *Scope
}

// NewChecker creates a Checker with given file name.
func NewChecker(filename string) *Checker {
	return &Checker{Filename: filename, scope: NewScope(filename)}
}

func (c *Checker) checkType(s *Scope, t ast.Type) (Type, error) {
	switch t := t.(type) {
	case *ast.TypeIdent:
		obj := s.LookupParent(t.Name)
		if obj == nil {
			return nil, &Error{Message: fmt.Sprintf("Unknown type: %s", t.Name), Pos: t.Position()}
		}
		tobj, ok := obj.(*TypeName)
		if !ok {
			return nil, &Error{Message: fmt.Sprintf("Unknown type: %s", t.Name), Pos: t.Position()}
		}
		return tobj.Type(), nil
	case *ast.FuncType:
		vars := make([]*Var, len(t.Params))
		for i, p := range t.Params {
			pty, err := c.checkType(s, p.Type)
			if err != nil {
				return nil, err
			}
			vars[i] = NewVar(p.Ident.Name, pty)
		}
		var retty Type
		if t.Result != nil {
			r, err := c.checkType(s, t.Result)
			if err != nil {
				return nil, err
			}
			retty = r
		} else {
			retty = BuiltinTypes[Unit]
		}
		return NewSignature(nil, NewVars(vars...), retty), nil
	case *ast.InstType:
		basety, err := c.checkType(s, t.Base)
		if err != nil {
			return nil, err
		}
		base, ok := basety.(*GenericType)
		if !ok {
			return nil, &Error{Message: fmt.Sprintf("%s is not generic type", basety.String()), Pos: t.Position()}
		}
		args := make([]Type, len(t.Args))
		for i, a := range t.Args {
			aty, err := c.checkType(s, a)
			if err != nil {
				return nil, err
			}
			args[i] = aty
		}
		return base.Instantiate(args), nil
	default:
		panic("internal error: unreachable clause")
	}
}

// Check checks the types of given ast declarations.
func (c *Checker) Check(ds []ast.Decl) error {
	for _, d := range ds {
		if err := c.checkDecl(c.scope, d); err != nil {
			return err
		}
	}
	return nil
}
