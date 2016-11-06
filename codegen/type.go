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
	}, true)
	c.basicTypes = []llvm.Type{
		token.UNIT:   c.llcontext.VoidType(),
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

func (c *Context) sigType(sig *types.Signature) (t llvm.Type, err error) {
	params := make([]llvm.Type, sig.Params().Len())
	for i := 0; i < sig.Params().Len(); i++ {
		param, err := c.genType(sig.Params().At(i).Type())
		if err != nil {
			return t, err
		}
		params[i] = param
	}
	result, err := c.genType(sig.Result())
	if err != nil {
		return t, err
	}
	return llvm.FunctionType(result, params, false), nil
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
	case *types.Signature:
		return c.sigType(typ)
	default:
		panic(fmt.Sprintf("unimplemented yet: %T", typ))
	}
}
