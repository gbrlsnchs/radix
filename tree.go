package radix

import (
	"strings"
	"sync"

	"github.com/fatih/color"
)

const (
	// Tsafe activates thread safety.
	Tsafe = 1 << iota
	// Tdebug adds more information to the tree's string representation.
	Tdebug
	// Tbinary uses a binary PATRICIA tree instead of a prefix Radix tree.
	Tbinary
	// Tnocolor disables colorful output.
	Tnocolor
)

// Tree is a radix tree.
type Tree struct {
	root        *Node
	length      int // total number of nodes
	size        int // total byte size
	safe        bool
	binary      bool
	placeholder byte
	delim       byte
	mu          *sync.RWMutex
	bd          *builder
}

// New creates a named radix tree with a single node (its root).
func New(flags int) *Tree {
	tr := &Tree{
		root:   &Node{},
		length: 1,
	}
	if flags&Tbinary > 0 {
		tr.binary = true
		tr.root.edges = make([]*edge, 2) // create two empty edges
	}
	if flags&Tsafe > 0 {
		tr.mu = &sync.RWMutex{}
		tr.safe = true
	}
	tr.bd = &builder{
		Builder: &strings.Builder{},
		debug:   flags&Tdebug > 0,
	}
	tr.bd.colors[colorRed] = color.New(color.FgRed)
	tr.bd.colors[colorGreen] = color.New(color.FgGreen)
	tr.bd.colors[colorMagenta] = color.New(color.FgMagenta)
	tr.bd.colors[colorBold] = color.New(color.Bold)
	for _, c := range tr.bd.colors {
		if flags&Tnocolor > 0 {
			c.DisableColor()
		}
	}
	return tr
}

// Add adds a new node to the tree.
func (tr *Tree) Add(label string, v interface{}) {
	// No empty strings or interfaces allowed.
	if label == "" || v == nil {
		return
	}
	if tr.safe {
		defer tr.mu.Unlock()
		tr.mu.Lock()
	}
	tnode := tr.root
	if tr.binary {
		nn := tnode.addBinary(label, v)
		tr.length += nn
		tr.size += nn / 8
		return
	}
	for {
		var next *edge
		var slice string
		for _, edge := range tnode.edges {
			var found int
			slice = edge.label
			for i := range slice {
				if i < len(label) && slice[i] == label[i] {
					found++
					continue
				}
				break
			}
			if found > 0 {
				label = label[found:]
				slice = slice[found:]
				next = edge
				break
			}
		}
		if next != nil {
			tnode = next.n
			tnode.priority++
			// Match the whole word.
			if len(label) == 0 {
				// The label is exactly the same as the edge's label,
				// so just replace its node's value.
				//
				// Example:
				// 	(root) -> tnode("tomato", v1)
				// 	becomes
				// 	(root) -> tnode("tomato", v2)
				if len(slice) == 0 {
					tnode.Value = v
					return
				}
				// The label is a prefix of the edge's label.
				//
				// Example:
				// 	(root) -> tnode("tomato", v1)
				// 	then add "tom"
				// 	(root) -> ("tom", v2) -> ("ato", v1)
				next.label = next.label[:len(next.label)-len(slice)]
				tnode.edges = []*edge{
					&edge{
						label: slice,
						n:     tnode.split(),
					},
				}
				tnode.Value = v
				tr.length++
				tr.size += len(slice)
				return
			}
			// Add a new node but break its parent into prefix and
			// the remaining slice as a new edge.
			//
			// Example:
			// 	(root) -> ("tomato", v1)
			// 	then add "tornado"
			// 	(root) -> ("to", nil) -> ("mato", v1)
			// 	                      +> ("rnado", v2)
			if len(slice) > 0 {
				c := tnode.split()
				c.priority--
				tnode.edges = []*edge{
					&edge{ // the suffix that is split into a new node
						label: slice,
						n:     c,
					},
					&edge{ // the new node
						label: label,
						n: &Node{
							Value:    v,
							depth:    tnode.depth + 1,
							priority: 1,
						},
					},
				}
				next.label = next.label[:len(next.label)-len(slice)]
				tnode.Value = nil
				tr.length += 2
				tr.size += len(label)
				return
			}
			continue
		}
		tnode.edges = append(tnode.edges, &edge{
			label: label,
			n: &Node{
				Value:    v,
				depth:    tnode.depth + 1,
				priority: 1,
			},
		})
		tr.length++
		tr.size += len(label)
		return
	}
}

// Del deletes a node.
//
// If a parent node that holds no value ends up holding only one edge
// after a deletion of one of its edges, it gets merged with the remaining edge.
func (tr *Tree) Del(label string) {
	if string(label) == "" {
		return
	}
	if tr.safe {
		defer tr.mu.Unlock()
		tr.mu.Lock()
	}
	tnode := tr.root
	if tr.binary {
		del := tnode.delBinary(label)
		tr.length--
		tr.size = (tr.size*8 - del) / 8
		return
	}
	var edgex int
	var parent *edge
	var ptrs []*int
	for tnode != nil && label != "" {
		var next *edge
		// Look for exact matches.
		for i, e := range tnode.edges {
			if strings.HasPrefix(label, e.label) {
				next = e
				edgex = i
				break
			}
		}
		if next != nil {
			tnode = next.n
			label = label[len(next.label):]
			ptrs = append(ptrs, &tnode.priority)
			// While not the exact match, set the tnode's parent.
			if label != "" {
				parent = next
			}
			continue
		}
		// No matches.
		parent = nil
		tnode = nil
	}
	if tnode != nil {
		pnode := tr.root // in case label matched in the first try
		if parent != nil {
			pnode = parent.n
		}
		// Decrement the priority of upper nodes.
		done := make(chan struct{})
		if tnode.Value != nil {
			go func() {
				for _, p := range ptrs {
					*p--
				}
				close(done)
			}()
		}
		// Merge tnode's edges with the parent's.
		pnode.edges = append(pnode.edges, tnode.edges...)
		// Remove tnode from the parent, leaving only its edges behind.
		pnode.edges = append(pnode.edges[:edgex], pnode.edges[edgex+1:]...)
		// When only one edge remained in pnode and its value is nil, they can be merged.
		if len(pnode.edges) == 1 && pnode.Value == nil && parent != nil {
			e := pnode.edges[0]
			parent.label += e.label
			pnode.Value = e.n.Value
			pnode.edges = e.n.edges
			tr.length--
		}
		tr.length--
		if tnode.Value != nil {
			<-done
		}
	}
}

// Get retrieves a node.
func (tr *Tree) Get(label string) (*Node, map[string]string) {
	if label == "" {
		return nil, nil
	}
	if tr.safe {
		defer tr.mu.RUnlock()
		tr.mu.RLock()
	}
	tnode := tr.root
	if tr.binary {
		return tnode.getBinary(label), nil
	}
	var params map[string]string
	for tnode != nil && label != "" {
		var next *edge
	walk:
		for _, edge := range tnode.edges {
			slice := edge.label
			for {
				phIndex := len(slice)
				// Check if there are any placeholders.
				// If there are none, then use the whole word for comparison.
				if i := strings.IndexByte(slice, tr.placeholder); i >= 0 {
					phIndex = i
				}
				prefix := slice[:phIndex]
				// If "slice" (until placeholder) is not prefix of
				// "label", then keep walking.
				if !strings.HasPrefix(label, prefix) {
					continue walk
				}
				label = label[len(prefix):]
				// If "slice" is the whole label,
				// then the match is complete and the algorithm
				// is ready to go to the next edge.
				if len(prefix) == len(slice) {
					next = edge
					break walk
				}
				// Check whether there is a delimiter.
				// If there isn'tr, then use the whole world as parameter.
				var delimIndex int
				slice = slice[phIndex:]
				if delimIndex = strings.IndexByte(slice[1:], tr.delim) + 1; delimIndex <= 0 {
					delimIndex = len(slice)
				}
				key := slice[1:delimIndex] // remove the placeholder from the map key
				slice = slice[delimIndex:]
				if delimIndex = strings.IndexByte(label[1:], tr.delim) + 1; delimIndex <= 0 {
					delimIndex = len(label)
				}
				if params == nil {
					params = make(map[string]string)
				}
				params[key] = label[:delimIndex]
				label = label[delimIndex:]
				if slice == "" && label == "" {
					next = edge
					break walk
				}
			}
		}
		if next != nil {
			tnode = next.n
			continue
		}
		tnode = nil
	}
	return tnode, params
}

// Len returns the total numbers of nodes,
// including the tree's root.
func (tr *Tree) Len() int {
	if tr.safe {
		defer tr.mu.RUnlock()
		tr.mu.RLock()
	}
	return tr.length
}

// SetBoundaries sets a placeholder and a delimiter for
// the tree to be able to search for named labels.
func (tr *Tree) SetBoundaries(placeholder, delim byte) {
	tr.placeholder = placeholder
	tr.delim = delim
}

// Size returns the total byte size stored in the tree.
func (tr *Tree) Size() int {
	return tr.size
}

// Sort sorts the tree nodes and its children recursively
// according to their priority lengther.
func (tr *Tree) Sort(st SortingTechnique) {
	if tr.safe {
		defer tr.mu.Unlock()
		tr.mu.Lock()
	}
	tr.root.sort(st, tr.binary)
}

// String returns a string representation of the tree structure.
func (tr *Tree) String() string {
	if tr.safe {
		defer tr.mu.RUnlock()
		tr.mu.RLock()
	}
	tr.bd.Reset()
	tr.bd.colors[colorBold].Fprint(tr.bd, "\n.")
	if tr.bd.debug {
		mag := tr.bd.colors[colorMagenta]
		mag.Fprintf(tr.bd, " (%d node", tr.length)
		if tr.length != 1 {
			mag.Fprint(tr.bd, "s") // avoid writing "1 nodes"
		}
		mag.Fprint(tr.bd, ")")
	}
	tr.bd.WriteByte('\n')
	if !tr.binary {
		// TODO: implement binary string representation
		tr.root.writeTo(tr.bd)
	}
	return tr.bd.String()
}
