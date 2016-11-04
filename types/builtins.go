package types

import "strings"

// This file initializes the global builtin object scope.

var builtinScope *Scope

var BasicTypes = []*Basic{
	Invalid: {kind: Invalid, name: "invalid type"},

	Unit:   {kind: Unit, name: "Unit"},
	Bool:   {kind: Bool, name: "Bool"},
	Int:    {kind: Int, name: "Int"},
	String: {kind: String, name: "String"},
}

func defBasicTypes() {
	for _, t := range BasicTypes {
		defObject(NewTypeName(t.Name(), t))
	}
}

type builtinType int

const (
	invalidBuiltinType builtinType = iota
	refType
)

var builtinGenericTypes = []*GenericType{
	invalidBuiltinType: nil,
	refType:            NewGenericType("Ref", NewTypeName("T", nil)),
}

func defBuiltinGenericTypes() {
	for _, t := range builtinGenericTypes {
		if t == nil {
			continue
		}
		defObject(NewTypeName(t.Name(), t))
	}
}

const (
	BuiltinPrintInt = "print_int"
)

var builtinFunctions = []struct {
	name string
	sig  *Signature
}{
	{BuiltinPrintInt, NewSignature(nil, NewVars(NewVar("x", BasicTypes[Int])), BasicTypes[Unit])},
}

func defBuiltinFunctions() {
	for _, t := range builtinFunctions {
		defObject(NewFunc(t.name, t.sig))
	}
}

func defObject(obj Object) {
	name := obj.Name()
	if strings.Contains(name, " ") {
		return // invalid object (e.g. Invalid type)
	}
	if builtinScope.Insert(obj) != nil {
		panic("internal error: builtin object's double declaration")
	}
}

func init() {
	builtinScope = &Scope{name: "builtin"}
	defBasicTypes()
	defBuiltinGenericTypes()
	defBuiltinFunctions()
}
