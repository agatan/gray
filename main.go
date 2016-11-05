package main

import (
	"encoding/json"
	"fmt"
	"os"

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
	scope, _, err := checker.Check(ds)
	if err != nil {
		panic(err)
	}
	scope.Dump(os.Stdout, 0, true)
	bs, err := json.Marshal(ds)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(bs))
	}
}
