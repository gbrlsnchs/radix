package radix_test

import (
	"fmt"

	"github.com/gbrlsnchs/radix"
)

func Example() {
	// Example from https://upload.wikimedia.org/wikipedia/commons/a/ae/Patricia_trie.svg.
	t := radix.New("Example")

	t.Add("romane", 1)
	t.Add("romanus", 2)
	t.Add("romulus", 3)
	t.Add("rubens", 4)
	t.Add("ruber", 5)
	t.Add("rubicon", 6)
	t.Add("rubicundus", 7)
	t.Sort(radix.AscLabelSort)

	err := t.Debug()

	if err != nil {
		// ...
	}

	t.Sort(radix.DescLabelSort)

	err = t.Debug()

	if err != nil {
		// ...
	}

	t.Sort(radix.PrioritySort)

	err = t.Debug()

	if err != nil {
		// ...
	}

	n := t.Get("romanus")

	fmt.Println(n.Value)
	// Output: 2
}

func Example_named() {
	t := radix.New("Named Edge Example")

	t.Add("foo@bar!@baz", nil)

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
