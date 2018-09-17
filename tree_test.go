package radix_test

import (
	"reflect"
	"testing"

	. "github.com/gbrlsnchs/radix"
)

type testWrapper struct {
	label    string
	priority int
	depth    int
	value    interface{}
}

type testValues []*testWrapper

func TestPrefixTree(t *testing.T) {
	testCases := []struct {
		labels      []string
		values      testValues
		length      int
		size        int
		params      map[string]string
		placeholder byte
		delim       byte
	}{
		{
			labels: []string{"foobar"},
			values: testValues{&testWrapper{"foobar", 1, 1, "bazqux"}},
			length: 2,
			size:   len("bazqux"),
		},
		{
			labels: []string{"a", "ab", "abc"},
			values: testValues{&testWrapper{"a", 3, 1, 1}, &testWrapper{"ab", 2, 2, 2}, &testWrapper{"abc", 1, 3, 3}},
			length: 4,
			size:   len("abc"),
		},
		{
			labels: []string{"ab", "a", "abc"},
			values: testValues{&testWrapper{"ab", 2, 2, 2}, &testWrapper{"a", 3, 1, 1}, &testWrapper{"abc", 1, 3, 3}},
			length: 4,
			size:   len("abc"),
		},
		{
			labels: []string{"ab", "abc", "a"},
			values: testValues{&testWrapper{"ab", 2, 2, 2}, &testWrapper{"abc", 1, 3, 3}, &testWrapper{"a", 3, 1, 1}},
			length: 4,
			size:   len("abc"),
		},
		{
			labels: []string{"abc", "a", "ab"},
			values: testValues{&testWrapper{"abc", 1, 3, 3}, &testWrapper{"a", 3, 1, 1}, &testWrapper{"ab", 2, 2, 2}},
			length: 4,
			size:   len("abc"),
		},
		{
			labels: []string{"a", "b", "c"},
			values: testValues{&testWrapper{"a", 1, 1, 1}, &testWrapper{"b", 1, 1, 2}, &testWrapper{"c", 1, 1, 3}},
			length: 4,
			size:   len("a") + len("b") + len("c"),
		},
		{
			labels:      []string{"/path/123"},
			values:      testValues{&testWrapper{"/path/@id", 1, 1, "foobar"}},
			length:      2,
			size:        len("/path/@id"),
			params:      map[string]string{"id": "123"},
			placeholder: '@',
			delim:       '/',
		},
		{
			labels:      []string{"/path/123/subpath/456"},
			values:      testValues{&testWrapper{"/path/@id/subpath/@id2", 1, 1, "foobar"}},
			length:      2,
			size:        len("/path/@id/subpath/@id2"),
			params:      map[string]string{"id": "123", "id2": "456"},
			placeholder: '@',
			delim:       '/',
		},
		{
			labels:      []string{"/path/123", "/path/123/subpath/456"},
			values:      testValues{&testWrapper{"/path/@id", 2, 1, "foobar"}, &testWrapper{"/path/@id/subpath/@id2", 1, 2, "foobar"}},
			length:      3,
			size:        len("/path/@id/subpath/@id2"),
			params:      map[string]string{"id": "123", "id2": "456"},
			placeholder: '@',
			delim:       '/',
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			tr := New(Tdebug)
			if tc.placeholder > 0 && tc.delim > 0 {
				tr.SetBoundaries(tc.placeholder, tc.delim)
			}
			for _, v := range tc.values {
				tr.Add(v.label, v.value)
			}
			t.Log(tr.String())

			if want, got := tc.length, tr.Len(); want != got {
				t.Errorf("want %d, got %d", want, got)
			}
			if want, got := tc.size, tr.Size(); want != got {
				t.Errorf("want %d, got %d", want, got)
			}
			var (
				n *Node
				p map[string]string
			)
			for i, v := range tc.values {
				n, p = tr.Get(tc.labels[i])
				if want, got := true, n != nil; want != got {
					t.Fatalf("want %t, got %t", want, got)
				}
				if want, got := v.value, n.Value; !reflect.DeepEqual(want, got) {
					t.Errorf("want %v, got %v", want, got)
				}
				if want, got := v.priority, n.Priority(); want != got {
					t.Errorf("want %d, got %d", want, got)
				}
				if want, got := v.depth, n.Depth(); want != got {
					t.Errorf("want %d, got %d", want, got)
				}
			}
			if want, got := tc.params, p; !reflect.DeepEqual(want, got) {
				t.Errorf("want %v, got %v", want, got)
			}

			for i, v := range tc.values {
				tr.Del(v.label)
				n, _ = tr.Get(tc.labels[i])
				if want, got := (*Node)(nil), n; want != got {
					t.Errorf("want %v, got %v", want, got)
				}
			}
		})
	}
}
