package main

import (
	"fmt"
	"os"

	"github.com/agatan/gray/codegen"
	"github.com/agatan/gray/parser"
	"github.com/agatan/gray/types"
)

func main() {
	l := parser.NewLexer("<stdin>", os.Stdin)
	ds, err := parser.Parse(l)
	if err != nil {
		fmt.Println(err)
		return
	}
	checker := types.NewChecker("<stdin>")
	scope, typemap, err := checker.Check(ds)
	if err != nil {
		panic(err)
	}
	scope.Dump(os.Stdout, 0, true)
	ctx, err := codegen.NewContext("<stdin>", scope, typemap)
	if err != nil {
		panic(err)
	}
	defer ctx.Dispose()

	mod, err := ctx.Generate(ds)
	if err != nil {
		panic(err)
	}
	mod.Dump()
}
