package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/agatan/gray/codegen"
	"github.com/agatan/gray/parser"
	"github.com/agatan/gray/types"
)

var output string
var dumpIR bool

func init() {
	flag.StringVar(&output, "o", "", "output file name")
	flag.BoolVar(&dumpIR, "dump", false, "dump LLVM IR")
	flag.Parse()
}

func moduleName(filename string) string {
	name := strings.Replace(filename, "/", ".", -1)
	return name[:len(name)-3]
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
	modname := moduleName(filename)
	l := parser.NewLexer(modname, fp)
	ds, err := parser.Parse(l)
	if err != nil {
		fmt.Println(err)
		return
	}
	checker := types.NewChecker(modname)
	scope, typemap, err := checker.Check(ds)
	if err != nil {
		panic(err)
	}
	scope.Dump(os.Stdout, 0, true)
	ctx, err := codegen.NewContext(modname, scope, typemap)
	if err != nil {
		panic(err)
	}
	defer ctx.Dispose()

	llmod, err := ctx.Generate(ds)
	if err != nil {
		panic(err)
	}
	if dumpIR {
		llmod.Dump()
		os.Exit(0)
	}

	if output == "" {
		base := path.Base(filename)
		output = base[:len(base)-2] + "o"
	}
	err = ctx.EmitObject(output, llmod)
	if err != nil {
		panic(err)
	}
}
