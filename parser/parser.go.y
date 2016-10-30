%{
package parser

import (
	"github.com/agatan/gray/ast"
	"github.com/agatan/gray/token"
)
%}

%type<expr> expr

%union{
	expr ast.Expr
	tok token.Token
}

%token<tok> IDENT INT TRUE FALSE

%%

expr:
	IDENT
	{
		$$ = &ast.Ident{Name: $1.Lit}
		$$.SetPosition($1.Position())
		if l, ok := parserlex.(*Lexer); ok {
			l.expr = $$
		}
	}
	| INT
	{
		$$ = &ast.BasicLit{Kind: token.INT, Lit: $1.Lit}
		$$.SetPosition($1.Position())
		if l, ok := parserlex.(*Lexer); ok {
			l.expr = $$
		}
	}
	| TRUE
	{
		$$ = &ast.BasicLit{Kind: token.BOOL, Lit: $1.Lit}
		$$.SetPosition($1.Position())
		if l, ok := parserlex.(*Lexer); ok {
			l.expr = $$
		}
	}
	| FALSE
	{
		$$ = &ast.BasicLit{Kind: token.BOOL, Lit: $1.Lit}
		$$.SetPosition($1.Position())
		if l, ok := parserlex.(*Lexer); ok {
			l.expr = $$
		}
	}
