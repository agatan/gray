package main

import (
	"fmt"

	"github.com/agatan/gray/ast"
	"github.com/agatan/gray/token"
)

func main() {
	ds := []ast.Decl{
		&ast.FuncDecl{
			Ident: &ast.Ident{Name: "main"},
			Type: &ast.FuncType{
				Params: []*ast.Param{},
				Result: &ast.TypeIdent{Name: "Unit"},
			},
			Body: &ast.BlockExpr{
				Stmts: []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.BasicLit{
							Kind: token.UNIT,
							Lit:  "()",
						},
					},
				},
			},
		},
	}
	fmt.Println(ds)
}
