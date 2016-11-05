package types

import "github.com/agatan/gray/ast"

// TypeMap holds an associative map between Node ID and its type.
type TypeMap struct {
	typemap []Type
}

// NewTypeMap returns a new TypeMap instance.
func NewTypeMap() *TypeMap {
	return &TypeMap{typemap: make([]Type, 1024)}
}

// extendTo allocates enough buffer to record the id.
func (t *TypeMap) extendTo(id uint) {
	for i := uint(len(t.typemap)) - 1; i < id; i++ {
		t.typemap = append(t.typemap, nil)
	}
}

// Record records the type information with the Node n.
func (t *TypeMap) Record(n ast.Node, typ Type) {
	if uint(len(t.typemap)) < n.ID()+1 {
		t.extendTo(n.ID())
	}
	t.typemap[n.ID()] = typ
}

// Type returns the type of given Node.
func (t *TypeMap) Type(n ast.Node) Type {
	return t.typemap[n.ID()]
}
