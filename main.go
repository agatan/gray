package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/agatan/gray/parser"
)

func main() {
	l := parser.NewLexer("<stdin>", os.Stdin)
	ds, err := parser.Parse(l)
	if err != nil {
		fmt.Println(err)
		return
	}
	bs, err := json.Marshal(ds)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(bs))
	}
}
