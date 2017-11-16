package patricia_test

import (
	"fmt"

	"github.com/gbrlsnchs/patricia"
)

func Example() {
	// Example from https://upload.wikimedia.org/wikipedia/commons/a/ae/Patricia_trie.svg.
	t := patricia.New("Example").
		WithNode("romane", 1).
		WithNode("romanus", 2).
		WithNode("romulus", 3).
		WithNode("rubens", 4).
		WithNode("ruber", 5).
		WithNode("rubicon", 6).
		WithNode("rubicundus", 7)

	t.Sort(patricia.AscLabelSort)

	err := t.Debug()

	if err != nil {
		// ...
	}

	t.Sort(patricia.DescLabelSort)

	err = t.Debug()

	if err != nil {
		// ...
	}

	t.Sort(patricia.PrioritySort)

	err = t.Debug()

	if err != nil {
		// ...
	}

	n := t.Get("romanus")

	fmt.Println(n.Value)
	// Output: 2
}

func Example_named() {
	t := patricia.New("Named Edge Example").
		WithNode("foo@bar!@baz", nil)

	err := t.Debug()

	if err != nil {
		// ...
	}

	_, params := t.GetByRune("foo123!456", '@', '!')

	fmt.Println(params["bar"])
	fmt.Println(params["baz"])
	// Output:
	// 123
	// 456
}
