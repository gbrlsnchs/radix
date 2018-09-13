package radix_test

import (
	"reflect"
	"testing"

	. "github.com/gbrlsnchs/radix"
)

func TestTree(t *testing.T) {
	testCases := []struct {
		flags  int
		values map[string]interface{}
		query  string
		v      interface{}
		length int
		size   int
	}{
		{Tdebug, map[string]interface{}{"foobar": "raboof"}, "foobar", "raboof", 2, len("raboof")},
		{Tdebug | Tbinary, map[string]interface{}{"foobar": "raboof"}, "foobar", "raboof", 1 + len("raboof")*8, len("raboof")},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			tr := New(tc.flags)
			for k, v := range tc.values {
				tr.Add(k, v)
			}
			var n *Node
			n, _ = tr.Get(tc.query)
			if want, got := tc.v, n.Value; !reflect.DeepEqual(want, got) {
				t.Errorf("want %v, got %v", want, got)
			}
			if want, got := tc.length, tr.Len(); want != got {
				t.Errorf("want %d, got %d", want, got)
			}
			if want, got := tc.size, tr.Size(); want != got {
				t.Errorf("want %d, got %d", want, got)
			}
			t.Log(tr.String())

			tr.Del(tc.query)
			n, _ = tr.Get(tc.query)
			if want, got := (*Node)(nil), n; want != got {
				t.Errorf("want %v, got %v", want, got)
			}
			t.Log(tr.String())
		})
	}
}
