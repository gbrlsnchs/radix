package radix

import "sort"

// Node is a node of a radix tree.
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

// Priority returns the node's priority.
func (n *Node) Priority() int {
	return n.priority
}

func (n *Node) addBinary(label string, v interface{}) (size, length int) {
	val := make([]byte, 0)
	for i := range label {
		for j := uint8(8); j > 0; j-- {
			bbit := bit(j, label[i])
			if bbit == 0 {
				val = append(val, '0')
			} else {
				val = append(val, '1')
			}
			done := i == len(label)-1 && j == 1
			if e := n.edges[bbit]; e != nil {
				if done {
					e.n.Value = v
					return
				}
				goto next
			}
			n.edges[bbit] = &edge{
				n: &Node{
					depth: n.depth + 1,
					edges: make([]*edge, 2),
				},
			}
			if done {
				n.edges[bbit].n.Value = v
			}
			size++
			length++
		next:
			n = n.edges[bbit].n
		}
	}
	return
}

func (n *Node) delBinary(label string) (d int) {
	for i := range label {
		for j := uint8(8); j > 0; j-- {
			bbit := bit(j, label[i])
			done := i == len(label)-1 && j == 1
			if e := n.edges[bbit]; e != nil {
				d++
				if done && n.IsLeaf() { // only delete if node is leaf, otherwise it would break the tree
					n.edges = make([]*edge, 2)
					return
				}
			}
			return
		}
	}
	return
}

func (n *Node) getBinary(label string) *Node {
	for i := range label {
		for j := uint8(8); j > 0; j-- {
			bbit := bit(j, label[i])
			done := i == len(label)-1 && j == 1
			if e := n.edges[bbit]; e != nil {
				if done {
					return e.n
				}
				n = e.n
				continue
			}
			return nil
		}
	}
	return nil
}

func (n *Node) incrDepth() {
	n.depth++
	for _, e := range n.edges {
		e.n.incrDepth()
	}
}

// sort sorts the node and its children recursively.
func (n *Node) sort(st SortingTechnique, binary bool) {
	s := &sorter{
		n:      n,
		st:     st,
		binary: binary,
	}
	sort.Sort(s)
	for _, e := range n.edges {
		e.n.sort(st, binary)
	}
}

func (n *Node) split() *Node {
	c := *n
	c.incrDepth()
	return &c
}

func (n *Node) writeTo(bd *builder) {
	for i, e := range n.edges {
		e.writeTo(bd, []bool{i == len(n.edges)-1})
	}
}
