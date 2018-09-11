package radix_test

import (
	"os"
	"testing"

	. "github.com/gbrlsnchs/radix"
)

var benchTree = New(Tdebug)

func TestMain(m *testing.M) {
	benchTree.Add("romane", 1)
	benchTree.Add("romanus", 2)
	benchTree.Add("romulus", 3)
	benchTree.Add("rubens", 4)
	benchTree.Add("ruber", 5)
	benchTree.Add("rubicon", 6)
	benchTree.Add("rubicundus", 7)
	os.Exit(m.Run())
}

func BenchmarkTree(b *testing.B) {
}
