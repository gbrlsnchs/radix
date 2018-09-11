package radix_test

import (
	"fmt"

	"github.com/gbrlsnchs/radix"
)

func ExampleTree() {
	tr := radix.New(radix.Tdebug)
	tr.Add("romane", 1)
	tr.Add("romanus", 2)
	tr.Add("romulus", 3)
	tr.Add("rubens", 4)
	tr.Add("ruber", 5)
	tr.Add("rubicon", 6)
	tr.Add("rubicundus", 7)
	fmt.Printf("%v\n", tr)
}
