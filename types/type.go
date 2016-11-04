package types

import (
	"bytes"
	"fmt"
)

// Type provides an interface of types' type.
type Type interface {
	typ()
	String() string // for debugging
}

// BasicKind describes the kind of basic types.
type BasicKind int

const (
	Invalid BasicKind = iota

	Unit
	Bool
	Int
	String
)

// Basic represent basic types.
type Basic struct {
	kind BasicKind
	name string
}

func (*Basic) typ() {}
func (b *Basic) String() string {
	return b.name
}
func (b *Basic) Kind() BasicKind { return b.kind }
func (b *Basic) Name() string    { return b.name }

// Vars represent a sequence of *Var s.
type Vars struct {
	vars []*Var
}

// NewVars returns a new sequence of variables.
func NewVars(x ...*Var) *Vars {
	return &Vars{vars: x}
}

// Len returns the number of variables.
func (v *Vars) Len() int {
	return len(v.vars)
}

// At returns the i'th variable.
func (v *Vars) At(i int) *Var { return v.vars[i] }

// String returns a string representation of Vars v.
func (vs *Vars) String() string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "(")
	for i, v := range vs.vars {
		if i != 0 {
			fmt.Fprintf(buf, ", ")
		}
		if v.Name() == "" {
			fmt.Fprintf(buf, v.Type().String())
		} else {
			fmt.Fprintf(buf, "%s: %s", v.Name(), v.Type().String())
		}
	}
	fmt.Fprintf(buf, ")")
	return buf.String()
}

// Signature represents a function type.
type Signature struct {
	scope  *Scope
	params *Vars
	result Type
}

func (*Signature) typ() {}
func (s *Signature) String() string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "%s -> %s", s.params.String(), s.result.String())
	return buf.String()
}

// NewSignature returns a new function type.
func NewSignature(scope *Scope, params *Vars, result Type) *Signature {
	return &Signature{scope: scope, params: params, result: result}
}

// Params returns the parameters of the signature.
func (s *Signature) Params() *Vars { return s.params }

// Result returns the return type of the signature.
func (s *Signature) Result() Type { return s.result }

// GenericType represents a generic type (e.g. Ref<Int>).
type GenericType struct {
	instances [][]Type // a set of instantiate args
	name      string
	params    []*TypeName
}

// NewGenericType returns a new generic type.
func NewGenericType(name string, params []*TypeName) *GenericType {
	return &GenericType{instances: nil, name: name, params: params}
}

func (*GenericType) typ()             {}
func (g *GenericType) String() string { return fmt.Sprintf("%s<with %d types>", g.name, len(g.params)) }

// Name returns the name of GenericType g.
func (g *GenericType) Name() string { return g.name }

// Params returns the parameters of GenericType g.
func (g *GenericType) Params() []*TypeName { return g.params }

// InstType returns a new generic type's instance.
type InstType struct {
	base *GenericType
	args []Type
}

// NewInstType returns a new instance of generic type g.
func NewInstType(g *GenericType, args []Type) *InstType {
	return &InstType{base: g, args: args}
}

func (i *InstType) typ() {}
func (i *InstType) String() string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "%s<", i.base.Name())
	for i, a := range i.args {
		if i != 0 {
			fmt.Fprintf(buf, ", ")
		}
		fmt.Fprintf(buf, a.String())
	}
	fmt.Fprintf(buf, ">")
	return buf.String()
}

// Base returns a base type of InstType i.
func (i *InstType) Base() *GenericType { return i.base }

// Args returns args of InstType i.
func (i *InstType) Args() []Type { return i.args }
