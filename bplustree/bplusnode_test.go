package bplustree

import (
	"log"
	"math/rand"
	"sort"
	"testing"
)

func TestPopElement(t *testing.T) {
	node := newNode(5, true)

	//rand.Seed(time.Now().Unix())
	rand.Seed(1)
	var i int

	vals := []int{10, 20, 30, 40, 50, 60, 5, 1, 70}

	for i < len(vals) {
		node.insert(vals[i])
		i++
	}

	val := node.popElement()
	if val != 70 {
		log.Fatalln("Should pop 85 got ", val)
	}

	if node.elementsCount != len(vals)-1 {
		log.Fatalln("Elements count wrong")
	}
}

func TestPrependElement(t *testing.T) {
	node := newNode(5, true)

	//rand.Seed(time.Now().Unix())
	rand.Seed(1)

	vals := []int{10, 20, 30, 40, 50, 60, 5, 70}

	for i := 0; i < len(vals); i++ {
		node.insert(vals[i])
	}

	node.prependElement(1)
	if node.values[0] != 1 {
		log.Fatalln("Should prepend element got ", node.values[0])
	}

	sort.Ints(vals)
	for i := 0; i < len(vals); i++ {
		if node.values[i+1] != vals[i] {
			log.Fatalln("Should retain order. Not equal for", i)
		}
	}
}

func TestPopChildren(t *testing.T) {
	node := newNode(5, true)

	//rand.Seed(time.Now().Unix())
	rand.Seed(1)

	vals := []*treeNode{
		&treeNode{},
		&treeNode{},
		&treeNode{degree: 5},
	}

	for i := 0; i < len(vals); i++ {
		node.children[i] = vals[i]
	}

	node.elementsCount = 2

	c := node.popChildren()

	if c == nil {
		log.Fatalln("Should pop existing child")
	}

	if node.children[node.elementsCount] != nil {
		log.Fatalln("Should delete popped children")
	}

	if c.degree != vals[2].degree {
		log.Fatalln("Wrong children popped")
	}
}

func TestPrependChildren(t *testing.T) {
	node := newNode(5, true)

	//rand.Seed(time.Now().Unix())
	rand.Seed(1)

	vals := []*treeNode{
		&treeNode{},
		&treeNode{},
		&treeNode{},
	}

	for i := 0; i < len(vals); i++ {
		node.children[i] = vals[i]
	}

	node.elementsCount = 2

	node.prependChildren(&treeNode{degree: 10})

	if node.children[0].degree != 10 {
		log.Fatalln("Children appended incorrectly")
	}
}

func TestBPlusNode_deleteKey(t *testing.T) {
	type args struct {
		elements []int
		val      int
	}
	tests := []struct {
		name string
		n    *treeNode
		args args
		want bool
	}{
		{
			"Simple test",
			newNode(5, true),
			args{[]int{10, 40, 60, 70}, 10},
			true,
		},

		{
			"No key",
			newNode(5, true),
			args{[]int{10, 40, 60, 70}, 111},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			for _, v := range tt.args.elements {
				tt.n.insert(v)
			}

			if got := tt.n.delete(tt.args.val); got != tt.want {
				t.Errorf("treeNode.delete() = %v, want %v", got, tt.want)
			}
		})
	}
}
