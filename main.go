package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/agatan/gray/codegen"
	"github.com/agatan/gray/parser"
	"github.com/agatan/gray/types"
)

var output string

func init() {
	flag.StringVar(&output, "o", "", "output file name")
	flag.Parse()
}

func main() {
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	filename := flag.Arg(0)
	if path.Ext(filename) != ".gy" {
		panic("given file is not gray's source code")
	}
	fp, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	l := parser.NewLexer(filename, fp)
	ds, err := parser.Parse(l)
	if err != nil {
		fmt.Println(err)
		return
	}
	checker := types.NewChecker(filename)
	scope, typemap, err := checker.Check(ds)
	if err != nil {
		panic(err)
	}
	scope.Dump(os.Stdout, 0, true)
	ctx, err := codegen.NewContext(filename, scope, typemap)
	if err != nil {
		panic(err)
	}
	defer ctx.Dispose()

	if output == "" {
		base := path.Base(filename)
		output = base[:len(base)-2] + "o"
	}
	err = ctx.EmitObject(output, ds)
	if err != nil {
		panic(err)
	}
}
