package codegen

import (
	"fmt"

	"github.com/agatan/gray/token"
	"github.com/agatan/gray/types"
	"llvm.org/llvm/bindings/go/llvm"
)

func (c *Context) defBasicTypes() {
	// string representation is a pair of length and pointer.
	str := c.llcontext.StructCreateNamed("string")
	str.StructSetBody([]llvm.Type{
		c.llcontext.Int32Type(),                     // length
		llvm.PointerType(c.llcontext.Int8Type(), 0), // pointer to characters
	}, false)
	c.basicTypes = []llvm.Type{
		token.UNIT:   c.llcontext.Int1Type(),
		token.INT:    c.llcontext.Int32Type(),
		token.BOOL:   c.llcontext.Int1Type(),
		token.STRING: str,
	}
}

func (c *Context) intType() llvm.Type {
	return c.basicTypes[token.INT]
}

func (c *Context) boolType() llvm.Type {
	return c.basicTypes[token.BOOL]
}

func (c *Context) unitType() llvm.Type {
	return c.basicTypes[token.UNIT]
}

func (c *Context) stringType() llvm.Type {
	return c.basicTypes[token.STRING]
}

func (c *Context) sigType(sig *types.Signature) llvm.Type {
	params := make([]llvm.Type, sig.Params().Len())
	for i := 0; i < sig.Params().Len(); i++ {
		params[i] = c.genType(sig.Params().At(i).Type())
	}
	result := c.genType(sig.Result())
	return llvm.FunctionType(result, params, false)
}

func (c *Context) genType(typ types.Type) llvm.Type {
	switch typ := typ.(type) {
	case *types.Basic:
		switch typ.Kind() {
		case types.Unit:
			return c.unitType()
		case types.Bool:
			return c.boolType()
		case types.Int:
			return c.intType()
		case types.String:
			return c.stringType()
		default:
			panic("internal error: unreachable")
		}
	case *types.Signature:
		return c.sigType(typ)
	default:
		panic(fmt.Sprintf("unimplemented yet: %T", typ))
	}
}
