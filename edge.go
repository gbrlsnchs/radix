package patricia

import (
	"bytes"

	"github.com/fatih/color"
)

// edge is a PATRICIA tree edge.
type edge struct {
	label string
	node  *Node
}

// newEdge creates a new edge.
func newEdge(label string, node *Node) *edge {
	return &edge{
		label: label,
		node:  node,
	}
}

// buffer returns a pointer to a bytes.Buffer containing
// a subtree structure plus, if debug is truthy, its metadata.
func (e *edge) buffer(debug bool, tabList []bool) (*bytes.Buffer, error) {
	branches := []byte{}
	buf := &bytes.Buffer{}
	red := color.New(color.FgRed).SprintfFunc()
	green := color.New(color.FgGreen).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintfFunc()

	if len(tabList) > 1 {
		for _, tlist := range tabList[:len(tabList)-1] {
			if tlist {
				branches = append(branches, "    "...)

				continue
			}

			branches = append(branches, "|   "...)
		}
	}

	branches = append(branches, "+-- "...)

	_, err := buf.Write(branches)

	if err != nil {
		return nil, err
	}

	if debug {
		_, err = buf.WriteString(red("(%d) ", e.node.priority))

		if err != nil {
			return nil, err
		}
	}

	_, err = buf.WriteString(e.label)

	if err != nil {
		return nil, err
	}

	if e.node.IsLeaf() {
		_, err = buf.WriteString(green(" 🍂"))

		if err != nil {
			return nil, err
		}
	}

	if debug {
		_, err = buf.WriteString(magenta(" -> %#v", e.node.Value))

		if err != nil {
			return nil, err
		}
	}

	err = buf.WriteByte('\n')

	if err != nil {
		return nil, err
	}

	for i, next := range e.node.edges {
		var nbuf *bytes.Buffer

		if len(tabList) < next.node.depth {
			tabList = append(tabList, i == len(e.node.edges)-1)
		} else {
			tabList[next.node.depth-1] = i == len(e.node.edges)-1
		}

		nbuf, err = next.buffer(debug, tabList)

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
