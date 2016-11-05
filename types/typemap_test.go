package types

import (
	"testing"

	"github.com/agatan/gray/ast"
)

func TestRecordAndType(t *testing.T) {
	typemap := NewTypeMap()
	typ := NewBangType()
	tests := []uint{1, 10, 1024, 1025, 2048, 2049}
	nodes := make([]ast.Node, len(tests))
	for i, test := range tests {
		node := &ast.NodeImpl{}
		node.SetID(test)
		nodes[i] = node
	}
	for _, node := range nodes {
		typemap.Record(node, typ)
	}
	for i, node := range nodes {
		if typemap.Type(node) == nil {
			t.Errorf("#%d: TypeMap doesn't have the type of Node(ID: %d)", i, node.ID())
		}
	}
}
