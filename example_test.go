package patricia_test

import (
	"fmt"

	"github.com/gbrlsnchs/patricia"
)

func Example() {
	// Example from https://upload.wikimedia.org/wikipedia/commons/a/ae/Patricia_trie.svg.
	t := patricia.New("Example").
		Add("romane", 1).
		Add("romanus", 2).
		Add("romulus", 3).
		Add("rubens", 4).
		Add("ruber", 5).
		Add("rubicon", 6).
		Add("rubicundus", 7)

	err := t.Debug()

	if err != nil {
		// ...
	}

	n := t.Get("romanus")

	fmt.Println(n.Value)
	// Output: 2
}

func Example_named() {
	t := patricia.New("Named Edge Example").
		Add("foo@bar!@baz", nil)

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
