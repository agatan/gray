package parser

import (
	"fmt"
	"strings"
	"testing"
)

func TestExpr(t *testing.T) {
	tests := []struct {
		src string
	}{
		{" 123 "},
		{" ident "},
		{" true "},
		{" false "},
		{" ( 123 ) "},
		{" f() "},
		{" f(1, 2) "},
		{" f(1, 2,) "},
		{" (f)(1, 2,) "},
	}
	for _, test := range tests {
		l := NewLexer("test.gy", strings.NewReader(fmt.Sprintf("def test() { %s }", test.src)))
		expr, err := Parse(l)
		if err != nil {
			t.Error(err)
		}
		if expr == nil {
			t.Error("result is nil")
		}
	}
}

func TestLetStmt(t *testing.T) {
	tests := []string{` let x = 1 `, ` let x : Int = 1 `}
	for _, test := range tests {
		l := NewLexer("test.gy", strings.NewReader(fmt.Sprintf("def test() { %s }", test)))
		ss, err := Parse(l)
		if err != nil {
			t.Error(err)
		}
		if len(ss) != 1 {
			t.Error("let statement does not work")
		}
	}
}

func TestIfStmt(t *testing.T) {
	tests := []string{
		` if x { 0 } `,
		` if x { 0 } else { 1 } `,
		` if x { 0 } else if y { 1 } `,
		` if x { 0 } else if y { 1 } else { 2 } `,
	}
	for i, test := range tests {
		l := NewLexer("test.gy", strings.NewReader(fmt.Sprintf("def test(x: Bool) -> Int { %s }", test)))
		_, err := Parse(l)
		if err != nil {
			t.Fatalf("#%d: parsing '%s': %s", i, test, err)
		}
	}
}

func TestDefDecl(t *testing.T) {
	test := `
	def test(
		x: Int,
		y: Bool,
	) -> Bool {
		y
	}
	`
	l := NewLexer("test.gy", strings.NewReader(test))
	ds, err := Parse(l)
	if err != nil {
		t.Error(err)
	}
	if len(ds) != 1 {
		t.Error("def declaration does not work")
	}
}
