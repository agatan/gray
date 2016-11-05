package types

import (
	"fmt"

	"github.com/agatan/gray/ast"
)

// Checker contains a type checking state.
type Checker struct {
	Filename    string
	scope       *Scope
	returnInfos []returnInfo // returnTypes is a set of return types in the current function body.
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
			// Result type is Unit.
			retty = BasicTypes[Unit]
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

func (c *Checker) isSameType(lhs, rhs Type) bool {
	switch lhs := lhs.(type) {
	case *Basic:
		rhs, ok := rhs.(*Basic)
		if !ok {
			return false
		}
		return lhs.kind == rhs.kind
	case *Signature:
		rhs, ok := rhs.(*Signature)
		if !ok {
			return false
		}
		if lhs.Params().Len() != rhs.Params().Len() {
			return false
		}
		for i := 0; i < lhs.Params().Len(); i++ {
			if !c.isSameType(lhs.Params().At(i).Type(), rhs.Params().At(i).Type()) {
				return false
			}
		}
		return c.isSameType(lhs.Result(), rhs.Result())
	case *GenericType:
		rhs, ok := rhs.(*GenericType)
		if !ok {
			return false
		}
		// Address comparison (every same generic types should be the same object)
		return lhs == rhs
	case *InstType:
		rhs, ok := rhs.(*InstType)
		if !ok {
			return false
		}
		if !c.isSameType(lhs.Base(), rhs.Base()) {
			return false
		}
		if len(lhs.Args()) != len(rhs.Args()) {
			return false
		}
		for i := 0; i < len(lhs.Args()); i++ {
			if !c.isSameType(lhs.Args()[i], rhs.Args()[i]) {
				return false
			}
		}
		return true
	default:
		panic("internal error: unreachable")
	}
}

func (c *Checker) compatibleType(lhs, rhs Type) (Type, bool) {
	if _, ok := lhs.(*BangType); ok {
		return rhs, true
	}
	if _, ok := rhs.(*BangType); ok {
		return lhs, true
	}
	if c.isSameType(lhs, rhs) {
		return lhs, true
	}
	return nil, false
}

func (c *Checker) isCompatibleType(lhs, rhs Type) bool {
	_, ok := c.compatibleType(lhs, rhs)
	return ok
}

// Check checks the types of given ast declarations.
func (c *Checker) Check(ds []ast.Decl) (*Scope, error) {
	for _, d := range ds {
		if err := c.checkDecl(c.scope, d); err != nil {
			if e, ok := err.(*Error); ok {
				e.Filename = c.Filename
			}
			return nil, err
		}
	}
	return c.scope, nil
}
