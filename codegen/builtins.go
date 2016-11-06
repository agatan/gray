package codegen

import (
	"errors"

	"github.com/agatan/gray/types"
	"llvm.org/llvm/bindings/go/llvm"
)

func (c *Context) unitValue() llvm.Value {
	return llvm.ConstInt(c.unitType(), 0, false)
}

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

var errNoGrayMain = errors.New("'gray.main' does not exist")

func (c *Context) defGrayMain() error {
	v := c.llmodule.NamedFunction("gray.main")
	if v.IsNil() || v.IsNull() {
		return errNoGrayMain
	}
	fty := llvm.FunctionType(c.intType(), []llvm.Type{}, false)
	f := llvm.AddFunction(c.llmodule, "main", fty)
	bb := llvm.AddBasicBlock(f, "entry")
	c.llbuilder.SetInsertPointAtEnd(bb)
	ret := c.llbuilder.CreateCall(v, []llvm.Value{}, "calltmp")
	c.llbuilder.CreateRet(ret)
	return nil
}
