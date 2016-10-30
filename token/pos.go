package token

// Position has line and column infomation.
type Position struct {
	Line   int
	Column int
}

// Pos is an interface that provides store and fetch source locations.
type Pos interface {
	Position() Position
	SetPosition(Position)
}

// PosImpl provides default implementations for Pos.
type PosImpl struct {
	pos Position
}

// Position returns the position of the AST node.
func (x *PosImpl) Position() Position {
	return x.pos
}

// SetPosition is a function to specify position of the AST node.
func (x *PosImpl) SetPosition(pos Position) {
	x.pos = pos
}
