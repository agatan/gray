package codegen

import (
	"github.com/agatan/gray/types"
	"llvm.org/llvm/bindings/go/llvm"
)

func (c *Context) intType() llvm.Type {
	return c.llcontext.Int32Type()
}

func (c *Context) boolType() llvm.Type {
	return c.llcontext.Int1Type()
}

func (c *Context) unitType() llvm.Type {
	return c.llcontext.VoidType()
}

func (c *Context) stringType() llvm.Type {
	// string representation is a pair of length and pointer.
	str := c.llcontext.StructCreateNamed("string")
	str.StructSetBody([]llvm.Type{
		c.intType(),                              // length
		llvm.PointerType(c.llcontext.Int8Type()), // pointer to characters
	}, true)
	return str
}

func (c *Context) genType(typ types.Type) (llvm.Type, error) {
	switch typ := typ.(type) {
	case *types.Basic:
		switch typ.Kind() {
		case types.Unit:
			return c.unitType(), nil
		case types.Bool:
			return c.boolType(), nil
		case types.Int:
			return c.intType(), nil
		case types.String:
			return c.stringType(), nil
		default:
			panic("internal error: unreachable")
		}
	}
}
