%{
package parser

import (
	"github.com/agatan/gray/ast"
	"github.com/agatan/gray/token"
)
%}

%type<expr> primitive_expr
%type<expr> expr

%union{
	expr ast.Expr
	tok token.Token
}

%token<tok> IDENT INT TRUE FALSE
%token<tok> LPAREN RPAREN

%%

expr:
	primitive_expr
	{
		$$ = $1
		if l, ok := yylex.(*Lexer); ok {
			l.expr = $$
		}
	}

primitive_expr:
	IDENT
	{
		$$ = &ast.Ident{Name: $1.Lit}
		$$.SetPosition($1.Position())
	}
	| INT
	{
		$$ = &ast.BasicLit{Kind: token.INT, Lit: $1.Lit}
		$$.SetPosition($1.Position())
	}
	| TRUE
	{
		$$ = &ast.BasicLit{Kind: token.BOOL, Lit: $1.Lit}
		$$.SetPosition($1.Position())
	}
	| FALSE
	{
		$$ = &ast.BasicLit{Kind: token.BOOL, Lit: $1.Lit}
		$$.SetPosition($1.Position())
	}
	| LPAREN expr RPAREN
	{
		$$ = &ast.ParenExpr{X: $2}
		$$.SetPosition($1.Position())
	}
