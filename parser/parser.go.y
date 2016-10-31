%{
package parser

import (
	"github.com/agatan/gray/ast"
	"github.com/agatan/gray/token"
)
%}

%type<stmts> stmts
%type<stmt> stmt
%type<stmt> let_stmt

%type<expr> expr
%type<expr> primitive_expr

%type<ident> ident

%union{
	stmt ast.Stmt
	stmts []ast.Stmt
	expr ast.Expr

	ident *ast.Ident
	tok token.Token
}

%token<tok> IDENT INT TRUE FALSE
%token<tok> LPAREN RPAREN EQ
%token<tok> LET

%%

stmts:
	{
		$$ = nil
		if l, ok := yylex.(*Lexer); ok {
			l.stmts = $$
		}
	}
	| stmts stmt
	{
		$$ = append($1, $2)
		if l, ok := yylex.(*Lexer); ok {
			l.stmts = $$
		}
	}

stmt:
	expr
	{
		$$ = &ast.ExprStmt{X: $1}
		$$.SetPosition($1.Position())
	}
	| let_stmt
	{
		$$ = $1
	}

let_stmt:
	LET ident EQ expr
	{
		$$ = &ast.LetStmt {
			Ident: $2,
			Type: nil,
			Value: $4,
		}
		$$.SetPosition($1.Position())
	}

expr:
	primitive_expr
	{
		$$ = $1
	}

primitive_expr:
	ident
	{
		$$ = $1
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

ident:
	IDENT
	{
		$$ = &ast.Ident{Name: $1.Lit}
		$$.SetPosition($1.Position())
	}
