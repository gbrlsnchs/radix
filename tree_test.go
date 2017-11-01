package patricia_test

import (
	"strconv"
	"testing"

	. "github.com/gbrlsnchs/patricia"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	a := assert.New(t)
	tests := []*struct {
		tree           *Tree
		str            string
		ph             rune
		delim          rune
		expected       bool
		expectedSize   uint
		expectedVal    interface{}
		expectedParams map[string]string
	}{
		// #0
		{
			tree:         New("#0"),
			str:          "test",
			expectedSize: 1,
		},
		// #1
		{
			tree:         New("#1").Add("test", nil),
			str:          "test",
			expected:     true,
			expectedSize: 2,
		},
		// #2
		{
			tree:         New("#2").Add("test", "foo"),
			str:          "test",
			expected:     true,
			expectedSize: 2,
			expectedVal:  "foo",
		},
		// #3
		{
			tree:         New("#3").Add("test", "foo").Add("testing", "bar"),
			str:          "test",
			expected:     true,
			expectedSize: 3,
			expectedVal:  "foo",
		},
		// #4
		{
			tree:         New("#3").Add("test", "foo").Add("testing", "bar"),
			str:          "testing",
			expected:     true,
			expectedSize: 3,
			expectedVal:  "bar",
		},
		// #5
		{
			tree:           New("#5").Add("test:@param", nil),
			str:            "test:foo",
			ph:             '@',
			expected:       true,
			expectedSize:   2,
			expectedParams: map[string]string{"param": "foo"},
		},
		// #6
		{
			tree:           New("#6").Add("test:@param", "foobar"),
			str:            "test:foo",
			ph:             '@',
			expected:       true,
			expectedSize:   2,
			expectedVal:    "foobar",
			expectedParams: map[string]string{"param": "foo"},
		},
		// #7
		{
			tree:           New("#7").Add("test:@param1:@param2", "foobar"),
			str:            "test:foo:bar",
			ph:             '@',
			delim:          ':',
			expected:       true,
			expectedSize:   2,
			expectedVal:    "foobar",
			expectedParams: map[string]string{"param1": "foo", "param2": "bar"},
		},
		// #8
		{
			tree:           New("#8").Add("test:@param1", "foo").Add("test:@param1:@param2", "bar"),
			str:            "test:foo:bar",
			ph:             '@',
			delim:          ':',
			expected:       true,
			expectedSize:   3,
			expectedVal:    "bar",
			expectedParams: map[string]string{"param1": "foo", "param2": "bar"},
		},
		// #9
		{
			tree: New("#9").
				Add("test", nil).
				Add("test:@param1", "foo").
				Add("test:@param1:@param2", "bar").
				Add("test:@param1:@param2:@param3", "baz"),
			str:            "test:foo:bar:baz",
			ph:             '@',
			delim:          ':',
			expected:       true,
			expectedSize:   5,
			expectedVal:    "baz",
			expectedParams: map[string]string{"param1": "foo", "param2": "bar", "param3": "baz"},
		},
		// #10
		{
			tree: New("#10").
				Add("test", nil).
				Add("test:@param1", "foo").
				Add("test:@param1:@param2", "bar").
				Add("test:@param1:@param2:@param3", "baz").
				Del("test:@param1"),
			str:            "test:foo:bar",
			ph:             '@',
			delim:          ':',
			expected:       true,
			expectedSize:   3,
			expectedVal:    "baz",
			expectedParams: map[string]string{"param2": "foo", "param3": "bar"},
		},
		// #11
		{
			tree: New("#11").
				Add("test", "foo").
				Add("test:@param1", "bar").
				Add("test:@param1:@param2", "baz").
				Add("test:@param1:@param2:@param3", "qux").
				Del("test:@param1"),
			str:            "test:foo:bar",
			ph:             '@',
			delim:          ':',
			expected:       true,
			expectedSize:   4,
			expectedVal:    "qux",
			expectedParams: map[string]string{"param2": "foo", "param3": "bar"},
		},
		// #12
		{
			tree:           New("#12").Add("/foo/:bar", "baz"),
			str:            "/foo/123",
			ph:             ':',
			delim:          '/',
			expected:       true,
			expectedSize:   2,
			expectedVal:    "baz",
			expectedParams: map[string]string{"bar": "123"},
		},
	}

	for i, test := range tests {
		index := strconv.Itoa(i)
		err := test.tree.Debug()

		a.Nil(err, index)
		a.Exactly(test.expectedSize, test.tree.Size(), index)

		if test.ph == 0 {
			n := test.tree.Get(test.str)

			a.Exactly(test.expected, n != nil, index)

			if n != nil {
				a.Exactly(test.expectedVal, n.Val, index)
				t.Logf("n.Val = %#v\n", n.Val)
			}

			continue
		}

		n, params := test.tree.GetByRune(test.str, test.ph, test.delim)

		a.Exactly(test.expected, n != nil, index)
		a.Exactly(test.expectedParams, params, index)

		if n != nil {
			a.Exactly(test.expectedVal, n.Val, index)
			t.Logf("n.Val = %#v\n", n.Val)
		}
	}
}
