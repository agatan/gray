package codegen

import "llvm.org/llvm/bindings/go/llvm"

// Context represent a context of code generation.
type Context struct {
	llcontext       llvm.Context
	llbuilder       llvm.Builder
	lltarget        llvm.Target
	lltargetMachine llvm.TargetMachine
	lltargetData    llvm.TargetData
}

// NewContext returns a new Context instance.
func NewContext() (*Context, error) {
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
		lltarget:        target,
		lltargetMachine: tm,
		lltargetData:    tm.CreateTargetData(),
	}, nil
}
