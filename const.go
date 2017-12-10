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
