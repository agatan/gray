package types

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Scope represents literary scope.
// Each function has its own scope and each block also does.
// Scope should be created via NewScope.
type Scope struct {
	parent *Scope
	elems  map[string]Object
	name   string // for debug info. (may be "")
}

// NewScope creates a new root scope.
func NewScope(name string) *Scope {
	return &Scope{parent: builtinScope, elems: map[string]Object{}, name: name}
}

// Parent returns the scope's surrounding scope (if exists).
func (s *Scope) Parent() *Scope { return s.parent }

// Names returns a slice of object names.
func (s *Scope) Names() []string {
	ss := make([]string, 0, len(s.elems))
	for name := range s.elems {
		ss = append(ss, name)
	}
	sort.Strings(ss)
	return ss
}

// Lookup returns the object in the scope s with the given name (if exists)
func (s *Scope) Lookup(name string) Object {
	return s.elems[name]
}

// LookupParent returns the object in the scope s with the given name (if exists).
// In contrast of Lookup, LookupParent look up the object across parents of Scope s.
func (s *Scope) LookupParent(name string) Object {
	for ; s != nil; s = s.parent {
		if obj, ok := s.elems[name]; ok {
			return obj
		}
	}
	return nil
}

// Insert attempts to insert an object into the scope.
// If s already contains any object with the same name,
// Insert inserts obj and returns the old object.
// Otherwise, Insert returns nil.
func (s *Scope) Insert(obj Object) Object {
	name := obj.Name()
	old := s.elems[name]
	if s.elems == nil {
		s.elems = map[string]Object{}
	}
	s.elems[name] = obj
	if obj.Scope() == nil {
		obj.setScope(s)
	}
	return old
}

// Dump dumps the contents of scope s.
func (s *Scope) Dump(w io.Writer, depth int, verbose bool) {
	const ind = ".   "
	indent := strings.Repeat(ind, depth)
	fmt.Fprintf(w, "%s%s scope {", indent, s.name)
	if len(s.elems) == 0 {
		fmt.Fprintf(w, "}\n")
		return
	}

	fmt.Fprintln(w)
	nindent := indent + ind
	for _, name := range s.Names() {
		if verbose {
			fmt.Fprintf(w, "%s%s(%s)\n", nindent, name, s.elems[name].Type())
		} else {
			fmt.Fprintf(w, "%s%s\n", nindent, name)
		}
	}

	fmt.Fprintf(w, "%s}\n", indent)
}

// String returns a string representation of the scope (for debug)
func (s *Scope) String() string {
	var buf bytes.Buffer
	s.Dump(&buf, 0, false)
	return buf.String()
}
