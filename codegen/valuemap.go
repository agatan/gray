package codegen

import "llvm.org/llvm/bindings/go/llvm"

// ValueMap holds associative informations between identifiers and thier values.
type ValueMap struct {
	parent *ValueMap
	vmap   map[string]llvm.Value
}

// NewValueMap returns a new ValueMap instance.
func NewValueMap(parent *ValueMap) *ValueMap {
	return &ValueMap{
		parent: parent,
		vmap:   map[string]llvm.Value{},
	}
}

// Insert inserts an identifier and its value.
func (vm *ValueMap) Insert(name string, value llvm.Value) {
	vm.vmap[name] = value
}

// LookupParent look up the identifier value across the parents.
func (vm *ValueMap) LookupParent(name string) llvm.Value {
	if v, ok := vm.vmap[name]; ok {
		return v
	}
	return vm.parent.LookupParent(name)
}

func (c *Context) enterNewScope() {
	old := c.valuemap
	c.valuemap = NewValueMap(old)
}

func (c *Context) exitScope() {
	c.valuemap = c.valuemap.parent
}
