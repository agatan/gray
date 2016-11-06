package codegen

import (
	"io/ioutil"

	"github.com/agatan/gray/ast"
	"github.com/agatan/gray/types"

	"llvm.org/llvm/bindings/go/llvm"
)

// Context represent a context of code generation.
type Context struct {
	moduleName string

	llcontext       llvm.Context
	llbuilder       llvm.Builder
	llmodule        llvm.Module
	lltarget        llvm.Target
	lltargetMachine llvm.TargetMachine
	lltargetData    llvm.TargetData

	toplevelScope *types.Scope
	typemap       *types.TypeMap

	basicTypes []llvm.Type
	valuemap   *ValueMap
	typenames  map[string]llvm.Type
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
		moduleName:      modname,
		llcontext:       llvm.GlobalContext(),
		llbuilder:       llvm.NewBuilder(),
		llmodule:        llvm.GlobalContext().NewModule(modname),
		lltarget:        target,
		lltargetMachine: tm,
		lltargetData:    tm.CreateTargetData(),

		toplevelScope: toplevelScope,
		typemap:       typemap,

		valuemap: NewValueMap(nil),
	}
	ctx.defBasicTypes()
	ctx.defBuiltinFunctions()
	return ctx, nil
}

// Dispose is a destructor method for Context.
func (c *Context) Dispose() {
	c.llcontext.Dispose()
}

// Generate generates llvm IRs from ast.Decls.
func (c *Context) Generate(ds []ast.Decl) (llvm.Module, error) {
	c.forwardDecls(c.toplevelScope, ds)
	for _, d := range ds {
		if err := c.genDecl(d); err != nil {
			return c.llmodule, err
		}
	}
	return c.llmodule, nil
}

// EmitObject generates compiled object file for given AST.
func (c *Context) EmitObject(outname string, llmod llvm.Module) error {
	err := c.defGrayMain()
	if err != nil {
		return err
	}
	membuf, err := c.lltargetMachine.EmitToMemoryBuffer(llmod, llvm.ObjectFile)
	if err != nil {
		return err
	}
	defer membuf.Dispose()
	err = ioutil.WriteFile(outname, membuf.Bytes(), 0666)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) currentScope() *types.Scope {
	return c.toplevelScope
}
