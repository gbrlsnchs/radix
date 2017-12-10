package radix_test

import (
	"testing"

	. "github.com/gbrlsnchs/radix"
)

func BenchmarkSingleStatic(b *testing.B) {
	b.ReportAllocs()

	t := New("BenchmarkSingleStatic")

	t.Add("/foo/bar/baz/qux", nil)

	for i := 0; i < b.N; i++ {
		_ = t.Get("/foo/bar/baz/qux")
	}
}

func BenchmarkSingleDynamic(b *testing.B) {
	b.ReportAllocs()

	t := New("bench")

	t.Add("/foo/:bar/:baz/:qux", nil)

	for i := 0; i < b.N; i++ {
		_, _ = t.GetByRune("/foo/123/456/789", ':', '/')
	}
}

func BenchmarkMultipleStatic(b *testing.B) {
	b.ReportAllocs()

	t := New("BenchmarkMultipleStatic")

	t.Add("/foo", nil)
	t.Add("/foo/bar", nil)
	t.Add("/foo/bar/baz", nil)
	t.Add("/foo/bar/baz/qux", nil)

	for i := 0; i < b.N; i++ {
		_ = t.Get("/foo/bar/baz/qux")
	}
}

func BenchmarkMultipleDynamic(b *testing.B) {
	b.ReportAllocs()

	t := New("BenchmarkMultipleDynamic")

	t.Add("/foo", nil)
	t.Add("/foo/:bar", nil)
	t.Add("/foo/:bar/:baz", nil)
	t.Add("/foo/:bar/:baz/:qux", nil)

	for i := 0; i < b.N; i++ {
		_, _ = t.GetByRune("/foo/123/456/789", ':', '/')
	}
}

func BenchmarkLongString(b *testing.B) {
	b.ReportAllocs()

	t := New("BenchmarkLongString")

	t.Add("This is a very, very long string, so let's benchmark it.", nil)

	for i := 0; i < b.N; i++ {
		_ = t.Get("This is a very, very long string, so let's benchmark it.")
	}
}

func BenchmarkManyWords(b *testing.B) {
	b.ReportAllocs()

	t := New("BenchmarkManyWords")

	t.Add("romane", 1)
	t.Add("romanus", 2)
	t.Add("romulus", 3)
	t.Add("rubens", 4)
	t.Add("ruber", 5)
	t.Add("rubicon", 6)
	t.Add("rubicundus", 7)

	for i := 0; i < b.N; i++ {
		_ = t.Get("romanus")
	}
}
