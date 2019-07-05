package bplustree

import (
	"testing"
)

func BenchmarkSearchDegree5(b *testing.B) {
	tree := New(5)
	for n := 0; n < maxElements; n++ {
		tree.Insert(data[n])
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tree.Search(data[n])
	}
}

func BenchmarkSearchDegree20(b *testing.B) {
	tree := New(20)
	for n := 0; n < maxElements; n++ {
		tree.Insert(data[n])
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		tree.Search(data[n])
	}
}

func BenchmarkSearchDegree30(b *testing.B) {
	tree := New(30)
	for n := 0; n < maxElements; n++ {
		tree.Insert(data[n])
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		tree.Search(data[n])
	}
}

func BenchmarkSearchDegree50(b *testing.B) {
	tree := New(50)
	for n := 0; n < maxElements; n++ {
		tree.Insert(data[n])
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tree.Search(data[n])
	}
}

func BenchmarkSearchDegree100(b *testing.B) {
	tree := New(100)
	for n := 0; n < maxElements; n++ {
		tree.Insert(data[n])
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tree.Search(data[n])
	}
}
