package codegen

import (
	"github.com/agatan/gray/types"
	"llvm.org/llvm/bindings/go/llvm"
)

func (c *Context) defBuiltinFunctions() error {
	for _, t := range types.BuiltinFunctions {
		name := t.Name
		sig := t.Sig
		fty, err := c.sigType(sig)
		if err != nil {
			return err
		}
		llvm.AddFunction(c.llmodule, name, fty)
	}
	return nil
}
