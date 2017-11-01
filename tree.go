package patricia

import (
	"bytes"
	"log"
	"strings"

	"github.com/fatih/color"
)

// Tree is a PATRICIA Tree.
type Tree struct {
	// name is the tree's name.
	name string
	// root is the tree's root.
	root *Node
	// size is the total number of nodes
	// without the tree's root.
	size uint
}

// New creates a named PATRICIA tree with a single node (its root).
func New(name string) *Tree {
	return &Tree{name: name, root: &Node{}}
}

// Add adds a new word to the tree.
func (t *Tree) Add(s string, v interface{}) *Tree {
	found := 0
	tnode := t.root
	child := false
	s = strings.ToLower(s)

walk:
	for {
		next, _ := tnode.next(s[found:], child, 0, 0)

		if next != nil {
			tnode = next.node

			// keep looking
			if !child {
				found += len(next.label)

				if len(s[found:]) > 0 {
					continue
				}

				next.node.Value = v

				break walk
			}

			// in this case, tnode will become a child of a new node with the prefix
			// so deep copy it and clean afterwards
			tnode.edges = []*edge{
				newEdge(next.label[len(s[found:]):], tnode.clone()),
			}
			tnode.Value = v
			next.label = next.label[:len(s[found:])]
			t.size += 2

			break
		}

		if !child {
			child = true

			continue
		}

		// if there are remaining elements,
		// add a new node
		for _, e := range tnode.edges {
			cfound := 0

			for _, c := range e.label {
				if c == rune(s[found:][cfound]) {
					cfound++

					continue
				}

				if cfound > 0 {
					next = e
					tnode = next.node
					tnode.edges = []*edge{
						// clone from parent
						newEdge(next.label[cfound:], tnode.clone()),
						newEdge(s[found+cfound:], tnode.child(v)),
					}
					tnode.Value = nil
					next.label = next.label[:cfound]
					t.size += 2

					break walk
				}

				break
			}
		}

		tnode.edges = append(tnode.edges, newEdge(s[found:], tnode.child(v)))
		t.size++

		break
	}

	_ = t.root.count(0)
	t.root.sort()

	return t
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
func (t *Tree) Del(s string) *Tree {
	found := 0
	tnode := t.root
	edgeIndex := 0
	var parent *edge

	for tnode != nil && found < len(s) {
		next, _ := tnode.next(s[found:], false, 0, 0)

		if next != nil {
			for i, e := range tnode.edges {
				if next.label == e.label {
					edgeIndex = i

					break
				}
			}

			tnode = next.node
			found += len(next.label)

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

		parentNode.edges = append(parentNode.edges, tnode.edges...)
		parentNode.edges = append(parentNode.edges[:edgeIndex], parentNode.edges[edgeIndex+1:]...)

		if len(parentNode.edges) == 1 && parentNode.Value == nil && parent != nil {
			parent.label += parentNode.edges[0].label
			parentNode.Value = parentNode.edges[0].node.Value
			parentNode.edges = parentNode.edges[0].node.edges
			t.size--
		}
	}

	_ = t.root.count(0)
	t.root.sort()
	t.size--

	return t
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
	return t.size + 1
}

// String returns a string representation of the tree structure.
func (t *Tree) String(debug bool) (string, error) {
	buf := &bytes.Buffer{}
	green := color.New(color.FgGreen).SprintfFunc()
	magenta := color.New(color.FgMagenta).SprintfFunc()

	_, err := buf.WriteString(green("\n%s", t.name))

	if err != nil {
		return "", err
	}

	if debug {
		buf.WriteString(magenta(" (# of nodes: %d)", t.Size()))
	}

	_, err = buf.Write([]byte{'\n', '.', '\n'})

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
	found := 0
	tnode := t.root
	var params map[string]string

	for tnode != nil && found < len(s) {
		next, p := tnode.next(s[found:], false, ph, delim)

		if p != nil {
			if params == nil {
				params = make(map[string]string)
			}

			count := 0

			for k, v := range p {
				found += len(v) - len(k) - 1

				params[k] = v
				count++
			}
		}

		if next != nil {
			tnode = next.node
			found += len(next.label)

			continue
		}

		tnode = nil
		params = nil
	}

	return tnode, params
}
