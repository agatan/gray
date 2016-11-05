package ast

// Node is an interface of AST nodes.
type Node interface {
	ID() uint
	SetID(id uint)
}

// NodeImpl holds common AST node informations like ID.
type NodeImpl struct {
	id uint
}

// ID returns the id of Node n.
func (n *NodeImpl) ID() uint {
	if n.id == 0 {
		panic("internal error: Node ID is uninitialized")
	}
	return n.id
}

// SetID sets the id of Node n with the given id.
func (n *NodeImpl) SetID(id uint) { n.id = id }
