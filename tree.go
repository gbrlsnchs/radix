package radix

import (
	"bytes"
	"log"
	"strings"
	"sync"

	"github.com/fatih/color"
)

// Tree is a radix tree.
type Tree struct {
	// Safe tells whether the tree's operations
	// should be thread-safe. By default, the tree's
	// not thread-safe.
	Safe bool
	// name is the tree's name.
	name string
	// root is the tree's root.
	root *Node
	// size is the total number of nodes
	// without the tree's root.
	size uint
	// mtx controls the operations' safety.
	mtx *sync.RWMutex
}

// New creates a named radix tree with a single node (its root).
func New(name string) *Tree {
	return &Tree{name: name, root: &Node{}, mtx: &sync.RWMutex{}}
}

// Add adds a new node to the tree.
func (t *Tree) Add(s string, v interface{}) {
	if t.Safe {
		defer t.mtx.Unlock()
		t.mtx.Lock()
	}

	if s == "" || v == nil {
		return
	}

	sfound := 0
	cfound := 0
	tnode := t.root

walk:
	for {
		var next *edge

		for _, e := range tnode.edges {
			cfound = 0
			str := s[sfound:]

			for i := range e.label {
				if cfound < len(str) && e.label[i] == str[cfound] {
					cfound++

					continue
				}

				break
			}

			if cfound > 0 {
				sfound += cfound
				next = e

				break
			}
		}

		if next != nil {
			tnode = next.node

			if v != nil {
				tnode.priority++
			}

			// The conditions below only happen if
			// the whole string has been matched.
			if sfound == len(s) {
				// When the string already exists inside the tree and
				// there's nothing more to add to the latter,
				// only the value is substituted.
				if cfound == len(next.label) {
					tnode.Value = v

					break walk
				}

				// When the string is a prefix of the edge's label,
				// it splits the latter into the prefix and a new child,
				// the remaining label without the prefix.
				tnode.edges = []*edge{
					newEdge(next.label[len(s[sfound-cfound:]):], tnode.clone()),
				}
				tnode.Value = v
				next.label = next.label[:len(s[sfound-cfound:])]
				t.size++

				break walk
			}

			// The string "s" is a splitter or, in other words,
			// it shares a common prefix with an edge and thus
			// splits the edge label into the prefix (parent) and
			// two children, one being the remaning string per se and
			// the other being the part that doesn't match the string.
			if cfound > 0 && cfound < len(next.label) {
				tnode.edges = []*edge{
					newEdge(next.label[cfound:], tnode.clone()),
					newEdge(s[sfound:], tnode.child(v)),
				}

				for _, e := range tnode.edges {
					if e.node.Value != nil {
						e.node.priority = 1
					}
				}

				tnode.Value = nil
				next.label = next.label[:cfound]
				t.size += 2

				break walk
			}

			continue
		}

		// When the string has not been fully matched,
		// it appends a new child to the last node traversed.
		tnode.edges = append(tnode.edges, newEdge(s[sfound:], tnode.child(v)))
		t.size++
	}
}

// Debug prints the tree's structure plus its metadata.
func (t *Tree) Debug() error {
	s, err := t.String(true)

	if err != nil {
		return err
	}

	log.Println(s)

	return nil
}

// Del deletes a node.
//
// If a parent node that holds no value ends up holding only one edge
// after a deletion of one of its edges, it gets merged with the remaining edge.
func (t *Tree) Del(s string) {
	if t.Safe {
		defer t.mtx.Unlock()
		t.mtx.Lock()
	}

	found := 0
	tnode := t.root
	edgeIndex := 0
	var parent *edge
	var priorityPtrs []*int

	for tnode != nil && found < len(s) {
		var next *edge

		for _, e := range tnode.edges {
			if strings.HasPrefix(s[found:], e.label) {
				next = e

				break
			}
		}

		if next != nil {
			for i, e := range tnode.edges {
				if next.label == e.label {
					edgeIndex = i

					break
				}
			}

			tnode = next.node
			found += len(next.label)
			priorityPtrs = append(priorityPtrs, &tnode.priority)

			if found < len(s) {
				parent = next
			}

			continue
		}

		parent = nil
		tnode = nil
	}

	if tnode != nil {
		parentNode := t.root

		if parent != nil {
			parentNode = parent.node
		}

		if tnode.Value != nil {
			for _, p := range priorityPtrs {
				*p--
			}
		}

		parentNode.edges = append(parentNode.edges, tnode.edges...)
		parentNode.edges = append(parentNode.edges[:edgeIndex], parentNode.edges[edgeIndex+1:]...)

		if len(parentNode.edges) == 1 && parentNode.Value == nil && parent != nil {
			parent.label += parentNode.edges[0].label
			parentNode.Value = parentNode.edges[0].node.Value
			parentNode.edges = parentNode.edges[0].node.edges
			t.size--
		}

		t.size--
	}
}

// Get retrieves a node.
func (t *Tree) Get(s string) *Node {
	n, _ := t.get(s, 0, 0)

	return n
}

// GetByRune dynamically retrieves a node based on a placeholder and a delimiter.
// It also returns a map of "named parameters".
func (t *Tree) GetByRune(s string, ph, delim rune) (*Node, map[string]string) {
	return t.get(s, ph, delim)
}

// Print prints the tree's structure.
func (t *Tree) Print() error {
	s, err := t.String(false)

	if err != nil {
		return err
	}

	log.Println(s)

	return nil
}

// Size returns the total numbers of nodes the tree has,
// including the root.
func (t *Tree) Size() uint {
	if t.Safe {
		defer t.mtx.RUnlock()
		t.mtx.RLock()
	}

	return t.size + 1
}

// Sort sorts the tree nodes and its children recursively
// according to their priority counter.
func (t *Tree) Sort(st SortingTechnique) {
	t.root.sort(st)
}

// String returns a string representation of the tree structure.
func (t *Tree) String(debug bool) (string, error) {
	if t.Safe {
		defer t.mtx.RUnlock()
		t.mtx.RLock()
	}

	buf := &bytes.Buffer{}
	green := color.New(color.FgGreen).SprintfFunc()
	magenta := color.New(color.FgMagenta).SprintfFunc()
	bold := color.New(color.Bold).SprintFunc()

	_, err := buf.WriteString(green("\n%s", t.name))

	if err != nil {
		return "", err
	}

	_, err = buf.WriteString(bold("\n."))

	if err != nil {
		return "", err
	}

	if debug {
		_, err = buf.WriteString(magenta(" (%d nodes)", t.size+1))

		if err != nil {
			return "", err
		}
	}

	err = buf.WriteByte('\n')

	if err != nil {
		return "", err
	}

	nbuf, err := t.root.buffer(debug)

	if err != nil {
		return "", err
	}

	_, err = nbuf.WriteTo(buf)

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// get retrieves a node dynamically or not.
func (t *Tree) get(s string, ph, delim rune) (*Node, map[string]string) {
	if t.Safe {
		defer t.mtx.RUnlock()
		t.mtx.RLock()
	}

	sfound := 0
	tnode := t.root

	if tnode.IsLeaf() {
		return nil, nil
	}

	var params map[string]string

	for tnode != nil && sfound < len(s) {
		var next *edge

		for _, e := range tnode.edges {
			lfound := 0

			for lfound < len(e.label) {
				// Checks for a placeholder.
				i := strings.IndexRune(e.label[lfound:], ph)

				// If no placeholder is found,
				// then the limit is the end of the word.
				// Also, if the placeholder equals the delimiter,
				// disregard the label as a named parameter.
				if i < 0 || ph == delim {
					i = len(e.label[lfound:])
				}

				// Checks for a match of the label before the placeholder
				// in the remaining string.
				j := strings.Index(s[sfound:], e.label[lfound:lfound+i])

				// If the label before the placeholder is not a prefix
				// of the string, then the lookup fails.
				if j < 0 {
					break
				}

				// Sums the length of the label slice before the placeholder
				// to the "found" counter of both the label and the string.
				llen := len(e.label[lfound : lfound+i])
				sfound += llen
				lfound += llen

				// If there's no placeholder ahead,
				// move to the next edge traverse.
				if i == len(e.label) {
					next = e

					break
				}

				// Finds where the named parameter's key and value end.
				ldelim := strings.IndexRune(e.label[lfound:], delim)
				sdelim := strings.IndexRune(s[sfound:], delim)

				// If there's no delimiter, then it ends when
				// the label and the string themselves end.
				if ldelim < 0 {
					ldelim = len(e.label[lfound:])
				}

				if sdelim < 0 {
					sdelim = len(s[sfound:])
				}

				if len(e.label) < lfound+1 {
					next = e

					break
				}

				k := e.label[lfound+1 : lfound+ldelim]
				v := s[sfound : sfound+sdelim]

				if params == nil {
					params = make(map[string]string)
				}

				// Adds the named parameter to the "params" map and
				// sums the label's named parameter's length to "lfound" and
				// the parameter's value's length to "sfound".
				params[k] = v
				lfound += len(k) + 1
				sfound += len(v)
			}

			if lfound != len(e.label) {
				continue
			}

			next = e
		}

		if next != nil {
			tnode = next.node

			continue
		}

		tnode = nil
		params = nil
	}

	if sfound < len(s) {
		return nil, nil
	}

	return tnode, params
}
