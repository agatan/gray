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
