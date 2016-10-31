%{
package parser

import (
	"github.com/agatan/gray/ast"
	"github.com/agatan/gray/token"
)
%}

%type<stmts> compstmts
%type<stmts> stmts
%type<stmt> stmt
%type<stmt> let_stmt

%type<expr> expr
%type<expr> primitive_expr

%type<typ> typ

%type<ident> ident

%union{
	stmt      ast.Stmt
	stmts     []ast.Stmt
	expr      ast.Expr
	typ       ast.Type

	ident     *ast.Ident
	tok       token.Token

	term      token.Token
	terms     token.Token
	opt_terms token.Token
}

%token<tok> IDENT UIDENT INT TRUE FALSE
%token<tok> LET

%%

compstmts:
	opt_terms
	{
		$$ = nil
	}
	| stmts opt_terms
	{
		$$ = $1
	}

stmts:
	opt_terms stmt
	{
		$$ = []ast.Stmt{$2}
		if l, ok := yylex.(*Lexer); ok {
			l.stmts = $$
		}
	}
	| stmts terms stmt
	{
		$$ = append($1, $3)
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
	LET opt_terms ident opt_terms '=' opt_terms expr
	{
		$$ = &ast.LetStmt {
			Ident: $3,
			Type: nil,
			Value: $7,
		}
		$$.SetPosition($1.Position())
	}
	| LET opt_terms ident opt_terms ':' typ '=' opt_terms expr
	{
		$$ = &ast.LetStmt {
			Ident: $3,
			Type: $6,
			Value: $9,
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
	| '(' opt_terms expr opt_terms ')'
	{
		$$ = &ast.ParenExpr{X: $3}
		if l, ok := yylex.(*Lexer); ok {
			$$.SetPosition(l.pos)
		}
	}

typ:
	UIDENT
	{
		$$ = &ast.TypeIdent{Name: $1.Lit}
		$$.SetPosition($1.Position())
	}

ident:
	IDENT
	{
		$$ = &ast.Ident{Name: $1.Lit}
		$$.SetPosition($1.Position())
	}

opt_terms:
	 /* empty */ | terms
	 ;

terms:
	 term | terms term
	 ;

term:
	';' | '\n'
	;
