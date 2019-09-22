// Package bplustree provides implementation of b+ tree algorithm
// it offers search, insert, delete methods
package bplustree

import "fmt"

// Tree represents root structure for holding nodes and secondary data
type Tree struct {
	Root   *treeNode
	Degree int

	count int
}

// New creates new root node
func New(degree int) *Tree {
	return &Tree{Degree: degree, Root: newNode(degree, true)}
}

// Traverse display tree data
func (t *Tree) Traverse() {
	// TODO make it more human friendly

	if t.Root == nil {
		return
	}

	fmt.Println("-------------")
	fmt.Println(t.Root.traverse(0))
	fmt.Println("-------------")
}

// Insert inserts new value and splits every overflow node it meet
func (t *Tree) Insert(val int) {
	var parent, node *treeNode
	var nextNodeIndex int
	node = t.Root

	for node != nil && !node.leaf {
		if node.hasEmptyCell() {
			parent = node
			nextNodeIndex = node.searchNext(val)
			node = node.sibling(nextNodeIndex)

			continue
		}

		if parent == nil {
			node = node.splitNode(nil, t.Degree, -1)
			t.Root = node
		} else {
			node = node.splitNode(parent, t.Degree, nextNodeIndex)
			nextNodeIndex = node.searchNext(val)
			node = node.sibling(nextNodeIndex)
		}
		continue
	}

	if node == nil {
		return
	}

	node.insert(val)
	t.count++

	// guarantee we do not have full leafs
	if node.hasEmptyCell() {
		return
	}

	if parent == nil {
		node = node.splitNode(nil, t.Degree, -1)
		t.Root = node
	} else {
		node.splitNode(parent, t.Degree, nextNodeIndex)
	}
}

// Search searches for value
func (t *Tree) Search(val int) bool {
	node := t.Root
	for !node.leaf {
		var i int
		for i < node.elementsCount {
			if node.values[i] > val {
				break
			}

			if node.values[i] == 0 {
				break
			}

			i++
		}

		node = node.siblings[i]

		if node == nil {
			return false
		}
	}

	var i int
	for i < node.elementsCount {
		if val == node.values[i] {
			return true
		}
		i++
	}

	return false
}

// Delete deletes values and checks for underflow
func (t *Tree) Delete(val int) bool {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(val)
			panic(err)
		}
	}()

	result := t.Root.delete(val)
	if !result {
		return false
	}

	t.count--

	if t.Root.elementsCount == 0 && !t.Root.leaf {
		t.Root = t.Root.siblings[0]
		return true
	}

	if i := t.Root.keyIndex(val); i >= 0 {
		t.Root.values[i] = t.Root.siblings[i+1].values[0]
	}

	return true
}
