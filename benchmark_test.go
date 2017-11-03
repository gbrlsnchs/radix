package patricia_test

import (
	"testing"

	. "github.com/gbrlsnchs/patricia"
)

func BenchmarkSingleStatic(b *testing.B) {
	b.ReportAllocs()

	tree := New("bench")

	tree.Add("/foo/bar/baz/qux", nil)

	for i := 0; i < b.N; i++ {
		_ = tree.Get("/foo/bar/baz/qux")
	}
}

func BenchmarkSingleDynamic(b *testing.B) {
	b.ReportAllocs()

	tree := New("bench")

	tree.Add("/foo/:bar/:baz/:qux", nil)

	for i := 0; i < b.N; i++ {
		_, _ = tree.GetByRune("/foo/123/456/789", ':', '/')
	}
}

func BenchmarkMultipleStatic(b *testing.B) {
	b.ReportAllocs()

	tree := New("bench")

	tree.Add("/foo", nil)
	tree.Add("/foo/bar", nil)
	tree.Add("/foo/bar/baz", nil)
	tree.Add("/foo/bar/baz/qux", nil)

	for i := 0; i < b.N; i++ {
		_ = tree.Get("/foo/bar/baz/qux")
	}
}

func BenchmarkMultipleDynamic(b *testing.B) {
	b.ReportAllocs()

	tree := New("bench")

	tree.Add("/foo", nil)
	tree.Add("/foo/:bar", nil)
	tree.Add("/foo/:bar/:baz", nil)
	tree.Add("/foo/:bar/:baz/:qux", nil)

	for i := 0; i < b.N; i++ {
		_, _ = tree.GetByRune("/foo/123/456/789", ':', '/')
	}
}
