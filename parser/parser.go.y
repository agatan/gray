%{
package parser

import (
	"github.com/agatan/gray/ast"
	"github.com/agatan/gray/token"
)
%}

%type<decls> compdecls
%type<decls> decls
%type<decl> decl
%type<decl> func_decl

%type<stmts> compstmts
%type<stmts> stmts
%type<stmt> stmt
%type<stmt> let_stmt

%type<expr> expr
%type<if_expr> if_expr
%type<block_expr> block_expr
%type<expr> primitive_expr

%type<typ> typ

%type<ident> ident
%type<exprs> exprs
%type<exprs> rev_exprs
%type<params> params
%type<params> rev_params
%type<param> param

%union{
	decl       ast.Decl
	decls      []ast.Decl
	stmt       ast.Stmt
	stmts      []ast.Stmt
	expr       ast.Expr
	exprs      []ast.Expr
	block_expr *ast.BlockExpr
	if_expr    *ast.IfExpr
	typ        ast.Type

	ident      *ast.Ident
	tok        token.Token

	params     []*ast.Param
	param      *ast.Param

	term       token.Token
	terms      token.Token
	opt_terms  token.Token
}

%token<tok> IDENT UIDENT INT TRUE FALSE STRING
%token<tok> DEF LET IF ELSE RETURN
%token<tok> EQEQ NEQ GE LG OROR ANDAND
%token<tok> ARROW

%left OROR
%left ANDAND
%nonassoc EQEQ NEQ
%left '>' GE '<' LE

%left '+' '-'
%left '*' '/' '%'

%%

compdecls:
	opt_terms
	{
		$$ = nil
	}
	| decls opt_terms
	{
		$$ = $1
	}

decls:
	opt_terms decl
	{
		$$ = []ast.Decl{$2}
		if l, ok := yylex.(*Lexer); ok {
			l.decls = $$
		}
	}
	| decls terms decl
	{
		$$ = append($1, $3)
		if l, ok := yylex.(*Lexer); ok {
			l.decls = $$
		}
	}

decl:
	func_decl
	{
		$$ = $1
	}

func_decl:
	DEF opt_terms ident opt_terms '(' opt_terms params opt_terms ')' opt_terms block_expr
	{
		$$ = &ast.FuncDecl{
			Ident: $3,
			Type: &ast.FuncType{Params: $7, Result: nil},
			Body: $11,
		}
		$$.SetPosition($1.Position())
	}
	| DEF opt_terms ident opt_terms '(' opt_terms params opt_terms ')' opt_terms ARROW opt_terms typ opt_terms block_expr
	{
		$$ = &ast.FuncDecl{
			Ident: $3,
			Type: &ast.FuncType{Params: $7, Result: $13},
			Body: $15,
		}
		$$.SetPosition($1.Position())
	}

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
	}
	| stmts terms stmt
	{
		$$ = append($1, $3)
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
	| RETURN expr
	{
		$$ = &ast.ReturnStmt{X: $2}
		$$.SetPosition($1.Position())
	}
	| RETURN
	{
		$$ = &ast.ReturnStmt{}
		$$.SetPosition($1.Position())
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
	| block_expr
	{
		$$ = $1
	}
	| if_expr
	{
		$$ = $1
	}
	| expr '(' opt_terms exprs opt_terms ')'
	{
		$$ = &ast.CallExpr{Func: $1, Args: $4}
		$$.SetPosition($1.Position())
	}
	| expr '+' expr
	{
		$$ = &ast.InfixExpr{LHS: $1, Operator: "+", RHS: $3}
		$$.SetPosition($1.Position())
	}
	| expr '-' expr
	{
		$$ = &ast.InfixExpr{LHS: $1, Operator: "-", RHS: $3}
		$$.SetPosition($1.Position())
	}
	| expr '*' expr
	{
		$$ = &ast.InfixExpr{LHS: $1, Operator: "*", RHS: $3}
		$$.SetPosition($1.Position())
	}
	| expr '/' expr
	{
		$$ = &ast.InfixExpr{LHS: $1, Operator: "/", RHS: $3}
		$$.SetPosition($1.Position())
	}
	| expr '%' expr
	{
		$$ = &ast.InfixExpr{LHS: $1, Operator: "%", RHS: $3}
		$$.SetPosition($1.Position())
	}
	| expr OROR expr
	{
		$$ = &ast.InfixExpr{LHS: $1, Operator: "||", RHS: $3}
		$$.SetPosition($1.Position())
	}
	| expr ANDAND expr
	{
		$$ = &ast.InfixExpr{LHS: $1, Operator: "&&", RHS: $3}
		$$.SetPosition($1.Position())
	}
	| expr EQEQ expr
	{
		$$ = &ast.InfixExpr{LHS: $1, Operator: "==", RHS: $3}
		$$.SetPosition($1.Position())
	}
	| expr NEQ expr
	{
		$$ = &ast.InfixExpr{LHS: $1, Operator: "!=", RHS: $3}
		$$.SetPosition($1.Position())
	}
	| expr '<' expr
	{
		$$ = &ast.InfixExpr{LHS: $1, Operator: "<", RHS: $3}
		$$.SetPosition($1.Position())
	}
	| expr '>' expr
	{
		$$ = &ast.InfixExpr{LHS: $1, Operator: ">", RHS: $3}
		$$.SetPosition($1.Position())
	}
	| expr GE expr
	{
		$$ = &ast.InfixExpr{LHS: $1, Operator: ">=", RHS: $3}
		$$.SetPosition($1.Position())
	}
	| expr LE expr
	{
		$$ = &ast.InfixExpr{LHS: $1, Operator: "<=", RHS: $3}
		$$.SetPosition($1.Position())
	}

block_expr:
	'{' compstmts '}'
	{
		$$ = &ast.BlockExpr{Stmts: $2}
		if l, ok := yylex.(*Lexer); ok {
			$$.SetPosition(l.pos)
		}
	}

if_expr:
	IF expr block_expr
	{
		$$ = &ast.IfExpr{Cond: $2, Then: $3}
		$$.SetPosition($1.Position())
	}
	| IF expr block_expr ELSE block_expr
	{
		$$ = &ast.IfExpr{Cond: $2, Then: $3, Else: $5}
		$$.SetPosition($1.Position())
	}
	| IF expr block_expr ELSE if_expr
	{
		$$ = &ast.IfExpr{Cond: $2, Then: $3, Else: $5}
		$$.SetPosition($1.Position())
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
	| STRING
	{
		$$ = &ast.BasicLit{Kind: token.STRING, Lit: $1.Lit}
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

exprs:
	rev_exprs
	{
		$$ = make([]ast.Expr, len($1))
		for i, p := range $1 {
			$$[len($1)-i-1] = p
		}
	}

rev_exprs:
	/* empty */
	{
		$$ = nil
	}
	| expr
	{
		$$ = []ast.Expr{$1}
	}
	| expr ',' opt_terms rev_exprs
	{
		$$ = append($4, $1)
	}

params:
	rev_params
	{
		$$ = make([]*ast.Param, len($1))
		for i, p := range $1 {
			$$[len($1)-i-1] = p
		}
	}

rev_params:
	/* empty */
	{
		$$ = nil
	}
	| param
	{
		$$ = []*ast.Param{$1}
	}
	| param ',' opt_terms rev_params
	{
		$$ = append($4, $1)
	}

param:
	ident ':' typ
	{
		$$ = &ast.Param{Ident: $1, Type: $3}
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
