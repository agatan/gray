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
	if err := checker.Check(ds); err != nil {
		panic(err)
	}
	bs, err := json.Marshal(ds)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(bs))
	}
}
