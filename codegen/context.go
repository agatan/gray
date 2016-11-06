package codegen

import "llvm.org/llvm/bindings/go/llvm"

// Context represent a context of code generation.
type Context struct {
	llcontext llvm.Context
	llbuilder llvm.Builder
}

// NewContext returns a new Context instance.
func NewContext() *Context {
	return &Context{
		llcontext: llvm.GlobalContext(),
		llbuilder: llvm.NewBuilder(),
	}
}
