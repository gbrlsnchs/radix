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

func TestTree(t *testing.T) {
	testCases := []struct {
		labels      []string
		wrappers    []testWrapper
		length      int
		size        int
		params      map[string]string
		placeholder byte
		delim       byte
	}{
		{
			labels: []string{"foobar"},
			wrappers: []testWrapper{
				{label: "foobar", priority: 1, depth: 1, value: "bazqux"},
			},
			length: 2,
			size:   len("bazqux"),
		},
		{
			labels: []string{"a", "b"},
			wrappers: []testWrapper{
				{label: "a", priority: 1, depth: 1, value: 1},
				{label: "b", priority: 1, depth: 1, value: 2},
			},
			length: 3,
			size:   len("a") + len("b"),
		},
		{
			labels: []string{"a", "ab", "abc"},
			wrappers: []testWrapper{
				{label: "a", priority: 3, depth: 1, value: 1},
				{label: "ab", priority: 2, depth: 2, value: 2},
				{label: "abc", priority: 1, depth: 3, value: 3},
			},
			length: 4,
			size:   len("abc"),
		},
		{
			labels: []string{"a", "ab", "abc", "d"},
			wrappers: []testWrapper{
				{label: "a", priority: 3, depth: 1, value: 1},
				{label: "ab", priority: 2, depth: 2, value: 2},
				{label: "abc", priority: 1, depth: 3, value: 3},
				{label: "d", priority: 1, depth: 1, value: 4},
			},
			length: 5,
			size:   len("abc") + len("d"),
		},
		{
			labels: []string{"ab", "a", "abc"},
			wrappers: []testWrapper{
				{label: "ab", priority: 2, depth: 2, value: 2},
				{label: "a", priority: 3, depth: 1, value: 1},
				{label: "abc", priority: 1, depth: 3, value: 3},
			},
			length: 4,
			size:   len("abc"),
		},
		{
			labels: []string{"ab", "abc", "a"},
			wrappers: []testWrapper{
				{label: "ab", priority: 2, depth: 2, value: 2},
				{label: "abc", priority: 1, depth: 3, value: 3},
				{label: "a", priority: 3, depth: 1, value: 1},
			},
			length: 4,
			size:   len("abc"),
		},
		{
			labels: []string{"abc", "a", "ab"},
			wrappers: []testWrapper{
				{label: "abc", priority: 1, depth: 3, value: 3},
				{label: "a", priority: 3, depth: 1, value: 1},
				{label: "ab", priority: 2, depth: 2, value: 2},
			},
			length: 4,
			size:   len("abc"),
		},
		{
			labels: []string{"a", "b", "c"},
			wrappers: []testWrapper{
				{label: "a", priority: 1, depth: 1, value: 1},
				{label: "b", priority: 1, depth: 1, value: 2},
				{label: "c", priority: 1, depth: 1, value: 3},
			},
			length: 4,
			size:   len("a") + len("b") + len("c"),
		},
		{
			labels: []string{"/path/123"},
			wrappers: []testWrapper{
				{label: "/path/@id", priority: 1, depth: 1, value: "foobar"},
			},
			length:      2,
			size:        len("/path/@id"),
			params:      map[string]string{"id": "123"},
			placeholder: '@',
			delim:       '/',
		},
		{
			labels: []string{"/path/123/subpath/456"},
			wrappers: []testWrapper{
				{label: "/path/@id/subpath/@id2", priority: 1, depth: 1, value: "foobar"},
			},
			length:      2,
			size:        len("/path/@id/subpath/@id2"),
			params:      map[string]string{"id": "123", "id2": "456"},
			placeholder: '@',
			delim:       '/',
		},
		{
			labels: []string{"/path/123", "/path/123/subpath/456"},
			wrappers: []testWrapper{
				{label: "/path/@id", priority: 2, depth: 1, value: "foobar"},
				testWrapper{
					label:    "/path/@id/subpath/@id2",
					priority: 1,
					depth:    2,
					value:    "bazqux",
				},
			},
			length:      3,
			size:        len("/path/@id/subpath/@id2"),
			params:      map[string]string{"id": "123", "id2": "456"},
			placeholder: '@',
			delim:       '/',
		},
		{
			labels: []string{"/api/user/123"},
			wrappers: []testWrapper{
				{label: "/api/user/@id", priority: 1, depth: 1, value: "foobar"},
			},
			length:      2,
			size:        len("/api/user/@id"),
			params:      map[string]string{"id": "123"},
			placeholder: '@',
			delim:       '/',
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			tr := New(Tdebug)
			if tc.placeholder != 0 && tc.delim != 0 {
				tr.SetBoundaries(tc.placeholder, tc.delim)
			}
			for _, w := range tc.wrappers {
				tr.Add(w.label, w.value)
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
			for i, w := range tc.wrappers {
				n, p = tr.Get(tc.labels[i])
				if want, got := w.value, n.Value; !reflect.DeepEqual(want, got) {
					t.Errorf("want %v, got %v", want, got)
				}
				if want, got := w.priority, n.Priority(); want != got {
					t.Errorf("want %d, got %d", want, got)
				}
				if want, got := w.depth, n.Depth(); want != got {
					t.Errorf("want %d, got %d", want, got)
				}
			}
			if want, got := tc.params, p; !reflect.DeepEqual(want, got) {
				t.Errorf("want %v, got %v", want, got)
			}

			for i, w := range tc.wrappers {
				tr.Del(w.label)
				n, _ = tr.Get(tc.labels[i])
				if want, got := (*Node)(nil), n; want != got {
					t.Errorf("want %v, got %v", want, got)
				}
			}
		})
	}
}

func TestBinaryTree(t *testing.T) {
	testCases := []struct {
		labels []string
		values []interface{}
	}{
		{
			labels: []string{"foobar"},
			values: []interface{}{"bazqux"},
		},
		{
			labels: []string{"foo", "bar"},
			values: []interface{}{"baz", "qux"},
		},
		{
			labels: []string{"abc", "d"},
			values: []interface{}{"foo", "bar"},
		},
		{
			labels: []string{"a", "abc", "d"},
			values: []interface{}{"foo", "bar", "baz"},
		},
		{
			labels: []string{"foo", "bar", "baz", "qux"},
			values: []interface{}{1, 12, 123, 1234},
		},
		{
			labels: []string{
				"deck",
				"did",
				"doe",
				"dog",
				"doge",
				"dogs",
			},
			values: []interface{}{1, 2, 3, 4, 5, 6},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			tr := New(Tdebug | Tbinary)
			for i, v := range tc.values {
				tr.Add(tc.labels[i], v)
			}
			t.Log(tr.String())

			for i, v := range tc.values {
				label := tc.labels[i]
				n, _ := tr.Get(label)
				if want, got := v, n.Value; !reflect.DeepEqual(want, got) {
					t.Errorf("want %v, got %v", want, got)
				}
				if want, got := len(label)*8, n.Depth(); want != got {
					t.Errorf("want %d, got %d", want, got)
				}
			}
		})
	}
}
