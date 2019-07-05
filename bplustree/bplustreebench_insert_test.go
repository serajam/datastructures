package bplustree

import (
	"math/rand"
	"os"
	"testing"
	"time"
)

var data []int
var maxElements = 10000000

func TestMain(m *testing.M) {
	rand.Seed(time.Now().Unix())
	max := 10000000
	for n := 0; n < maxElements; n++ {
		data = append(data, rand.Intn(1+max))
	}

	retCode := m.Run()
	os.Exit(retCode)
}

func BenchmarkInsertDegree5(b *testing.B) {
	tree := New(5)
	for n := 0; n < b.N; n++ {
		tree.Insert(data[n])
	}
}

func BenchmarkInsertDegree20(b *testing.B) {
	tree := New(5)
	for n := 0; n < b.N; n++ {
		tree.Insert(data[n])
	}
}

func BenchmarkInsertDegree30(b *testing.B) {
	tree := New(30)
	for n := 0; n < b.N; n++ {
		tree.Insert(data[n])
	}
}

func BenchmarkInsertDegree50(b *testing.B) {
	tree := New(50)
	for n := 0; n < b.N; n++ {
		tree.Insert(data[n])
	}
}

func BenchmarkInsertDegree100(b *testing.B) {
	tree := New(100)
	for n := 0; n < b.N; n++ {
		tree.Insert(data[n])
	}
}
