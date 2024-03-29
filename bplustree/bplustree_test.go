package bplustree

import (
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestInsertLinkedList(t *testing.T) {
	tree := New(5)
	r := tree.Root
	var i int

	add := 50

	for i = add; i > 0; i-- {
		tree.Insert(i * 10)
	}

	if tree.count != add {
		log.Fatal("Expected", add, "elements", "but got", tree.count)
	}

	for i = add; i > 0; i-- {
		if tree.Search(i*10) != true {
			log.Fatal("Not found:", i*10)
		}
	}

	tree.Delete(310)
	tree.Delete(320)
	tree.Delete(330)
	tree.Delete(340)
	tree.Delete(350)

	leaf := r
	total := 0
	for leaf != nil {
		total += leaf.elementsCount
		leaf = leaf.nextLeaf
	}

	if total != tree.count {
		log.Fatalln("Elements count and real leaf values doesnt much")
	}

}

func TestBPlusDelete(t *testing.T) {
	tree := New(5)
	var i int

	vals := []int{10, 20, 30, 40, 50, 60, 5, 1, 70, 70, 80, 85, 25, 35}

	for i < len(vals) {
		tree.Insert(vals[i])
		i++
	}

	if tree.Delete(10) == false {
		log.Fatalln("Failed to delete", 10)
	}

	if tree.Search(10) == true {
		log.Fatalln("Element", 10, "was not deleted")
	}
}

func TestBPlusInsertLinear(t *testing.T) {
	tree := New(5)

	rand.Seed(time.Now().Unix())
	i := 200000

	for i > 0 {
		tree.Insert(i + rand.Intn(i))
		i--
	}
}

func TestInsertAndSearch(t *testing.T) {
	tree := New(5)
	var i int

	add := 2000000

	for i = add; i > 0; i-- {
		tree.Insert(i * 10)
	}

	if tree.count != add {
		log.Fatal("Expected", add, "elements", "but got", tree.count)
	}

	for i = add; i > 0; i-- {
		if tree.Search(i*10) != true {
			log.Fatal("Not found:", i*10)
		}
	}
}

func TestInsertAndDelete(t *testing.T) {
	tree := New(5)
	var i int

	add := 2000000

	for i = add; i > 0; i-- {
		tree.Insert(i * 10)
	}

	if tree.count != add {
		log.Fatal("Expected", add, "elements", "but got", tree.count)
	}

	del := 1999999
	for i = del; i > 0; i-- {
		tree.Delete(i * 10)
	}

	if tree.count != add-del {
		log.Fatal("Expected", add-del, "elements", "but got", tree.count)
	}
}

func TestInsertAndDelete2(t *testing.T) {
	tree := New(5)
	var i int

	add := 200000

	for i = add; i > 0; i-- {
		tree.Insert(i * 10)
	}

	if tree.count != add {
		log.Fatal("Expected", add, "elements", "but got", tree.count)
	}

	del := 199999

	for i = del; i > 0; i-- {
		tree.Delete(i * 10)
	}

	if tree.count != add-del {
		log.Fatal("Expected", add-del, "elements", "but got", tree.count)
	}

	add = 200000

	for i = add; i > 0; i-- {
		tree.Insert(i)
	}

	if tree.count != add+1 {
		log.Fatal("Expected", add, "elements", "but got", tree.count)
	}

	del = 199999

	for i = del; i > 0; i-- {
		tree.Delete(i)
	}

	if tree.count != add-del+1 {
		log.Fatal("Expected", add-del+1, "elements", "but got", tree.count)
	}
}

func TestInsertAndDelete3(t *testing.T) {
	tree := New(5)
	r := tree.Root
	var i int

	add := 200000

	for i = add; i > 0; i-- {
		tree.Insert(i * 10)
	}

	if tree.count != add {
		log.Fatal("Expected", add, "elements", "but got", tree.count)
	}

	del := 199999

	for i = del; i > 0; i-- {
		tree.Delete(i * 10)
	}

	if tree.count != add-del {
		log.Fatal("Expected", add-del, "elements", "but got", tree.count)
	}

	add = 200000

	for i = add; i > 0; i-- {
		tree.Insert(i)
	}

	if tree.count != add+1 {
		log.Fatal("Expected", add, "elements", "but got", tree.count)
	}

	del = 199999

	for i = del; i > 0; i-- {
		tree.Delete(i)
	}

	if tree.count != add-del+1 {
		log.Fatal("Expected", add-del+1, "elements", "but got", tree.count)
	}

	leaf := r
	total := 0
	for leaf != nil {
		total += leaf.elementsCount
		leaf = leaf.nextLeaf
	}

	delRandom := 10000
	rand.Seed(time.Now().Unix())

	for i = delRandom; i > 0; i-- {
		tree.Delete(rand.Intn(add) + 1)
	}

	if total != tree.count {
		log.Fatalln("Elements count and real leaf values doesn't much")
	}
}
