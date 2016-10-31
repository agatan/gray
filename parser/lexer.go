package parser

import (
	"fmt"
	"io"
	"log"
	"os"
	"text/scanner"

	"github.com/agatan/gray/ast"
	"github.com/agatan/gray/token"
)

//go:generate go tool yacc -o parser.go parser.go.y

var debugMode = os.Getenv("GRAY_DEBUG") != ""

func debugf(f string, args ...interface{}) {
	if debugMode {
		log.Printf("debug: "+f, args...)
	}
}

func init() {
	if debugMode {
		yyErrorVerbose = true
	}
}

const (
	EOF = -1
	EOL = '\n'
)

// Error provides a convenient interface for handling syntax error.
type Error struct {
	Message  string
	Pos      token.Position
	Filename string
	Fatal    bool
}

// Error returns the error message.
func (e *Error) Error() string {
	return e.Message
}

// Lexer represents lexing states.
type Lexer struct {
	scanner *scanner.Scanner
	expr    ast.Expr
	pos     token.Position
	err     error
}

func NewLexer(filename string, r io.Reader) *Lexer {
	s := &scanner.Scanner{}
	s.Filename = filename
	s.Init(r)
	return &Lexer{scanner: s}
}

var keywords map[string]int = map[string]int{
	"true":  TRUE,
	"false": FALSE,
}

func (l *Lexer) scan() (tok int, lit string, pos token.Position, err error) {
	spos := l.scanner.Pos()
	pos = token.Position{
		Line:   spos.Line,
		Column: spos.Column,
	}
	t := l.scanner.Scan()
	debugf("Scan: %s", scanner.TokenString(t))
	lit = l.scanner.TokenText()
	switch t {
	case scanner.EOF:
		tok = EOF
		return
	case scanner.Int:
		tok = INT
		return
	case scanner.Ident:
		if tok, ok := keywords[lit]; ok {
			return tok, lit, pos, err
		}
		tok = IDENT
	default:
		switch lit {
		case "(":
			tok = LPAREN
			return
		case ")":
			tok = RPAREN
			return
		default:
			err = fmt.Errorf("Unknown token: %v", lit)
			return
		}
	}
	return
}

func (l *Lexer) Lex(lval *yySymType) int {
	tok, lit, pos, err := l.scan()
	if err != nil {
		l.err = &Error{Message: err.Error(), Pos: pos, Fatal: true}
	}
	l.pos = pos
	lval.tok = token.Token{Kind: tok, Lit: lit}
	lval.tok.SetPosition(pos)
	return tok
}

func (l *Lexer) Error(msg string) {
	l.err = &Error{Message: msg, Pos: l.pos, Fatal: false}
}

func Parse(l *Lexer) (ast.Expr, error) {
	if yyParse(l) != 0 {
		return nil, l.err
	}
	return l.expr, l.err
}
