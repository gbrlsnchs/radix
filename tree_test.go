package radix_test

import (
	"strconv"
	"testing"
	"time"

	. "github.com/gbrlsnchs/radix"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	a := assert.New(t)
	tests := []*struct {
		tree           *Tree
		handlerFunc    func(*Tree)
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
			tree:          New("#1"),
			str:           "test",
			expected:      true,
			expectedSize:  2,
			expectedValue: "foo",
			handlerFunc: func(t *Tree) {
				t.Add("test", "foo")
			},
		},
		// #2
		{
			tree:          New("#2"),
			str:           "test",
			expected:      true,
			expectedSize:  2,
			expectedValue: "foo",
			handlerFunc: func(t *Tree) {
				t.Add("test", "foo")
			},
		},
		// #3
		{
			tree:          New("#3"),
			str:           "test",
			expected:      true,
			expectedSize:  3,
			expectedValue: "foo",
			handlerFunc: func(t *Tree) {
				t.Add("test", "foo")
				t.Add("testing", "bar")
			},
		},
		// #4
		{
			tree:          New("#4"),
			str:           "testing",
			expected:      true,
			expectedSize:  3,
			expectedValue: "bar",
			handlerFunc: func(t *Tree) {
				t.Add("test", "foo")
				t.Add("testing", "bar")
			},
		},
		// #5
		{
			tree:           New("#5"),
			str:            "test:foo",
			ph:             '@',
			expected:       true,
			expectedSize:   2,
			expectedValue:  "foo",
			expectedParams: map[string]string{"param": "foo"},
			handlerFunc: func(t *Tree) {
				t.Add("test:@param", "foo")
			},
		},
		// #6
		{
			tree:           New("#6"),
			str:            "test:foo",
			ph:             '@',
			expected:       true,
			expectedSize:   2,
			expectedValue:  "foobar",
			expectedParams: map[string]string{"param": "foo"},
			handlerFunc: func(t *Tree) {
				t.Add("test:@param", "foobar")
			},
		},
		// #7
		{
			tree:           New("#7"),
			str:            "test:foo:bar",
			ph:             '@',
			delim:          ':',
			expected:       true,
			expectedSize:   2,
			expectedValue:  "foobar",
			expectedParams: map[string]string{"param1": "foo", "param2": "bar"},
			handlerFunc: func(t *Tree) {
				t.Add("test:@param1:@param2", "foobar")
			},
		},
		// #8
		{
			tree:           New("#8"),
			str:            "test:foo:bar",
			ph:             '@',
			delim:          ':',
			expected:       true,
			expectedSize:   3,
			expectedValue:  "bar",
			expectedParams: map[string]string{"param1": "foo", "param2": "bar"},
			handlerFunc: func(t *Tree) {
				t.Add("test:@param1", "foo")
				t.Add("test:@param1:@param2", "bar")
			},
		},
		// #9
		{
			tree:           New("#9"),
			str:            "test:foo:bar:baz",
			ph:             '@',
			delim:          ':',
			expected:       true,
			expectedSize:   5,
			expectedValue:  "baz",
			expectedParams: map[string]string{"param1": "foo", "param2": "bar", "param3": "baz"},
			handlerFunc: func(t *Tree) {
				t.Add("test", ".")
				t.Add("test:@param1", "foo")
				t.Add("test:@param1:@param2", "bar")
				t.Add("test:@param1:@param2:@param3", "baz")
			},
		},
		// #10
		{
			tree:           New("#10"),
			str:            "test:foo:bar",
			ph:             '@',
			delim:          ':',
			expected:       true,
			expectedSize:   4,
			expectedValue:  "baz",
			expectedParams: map[string]string{"param2": "foo", "param3": "bar"},
			handlerFunc: func(t *Tree) {
				t.Add("test", "foo")
				t.Add("test:@param1", "foo")
				t.Add("test:@param1:@param2", "bar")
				t.Add("test:@param1:@param2:@param3", "baz")
				t.Del("test:@param1")
			},
		},
		// #11
		{
			tree:           New("#11"),
			str:            "test:foo:bar",
			ph:             '@',
			delim:          ':',
			expected:       true,
			expectedSize:   4,
			expectedValue:  "qux",
			expectedParams: map[string]string{"param2": "foo", "param3": "bar"},
			handlerFunc: func(t *Tree) {
				t.Add("test", "foo")
				t.Add("test:@param1", "bar")
				t.Add("test:@param1:@param2", "baz")
				t.Add("test:@param1:@param2:@param3", "qux")
				t.Del("test:@param1")
			},
		},
		// #12
		{
			tree:           New("#12"),
			str:            "/foo/123",
			ph:             ':',
			delim:          '/',
			expected:       true,
			expectedSize:   2,
			expectedValue:  "baz",
			expectedParams: map[string]string{"bar": "123"},
			handlerFunc: func(t *Tree) {
				t.Add("/foo/:bar", "baz")
			},
		},
		// #13
		{
			tree:         New("#13"),
			str:          "/foo/123/456",
			ph:           ':',
			delim:        '/',
			expected:     false,
			expectedSize: 2,
			handlerFunc: func(t *Tree) {
				t.Add("/foo/:bar", "baz")
			},
		},
		// #14
		{
			tree:           New("#14"),
			str:            "abc|def",
			ph:             '$',
			delim:          '|',
			expected:       true,
			expectedSize:   2,
			expectedValue:  "baz",
			expectedParams: map[string]string{"foo": "abc", "bar": "def"},
			handlerFunc: func(t *Tree) {
				t.Add("$foo|$bar", "baz")
			},
		},
		// #15
		{
			tree:           New("#15"),
			str:            "abc|def",
			ph:             '$',
			delim:          '|',
			expected:       true,
			expectedSize:   3,
			expectedValue:  "qux",
			expectedParams: map[string]string{"foo": "abc", "baz": "def"},
			handlerFunc: func(t *Tree) {
				t.Add("$foo", "bar")
				t.Add("$foo|$baz", "qux")
			},
		},
		// #16
		{
			tree:         New("#16"),
			str:          "/foo/123/qux",
			ph:           ':',
			delim:        '/',
			expected:     false,
			expectedSize: 2,
			handlerFunc: func(t *Tree) {
				t.Add("/foo/:bar/baz", "qux")
			},
		},
		// #17
		{
			tree:           New("#17"),
			str:            "/foo/123/456",
			ph:             ':',
			delim:          '/',
			expected:       true,
			expectedSize:   2,
			expectedValue:  "foo",
			expectedParams: map[string]string{"bar": "123", "baz": "456"},
			handlerFunc: func(t *Tree) {
				t.Add("/foo/:bar/:baz", "foo")
			},
		},
		// #18
		{
			tree:          New("#18"),
			str:           "testing",
			expected:      true,
			expectedSize:  3,
			expectedValue: "foo",
			handlerFunc: func(t *Tree) {
				t.Add("testing", "foo")
				t.Add("test", "bar")
			},
		},
		// #19
		{
			tree:           New("#19"),
			str:            "foo123",
			ph:             '*',
			expected:       true,
			expectedSize:   3,
			expectedValue:  "foo",
			expectedParams: map[string]string{"bar": "123"},
			handlerFunc: func(t *Tree) {
				t.Add("foo*bar", "foo")
				t.Add("foo", "bar")
			},
		},
	}

	for i, test := range tests {
		index := strconv.Itoa(i)

		if test.handlerFunc != nil {
			test.handlerFunc(test.tree)
		}

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

func TestRace(t *testing.T) {
	list := []string{
		"foo",
		"bar",
		"foobar",
		"foobarbaz",
		"qux",
		"barbazqux",
	}
	tree := New("TestRace")
	tree.Safe = true

	for i, n := range list {
		go func(i int, n string) {
			tree.Add(n, i)
			time.Sleep(time.Second * 3)
		}(i, n)
	}

	for _, n := range list {
		go func(n string) {
			_ = tree.Get(n)

			time.Sleep(time.Second * 3)
		}(n)
	}
}
