package types

import (
	"fmt"

	"github.com/agatan/gray/ast"
	"github.com/agatan/gray/token"
)

func (c *Checker) checkExpr(s *Scope, e ast.Expr) (Type, error) {
	ty, err := c.inferAndCheckExpr(s, e)
	if err != nil {
		return nil, err
	}
	c.typemap.Record(e, ty)
	return ty, nil
}

func (c *Checker) inferAndCheckExpr(s *Scope, e ast.Expr) (Type, error) {
	c.setID(e)
	switch e := e.(type) {
	case *ast.Ident:
		obj := s.LookupParent(e.Name)
		if obj == nil {
			return nil, &Error{
				Message: fmt.Sprintf("undefined variable %s", e.Name),
				Pos:     e.Position(),
			}
		}
		return obj.Type(), nil
	case *ast.BasicLit:
		switch e.Kind {
		case token.UNIT:
			return BasicTypes[Unit], nil
		case token.BOOL:
			return BasicTypes[Bool], nil
		case token.INT:
			return BasicTypes[Int], nil
		case token.STRING:
			return BasicTypes[String], nil
		default:
			panic("internal error: unreachable")
		}
	case *ast.RefExpr:
		ty, err := c.checkExpr(s, e.Value)
		if err != nil {
			return nil, err
		}
		ref := builtinGenericTypes[refType]
		t := ref.Instantiate([]Type{ty})
		return t, nil
	case *ast.DerefExpr:
		ty, err := c.checkExpr(s, e.Ref)
		if err != nil {
			return nil, err
		}
		ref, ok := ty.(*InstType)
		if !ok || !c.isCompatibleType(ref.Base(), builtinGenericTypes[refType]) {
			return nil, &Error{
				Message: fmt.Sprintf("type mismatch: expected reference type, but got %s", ty),
				Pos:     e.Position(),
			}
		}
		return ref.Args()[0], nil
	case *ast.InfixExpr:
		lty, err := c.checkExpr(s, e.LHS)
		if err != nil {
			return nil, err
		}
		rty, err := c.checkExpr(s, e.RHS)
		if err != nil {
			return nil, err
		}
		ty, err := c.checkInfixExpr(s, e.Operator, lty, rty, e.Position())
		if err != nil {
			return nil, err
		}
		return ty, nil
	case *ast.ParenExpr:
		ty, err := c.checkExpr(s, e.X)
		return ty, err
	case *ast.BlockExpr:
		if len(e.Stmts) == 0 {
			return BasicTypes[Unit], nil
		}
		blockScope := s.newChild("")
		isBang := false
		for _, stmt := range e.Stmts[:len(e.Stmts)-1] {
			if err := c.checkStmt(blockScope, stmt); err != nil {
				return nil, err
			}
			if isBangExitStmt(stmt) {
				isBang = true
			}
		}
		last := e.Stmts[len(e.Stmts)-1]
		if le, ok := last.(*ast.ExprStmt); ok {
			t, err := c.checkExpr(blockScope, le.X)
			return t, err
		}
		if err := c.checkStmt(blockScope, last); err != nil {
			return nil, err
		}
		if isBangExitStmt(last) {
			isBang = true
		}
		if isBang {
			return NewBangType(), nil
		}
		return BasicTypes[Unit], nil
	case *ast.CallExpr:
		ty, err := c.checkExpr(s, e.Func)
		if err != nil {
			return nil, err
		}
		sig, ok := ty.(*Signature)
		if !ok {
			return nil, &Error{
				Message: fmt.Sprintf("%s is not callable", ty),
				Pos:     e.Position(),
			}
		}
		if sig.Params().Len() != len(e.Args) {
			return nil, &Error{
				Message: fmt.Sprintf(
					"invalid number of argument: expected %d, but got %d",
					sig.Params().Len(),
					len(e.Args),
				),
				Pos: e.Position(),
			}
		}
		for i, arg := range e.Args {
			ty, err := c.checkExpr(s, arg)
			if err != nil {
				return nil, err
			}
			if !c.isCompatibleType(sig.Params().At(i).Type(), ty) {
				return nil, &Error{
					Message: fmt.Sprintf("type mismatch: expected %s, but got %s", sig.Params().At(i).Type(), ty),
					Pos:     arg.Position(),
				}
			}
		}
		return sig.Result(), nil
	case *ast.IfExpr:
		condTy, err := c.checkExpr(s, e.Cond)
		if err != nil {
			return nil, err
		}
		if !c.isCompatibleType(BasicTypes[Bool], condTy) {
			return nil, &Error{
				Message: fmt.Sprintf("type mismatch: expected %s, but got %s", BasicTypes[Bool], condTy),
				Pos:     e.Cond.Position(),
			}
		}
		thenTy, err := c.checkExpr(s, e.Then)
		if err != nil {
			return nil, err
		}
		var elseTy Type
		if e.Else != nil {
			ety, err := c.checkExpr(s, e.Else)
			if err != nil {
				return nil, err
			}
			elseTy = ety
		} else {
			elseTy = BasicTypes[Unit]
		}
		resTy, ok := c.compatibleType(thenTy, elseTy)
		if !ok {
			return nil, &Error{
				Message: fmt.Sprintf("type mismatch: %s and %s", thenTy, elseTy),
				Pos:     e.Position(),
			}
		}
		return resTy, nil
	default:
		panic("internal error: unreachable")
	}
}
