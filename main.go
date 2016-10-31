package main

import (
	"fmt"
	"os"

	"github.com/agatan/gray/parser"
)

func main() {
	l := parser.NewLexer("<stdin>", os.Stdin)
	ss, err := parser.Parse(l)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v\n", ss)
	}
}
