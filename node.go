package patricia

import (
	"bytes"
	"sort"
	"strings"
)

// Node is a node of a PATRICIA tree.
type Node struct {
	// Value is a value of any type held by a node.
	Value    interface{}
	edges    []*edge
	priority int
	depth    int
}

// Depth returns the node's depth.
func (n *Node) Depth() int {
	return n.depth
}

// IsLeaf returns whether the node is a leaf.
func (n *Node) IsLeaf() bool {
	return len(n.edges) == 0
}

// Len returns the number of edges the node has.
func (n *Node) Len() int {
	return len(n.edges)
}

// Less compares two nodes for sorting based on their priority.
func (n *Node) Less(i, j int) bool {
	return n.edges[i].node != nil &&
		n.edges[j].node != nil &&
		n.edges[i].node.priority > n.edges[j].node.priority
}

// Priority returns the node's priority.
func (n *Node) Priority() int {
	return n.priority
}

// Swap swaps two edges from the node.
func (n *Node) Swap(i, j int) {
	n.edges[i], n.edges[j] = n.edges[j], n.edges[i]
}

// buffer returns a pointer to a bytes.Buffer containing
// a subtree structure plus, if debug is truthy, its metadata.
func (n *Node) buffer(debug bool) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}

	for i, e := range n.edges {
		var nbuf *bytes.Buffer
		nbuf, err := e.buffer(debug, []bool{i == len(n.edges)-1})

		if err != nil {
			return nil, err
		}

		_, err = nbuf.WriteTo(buf)

		if err != nil {
			return nil, err
		}
	}

	return buf, nil
}

// child returns a child of the node.
func (n *Node) child(val interface{}) *Node {
	c := &Node{
		Value: val,
		depth: n.depth,
	}

	return c
}

// clone returns a clone of the node.
func (n *Node) clone() *Node {
	c := &Node{
		Value: n.Value,
		edges: n.edges,
		depth: n.depth,
	}

	return c
}

// count counts both the node's and its children's
// depth and priority recursively.
func (n *Node) count(depth int) int {
	n.priority = 0
	n.depth = depth

	if n.Value != nil {
		n.priority++
	}

	for _, e := range n.edges {
		n.priority += e.node.count(n.depth + 1)
	}

	return n.priority
}

// next returns the next edge to be traversed and, if there's a placeholder,
// may also return a map of named parameters.
func (n *Node) next(s string, child bool, ph, delim rune) (*edge, map[string]string) {
	for _, e := range n.edges {
		if i := strings.IndexRune(e.label, ph); i >= 0 && ph != delim {
			var params map[string]string
			lfound := 0
			sfound := 0

			for i >= 0 {
				if !strings.HasPrefix(s[sfound:], e.label[lfound:lfound+i]) {
					return nil, nil
				}

				sfound += len(e.label[lfound : lfound+i])
				sdelim := strings.IndexRune(s[sfound:], delim)

				if sdelim < 0 {
					sdelim = len(s[sfound:])
				}

				lfound += len(e.label[lfound : lfound+i])
				ldelim := strings.IndexRune(e.label[lfound:], delim)

				if ldelim < 0 {
					ldelim = len(e.label[lfound:])
				}

				if params == nil {
					params = make(map[string]string)
				}

				params[e.label[lfound+1:lfound+ldelim]] = s[sfound : sfound+sdelim]
				lfound += len(e.label[lfound : lfound+ldelim])
				sfound += len(s[sfound : sfound+sdelim])
				i = strings.IndexRune(e.label[lfound:], ph)
			}

			if e.label[lfound:] != s[sfound:] &&
				strings.IndexRune(e.label[lfound:], ph) != 0 &&
				len(e.label[lfound:]) > 0 {
				return nil, nil
			}

			return e, params
		}

		hasPrefix := strings.HasPrefix(s, e.label)

		if child {
			hasPrefix = strings.HasPrefix(e.label, s)
		}

		if hasPrefix {
			return e, nil
		}
	}

	return nil, nil
}

// sort sorts the node and its children recursively.
func (n *Node) sort() {
	sort.Sort(n)

	for _, e := range n.edges {
		e.node.sort()
	}
}
