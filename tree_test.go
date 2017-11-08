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
		expectedValue  interface{}
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
			tree:         New("#1").WithNode("test", nil),
			str:          "test",
			expected:     true,
			expectedSize: 2,
		},
		// #2
		{
			tree:          New("#2").WithNode("test", "foo"),
			str:           "test",
			expected:      true,
			expectedSize:  2,
			expectedValue: "foo",
		},
		// #3
		{
			tree:          New("#3").WithNode("test", "foo").WithNode("testing", "bar"),
			str:           "test",
			expected:      true,
			expectedSize:  3,
			expectedValue: "foo",
		},
		// #4
		{
			tree:          New("#3").WithNode("test", "foo").WithNode("testing", "bar"),
			str:           "testing",
			expected:      true,
			expectedSize:  3,
			expectedValue: "bar",
		},
		// #5
		{
			tree:           New("#5").WithNode("test:@param", nil),
			str:            "test:foo",
			ph:             '@',
			expected:       true,
			expectedSize:   2,
			expectedParams: map[string]string{"param": "foo"},
		},
		// #6
		{
			tree:           New("#6").WithNode("test:@param", "foobar"),
			str:            "test:foo",
			ph:             '@',
			expected:       true,
			expectedSize:   2,
			expectedValue:  "foobar",
			expectedParams: map[string]string{"param": "foo"},
		},
		// #7
		{
			tree:           New("#7").WithNode("test:@param1:@param2", "foobar"),
			str:            "test:foo:bar",
			ph:             '@',
			delim:          ':',
			expected:       true,
			expectedSize:   2,
			expectedValue:  "foobar",
			expectedParams: map[string]string{"param1": "foo", "param2": "bar"},
		},
		// #8
		{
			tree:           New("#8").WithNode("test:@param1", "foo").WithNode("test:@param1:@param2", "bar"),
			str:            "test:foo:bar",
			ph:             '@',
			delim:          ':',
			expected:       true,
			expectedSize:   3,
			expectedValue:  "bar",
			expectedParams: map[string]string{"param1": "foo", "param2": "bar"},
		},
		// #9
		{
			tree: New("#9").
				WithNode("test", nil).
				WithNode("test:@param1", "foo").
				WithNode("test:@param1:@param2", "bar").
				WithNode("test:@param1:@param2:@param3", "baz"),
			str:            "test:foo:bar:baz",
			ph:             '@',
			delim:          ':',
			expected:       true,
			expectedSize:   5,
			expectedValue:  "baz",
			expectedParams: map[string]string{"param1": "foo", "param2": "bar", "param3": "baz"},
		},
		// #10
		{
			tree: New("#10").
				WithNode("test", nil).
				WithNode("test:@param1", "foo").
				WithNode("test:@param1:@param2", "bar").
				WithNode("test:@param1:@param2:@param3", "baz").
				WithoutNode("test:@param1"),
			str:            "test:foo:bar",
			ph:             '@',
			delim:          ':',
			expected:       true,
			expectedSize:   3,
			expectedValue:  "baz",
			expectedParams: map[string]string{"param2": "foo", "param3": "bar"},
		},
		// #11
		{
			tree: New("#11").
				WithNode("test", "foo").
				WithNode("test:@param1", "bar").
				WithNode("test:@param1:@param2", "baz").
				WithNode("test:@param1:@param2:@param3", "qux").
				WithoutNode("test:@param1"),
			str:            "test:foo:bar",
			ph:             '@',
			delim:          ':',
			expected:       true,
			expectedSize:   4,
			expectedValue:  "qux",
			expectedParams: map[string]string{"param2": "foo", "param3": "bar"},
		},
		// #12
		{
			tree:           New("#12").WithNode("/foo/:bar", "baz"),
			str:            "/foo/123",
			ph:             ':',
			delim:          '/',
			expected:       true,
			expectedSize:   2,
			expectedValue:  "baz",
			expectedParams: map[string]string{"bar": "123"},
		},
		// #13
		{
			tree:         New("#13").WithNode("/foo/:bar", "baz"),
			str:          "/foo/123/456",
			ph:           ':',
			delim:        '/',
			expected:     false,
			expectedSize: 2,
		},
		// #14
		{
			tree:           New("#14").WithNode("$foo|$bar", "baz"),
			str:            "abc|def",
			ph:             '$',
			delim:          '|',
			expected:       true,
			expectedSize:   2,
			expectedValue:  "baz",
			expectedParams: map[string]string{"foo": "abc", "bar": "def"},
		},
		// #15
		{
			tree:           New("#15").WithNode("$foo", "bar").WithNode("$foo|$baz", "qux"),
			str:            "abc|def",
			ph:             '$',
			delim:          '|',
			expected:       true,
			expectedSize:   3,
			expectedValue:  "qux",
			expectedParams: map[string]string{"foo": "abc", "baz": "def"},
		},
		// #16
		{
			tree:         New("#16").WithNode("/foo/:bar/baz", "qux"),
			str:          "/foo/123/qux",
			ph:           ':',
			delim:        '/',
			expected:     false,
			expectedSize: 2,
		},
		// #17
		{
			tree:           New("#17").WithNode("/foo/:bar/:baz", nil),
			str:            "/foo/123/456",
			ph:             ':',
			delim:          '/',
			expected:       true,
			expectedSize:   2,
			expectedParams: map[string]string{"bar": "123", "baz": "456"},
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
				a.Exactly(test.expectedValue, n.Value, index)
				t.Logf("n.Value = %#v\n", n.Value)
			}

			continue
		}

		n, params := test.tree.GetByRune(test.str, test.ph, test.delim)

		a.Exactly(test.expected, n != nil, index)
		a.Exactly(test.expectedParams, params, index)

		if n != nil {
			a.Exactly(test.expectedValue, n.Value, index)
			t.Logf("n.Value = %#v\n", n.Value)
		}
	}
}
