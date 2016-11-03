package types

// Type provides an interface of types' type.
type Type interface {
	typ()
	String() string
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
