package types

// Type provides an interface of types' type.
type Type interface {
	typ()
	String() string
}

// BasicKind describes the kind of basic types.
type BasicKind int

//go:generate stringer -type=BasicKind
const (
	Invalid BasicKind = iota

	Unit
	Bool
	Int
	String
)

// Basic represent basic types.
type Basic struct {
	Kind BasicKind
}

func (*Basic) typ() {}
func (b *Basic) String() string {
	return b.Kind.String()
}
