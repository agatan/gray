package types

import (
	"testing"

	"github.com/agatan/gray/ast"
	"github.com/agatan/gray/token"
)

func TestExpr(t *testing.T) {
	tests := []struct {
		expr     ast.Expr
		expected Type
	}{
		{&ast.BasicLit{Kind: token.UNIT}, BasicTypes[Unit]},
		{
			&ast.RefExpr{Value: &ast.BasicLit{Kind: token.INT}},
			builtinGenericTypes[refType].Instantiate([]Type{BasicTypes[Int]}),
		},
		{
			&ast.DerefExpr{Ref: &ast.RefExpr{Value: &ast.BasicLit{Kind: token.INT}}},
			BasicTypes[Int],
		},
		{
			&ast.InfixExpr{
				Operator: "+",
				LHS:      &ast.BasicLit{Kind: token.INT},
				RHS:      &ast.BasicLit{Kind: token.INT},
			},
			BasicTypes[Int],
		},
		{
			&ast.CallExpr{
				Func: &ast.Ident{Name: BuiltinPrintInt},
				Args: []ast.Expr{&ast.BasicLit{Kind: token.INT}},
			},
			BasicTypes[Unit],
		},
		{
			&ast.IfExpr{
				Cond: &ast.BasicLit{Kind: token.BOOL},
				Then: &ast.BlockExpr{
					Stmts: []ast.Stmt{&ast.ExprStmt{X: &ast.BasicLit{Kind: token.INT}}},
				},
				Else: &ast.BlockExpr{
					Stmts: []ast.Stmt{&ast.ExprStmt{X: &ast.BasicLit{Kind: token.INT}}},
				},
			},
			BasicTypes[Int],
		},
		{
			&ast.IfExpr{
				Cond: &ast.BasicLit{Kind: token.BOOL},
				Then: &ast.BlockExpr{
					Stmts: []ast.Stmt{&ast.ReturnStmt{}},
				},
				Else: &ast.BlockExpr{
					Stmts: []ast.Stmt{&ast.ExprStmt{X: &ast.BasicLit{Kind: token.INT}}},
				},
			},
			BasicTypes[Int],
		},
	}

	for i, test := range tests {
		checker := NewChecker("TestExpr")
		ty, err := checker.checkExpr(checker.scope, test.expr)
		if err != nil {
			t.Errorf("#%d: %s", i, err)
		}
		if !checker.isSameType(test.expected, ty) {
			t.Errorf("#%d: type error: expected %s, but got %s", i, test.expected, ty)
		}
	}
}
