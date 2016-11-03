package types

import "strings"

// This file initializes the global builtin object scope.

var builtinScope *Scope

var BuiltinTypes = []*Basic{
	Invalid: {kind: Invalid, name: "invalid type"},

	Unit:   {kind: Unit, name: "Unit"},
	Bool:   {kind: Bool, name: "Bool"},
	Int:    {kind: Int, name: "Int"},
	String: {kind: String, name: "String"},
}

func defBuiltinTypes() {
	for _, t := range BuiltinTypes {
		defObject(NewTypeName(t.Name(), t))
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
	defBuiltinTypes()
}
