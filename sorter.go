package radix

// SortingTechnique is the technique used
// to sort the tree.
type SortingTechnique uint8

const (
	// AscLabelSort is the value for sorting
	// the tree's edges by label in ascending order.
	AscLabelSort SortingTechnique = iota
	// DescLabelSort is the value for sorting
	// the tree's edges by label in descending order.
	DescLabelSort
	// PrioritySort is the value for sorting
	// the tree's edges by the priority of their nodes.
	PrioritySort
)

type sorter struct {
	n      *Node
	st     SortingTechnique
	binary bool
}

func (s *sorter) Len() int {
	return len(s.n.edges)
}

func (s *sorter) Less(i, j int) bool {
	n := s.n
	if s.st == AscLabelSort {
		if s.binary {
			return i < j
		}
		return n.edges[i].label < n.edges[j].label
	}
	if s.st == DescLabelSort {
		if s.binary {
			return i > j
		}
		return n.edges[i].label > n.edges[j].label
	}
	return n.edges[i].n != nil &&
		n.edges[j].n != nil &&
		n.edges[i].n.priority > n.edges[j].n.priority
}

func (s *sorter) Swap(i, j int) {
	s.n.edges[i], s.n.edges[j] = s.n.edges[j], s.n.edges[i]
}
