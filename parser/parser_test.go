package parser

import (
	"strings"
	"testing"
)

func TestPrimitiveExpr(t *testing.T) {
	tests := []struct {
		src string
	}{
		{" 123 "},
		{" ident "},
		{" true "},
		{" false "},
		{" ( 123 ) "},
	}
	for _, test := range tests {
		l := NewLexer("test.gy", strings.NewReader(test.src))
		expr, err := Parse(l)
		if err != nil {
			t.Error(err)
		}
		if expr == nil {
			t.Error("result is nil")
		}
	}
}
