package codegen

import (
	"github.com/agatan/gray/ast"
	"github.com/agatan/gray/types"

	"llvm.org/llvm/bindings/go/llvm"
)

// Context represent a context of code generation.
type Context struct {
	llcontext       llvm.Context
	llbuilder       llvm.Builder
	llmodule        llvm.Module
	lltarget        llvm.Target
	lltargetMachine llvm.TargetMachine
	lltargetData    llvm.TargetData

	toplevelScope *types.Scope
	typemap       *types.TypeMap
}

// NewContext returns a new Context instance.
func NewContext(modname string, toplevelScope *types.Scope, typemap *types.TypeMap) (*Context, error) {
	if err := llvm.InitializeNativeTarget(); err != nil {
		return nil, err
	}
	triple := llvm.DefaultTargetTriple()
	target, err := llvm.GetTargetFromTriple(triple)
	if err != nil {
		return nil, err
	}
	tm := target.CreateTargetMachine(triple, "", "", llvm.CodeGenLevelNone, llvm.RelocDefault, llvm.CodeModelDefault)
	return &Context{
		llcontext:       llvm.GlobalContext(),
		llbuilder:       llvm.NewBuilder(),
		llmodule:        llvm.GlobalContext().NewModule(modname),
		lltarget:        target,
		lltargetMachine: tm,
		lltargetData:    tm.CreateTargetData(),

		toplevelScope: toplevelScope,
		typemap:       typemap,
	}, nil
}

// Dispose is a destructor method for Context.
func (c *Context) Dispose() {
	c.llcontext.Dispose()
}

// Generate generates llvm IRs from ast.Decls.
func (c *Context) Generate(ds []ast.Decl) (llvm.Module, error) {
	if err := c.forwardDecls(c.toplevelScope, ds); err != nil {
		return c.llmodule, err
	}
	return c.llmodule, nil
}
