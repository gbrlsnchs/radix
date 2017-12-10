package radix

import (
	"bytes"
	"sort"
)

// Node is a node of a radix tree.
type Node struct {
	// Value is a value of any type held by a node.
	Value    interface{}
	edges    []*edge
	priority int
	depth    int
	st       SortingTechnique
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
	if n.st == AscLabelSort {
		return n.edges[i].label < n.edges[j].label
	}

	if n.st == DescLabelSort {
		return n.edges[i].label > n.edges[j].label
	}

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
func (n *Node) child(v interface{}) *Node {
	c := &Node{
		Value: v,
		depth: n.depth + 1,
	}

	return c
}

// clone returns a clone of the node.
func (n *Node) clone() *Node {
	c := &Node{
		Value:    n.Value,
		edges:    n.edges,
		depth:    n.depth + 1,
		priority: n.priority - 1,
	}

	if c.Value != nil {
		c.priority++
	}

	c.incrDepth()

	return c
}

func (n *Node) incrDepth() {
	for _, e := range n.edges {
		e.node.depth++

		e.node.incrDepth()
	}
}

// sort sorts the node and its children recursively.
func (n *Node) sort(st SortingTechnique) {
	n.st = st

	sort.Sort(n)

	for _, e := range n.edges {
		e.node.sort(st)
	}
}
