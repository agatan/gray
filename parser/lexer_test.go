package parser

import (
	"strings"
	"testing"
)

func TestLex(t *testing.T) {
	tests := []struct {
		src string
		tok int
		lit string
	}{
		{"", EOF, ""},
		{" 123 ", INT, "123"},
		{" ident ", IDENT, "ident"},
		{" true ", TRUE, "true"},
		{" false ", FALSE, "false"},
		{" -> ", ARROW, "->"},
	}
	for _, test := range tests {
		l := NewLexer("test.gy", strings.NewReader(test.src))
		tok, lit, _, err := l.scan()
		if err != nil {
			t.Fatal(err)
		}
		if tok != test.tok {
			t.Errorf("expected %d, but got: %d", test.tok, tok)
		}
		if lit != test.lit {
			t.Errorf(`expected "%s", but got "%s"`, test.lit, lit)
		}
	}
}
