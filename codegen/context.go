package codegen

import (
	"io/ioutil"

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

	basicTypes []llvm.Type
}

// NewContext returns a new Context instance.
func NewContext(modname string, toplevelScope *types.Scope, typemap *types.TypeMap) (*Context, error) {
	if err := llvm.InitializeNativeTarget(); err != nil {
		return nil, err
	}
	if err := llvm.InitializeNativeAsmPrinter(); err != nil {
		return nil, err
	}
	triple := llvm.DefaultTargetTriple()
	target, err := llvm.GetTargetFromTriple(triple)
	if err != nil {
		return nil, err
	}
	tm := target.CreateTargetMachine(triple, "", "", llvm.CodeGenLevelNone, llvm.RelocDefault, llvm.CodeModelDefault)
	ctx := &Context{
		llcontext:       llvm.GlobalContext(),
		llbuilder:       llvm.NewBuilder(),
		llmodule:        llvm.GlobalContext().NewModule(modname),
		lltarget:        target,
		lltargetMachine: tm,
		lltargetData:    tm.CreateTargetData(),

		toplevelScope: toplevelScope,
		typemap:       typemap,
	}
	ctx.defBasicTypes()
	if err := ctx.defBuiltinFunctions(); err != nil {
		return nil, err
	}
	return ctx, nil
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
	for _, d := range ds {
		if err := c.genDecl(d); err != nil {
			return c.llmodule, err
		}
	}
	return c.llmodule, nil
}

// EmitObject generates compiled object file for given AST.
func (c *Context) EmitObject(ds []ast.Decl) error {
	mod, err := c.Generate(ds)
	membuf, err := c.lltargetMachine.EmitToMemoryBuffer(mod, llvm.ObjectFile)
	if err != nil {
		return err
	}
	defer membuf.Dispose()
	err = ioutil.WriteFile("test.o", membuf.Bytes(), 0666)
	if err != nil {
		return err
	}
	return nil
}
