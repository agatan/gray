package types

// Object provides an interface of gray's semantical objects.
type Object interface {
	Scope() *Scope
	Name() string
	Type() Type

	setScope(*Scope)
}

type object struct {
	scope *Scope
	name  string
	typ   Type
}

func (obj *object) Scope() *Scope     { return obj.scope }
func (obj *object) Name() string      { return obj.name }
func (obj *object) Type() Type        { return obj.typ }
func (obj *object) setScope(s *Scope) { obj.scope = s }

// Var represent a declared variable (including function parameters)
type Var struct {
	object
}

// NewVar creates a new variable object.
func NewVar(name string, typ Type) *Var {
	return &Var{object: object{scope: nil, name: name, typ: typ}}
}

// TypeName represent a declared type.
type TypeName struct {
	object
}

// NewTypeName creates a new typename object.
func NewTypeName(name string, typ Type) *TypeName {
	return &TypeName{object: object{scope: nil, name: name, typ: typ}}
}

// Func represent a declared function
type Func struct {
	object
}

// NewFunc creates a new function object.
func NewFunc(name string, typ Type) *Func {
	return &Func{object: object{scope: nil, name: name, typ: typ}}
}
