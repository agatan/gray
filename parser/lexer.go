package parser

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"text/scanner"
	"unicode/utf8"

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
	return fmt.Sprintf("%s:%d:%d: %s", e.Filename, e.Pos.Line, e.Pos.Column, e.Message)
}

// Lexer represents lexing states.
type Lexer struct {
	scanner *scanner.Scanner
	decls   []ast.Decl
	pos     token.Position
	err     error
}

func NewLexer(filename string, r io.Reader) *Lexer {
	s := &scanner.Scanner{}
	s.Init(r)
	s.Whitespace = 1<<'\t' | 1<<' ' | 1<<'\r'
	s.Filename = filename
	return &Lexer{scanner: s}
}

var keywords map[string]int = map[string]int{
	"true":  TRUE,
	"false": FALSE,
	"let":   LET,
	"def":   DEF,
	"if":    IF,
	"else":  ELSE,
}

func convertString(s string) (string, error) {
	ret := []rune{}
	src := []rune(s)
	// skip '"'
	src = src[1 : len(src)-1]
	for i := 0; i < len(src); i++ {
		if src[i] != '\\' {
			ret = append(ret, src[i])
			continue
		}
		if i == len(src)-1 {
			return "", errors.New("unexpected end of string")
		}
		i++
		switch src[i] {
		case 'b':
			ret = append(ret, '\b')
		case 'f':
			ret = append(ret, '\f')
		case 'r':
			ret = append(ret, '\r')
		case 'n':
			ret = append(ret, '\n')
		case 't':
			ret = append(ret, '\t')
		case '"':
			ret = append(ret, '"')
		case '\\':
			ret = append(ret, '\\')
		default:
			return "", errors.New("unexpected '\\'")
		}
	}
	return string(ret), nil
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
		ch, _ := utf8.DecodeRune([]byte(lit))
		if 'A' <= ch && ch <= 'Z' {
			tok = UIDENT
		} else {
			tok = IDENT
		}
	case scanner.String:
		tok = STRING
		lit, err = convertString(lit)
		return
	default:
		switch t {
		case '-':
			if l.scanner.Peek() == '>' {
				l.scanner.Next()
				tok = ARROW
				lit = "->"
				return
			}
		}
		tok = int(t)
		return
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

func Parse(l *Lexer) ([]ast.Decl, error) {
	if yyParse(l) != 0 {
		return nil, l.err
	}
	return l.decls, l.err
}
