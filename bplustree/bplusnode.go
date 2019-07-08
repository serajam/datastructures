package bplustree

import (
	"fmt"
	"strings"
)

type treeNode struct {
	degree           int
	elementsCount    int
	maxElementsCount int
	leaf             bool

	nextLeaf *treeNode
	prevLeaf *treeNode

	values   []int
	siblings []*treeNode
}

func newNode(degree int, leaf bool) *treeNode {
	n := &treeNode{
		values:           make([]int, degree*2-1),
		leaf:             leaf,
		degree:           degree,
		maxElementsCount: degree*2 - 1,
	}

	if !leaf {
		n.siblings = make([]*treeNode, degree*2)
	}

	return n
}

func (n treeNode) value(i int) int {
	return n.values[i]
}

func (n treeNode) lastValue() int {
	return n.values[n.elementsCount]
}

func (n treeNode) keyIndex(val int) int {
	var i int
	for i < n.elementsCount {
		if val == n.values[i] {
			return i
		}
		i++
	}

	return -1
}

func (n *treeNode) mergeLeft(node *treeNode) {
	if !node.leaf {
		n.prependElement(n.siblings[0].values[0])
		for i := n.elementsCount; i > 0; i-- {
			n.siblings[i+node.elementsCount] = n.siblings[i-1]
		}

		for i := 0; i < node.elementsCount+1; i++ {
			n.siblings[i] = node.siblings[i]
		}
	}

	if node.elementsCount == 0 {
		return
	}

	for i := n.elementsCount; i > 0; i-- {
		n.values[i+node.elementsCount-1] = n.values[i-1]
	}

	for i := 0; i < node.elementsCount; i++ {
		n.values[i] = node.values[i]
		n.elementsCount++
	}
}

func (n *treeNode) mergeRight(node *treeNode) {
	if !node.leaf {
		for i := 0; i < node.elementsCount+1; i++ {
			n.siblings[n.elementsCount+1+i] = node.siblings[i]
		}

		n.values[n.elementsCount] = node.siblings[0].values[0]
		n.elementsCount++
	}

	if node.elementsCount == 0 {
		return
	}

	for i := 0; i < node.elementsCount; i++ {
		n.values[n.elementsCount] = node.values[i]
		n.elementsCount++
	}

}

func (n *treeNode) deleteValue(i int) {
	var j int
	for j = i; j < n.elementsCount-1; j++ {
		n.values[j] = n.values[j+1]
	}
	n.values[j] = 0
	n.elementsCount--
}

func (n *treeNode) popElement() int {
	temp := n.values[n.elementsCount-1]
	n.values[n.elementsCount-1] = 0
	n.elementsCount--
	return temp
}

func (n *treeNode) prependElement(val int) {
	for i := n.elementsCount - 1; i >= 0; i-- {
		n.values[i+1] = n.values[i]
	}
	n.values[0] = val
	n.elementsCount++
}

func (n *treeNode) popChildren() *treeNode {
	temp := n.siblings[n.elementsCount]
	n.siblings[n.elementsCount] = nil
	return temp
}

func (n *treeNode) prependChildren(c *treeNode) {
	for i := n.elementsCount; i >= 0; i-- {
		n.siblings[i+1] = n.siblings[i]
	}
	n.siblings[0] = c
}

func (n treeNode) sibling(i int) *treeNode {
	if i > n.elementsCount {
		return nil
	}

	if i < 0 {
		return nil
	}

	d := n.siblings[i]
	if d == nil {
		return nil
	}

	return d
}

func (n *treeNode) deleteSibling(pos int) {

	if pos-1 >= 0 {
		n.siblings[pos-1].nextLeaf = n.siblings[pos].nextLeaf
	}

	if pos+1 <= n.elementsCount {
		n.siblings[pos+1].prevLeaf = n.siblings[pos].prevLeaf
	}

	for pos < n.elementsCount {
		n.siblings[pos] = n.siblings[pos+1]
		if pos < n.elementsCount-1 {
			n.values[pos] = n.values[pos+1]
		}
		pos++
	}

	n.siblings[pos] = nil
	n.values[pos-1] = 0
	n.elementsCount--
}

// traverse TODO implement more human readable traverse
func (n treeNode) traverse(lvl int) string {

	s := "|"

	for _, v := range n.values {
		s += fmt.Sprint(v, "|")
	}
	s += fmt.Sprint(" (", n.elementsCount, ")")

	var c string
	for _, v := range n.siblings {
		if v != nil {
			c += v.traverse(lvl + 1)
		}
	}

	if c == "" {
		return fmt.Sprint(strings.Repeat(" ", lvl*2)) + s
	}

	s = fmt.Sprint(strings.Repeat(" ", lvl*2)) + s + "\n\n" + c + "\n\n"

	return s
}

// searchNext returns child for next insertion lookup
func (n treeNode) searchNext(val int) (i int) {
	for i < n.elementsCount {
		if n.values[i] > val {
			return i
		}
		if n.values[i] == 0 {
			return i
		}

		i++
	}

	return i
}

func (n treeNode) hasEmptyCell() bool {
	return n.elementsCount < n.maxElementsCount
}

func (n treeNode) hasFreeElement() bool {
	return n.elementsCount > n.degree-1
}

func (n *treeNode) splitNode(parent *treeNode, degree int, pos int) *treeNode {
	// if splitting leaf Root
	if parent == nil {
		parent = newNode(degree, false)
	}

	rightNode := newNode(degree, n.leaf)

	mid := degree - 1

	midVal := n.values[mid]
	n.values[mid] = 0
	n.elementsCount--

	split := func() {
		midSlot := 0
		if !n.leaf {
			midSlot = 1
			for i := mid + 1; i <= n.elementsCount+1; i++ {
				rightNode.siblings[i-mid-1] = n.siblings[i]
				n.siblings[i] = nil
			}
		}

		count := n.elementsCount

		for i := mid + 1; i <= count; i++ {
			rightNode.values[i-mid-midSlot] = n.values[i]
			n.values[i] = 0

			rightNode.elementsCount++
			n.elementsCount--
		}

		if n.leaf {
			rightNode.values[0] = midVal
			rightNode.elementsCount++
		}
	}

	// split for root
	if pos == -1 {
		split()

		parent.values[parent.elementsCount] = midVal
		parent.siblings[parent.elementsCount] = n
		parent.siblings[parent.elementsCount+1] = rightNode
		parent.elementsCount++

		if rightNode.leaf {
			n.nextLeaf = rightNode
			rightNode.prevLeaf = n
		}

		return parent
	}

	split()
	for i := parent.elementsCount; i > pos; i-- {
		parent.values[i] = parent.values[i-1]
	}
	parent.values[pos] = midVal

	for i := parent.elementsCount; i > pos; i-- {
		parent.siblings[i+1] = parent.siblings[i]
	}
	parent.siblings[pos+1] = rightNode
	parent.elementsCount++

	if rightNode.leaf {
		rightNode.prevLeaf = parent.siblings[pos]
		if pos+2 <= parent.elementsCount && parent.siblings[pos+2] != nil {
			rightNode.nextLeaf = parent.siblings[pos+2]
			parent.siblings[pos+2].prevLeaf = rightNode
		}

		if parent.siblings[pos] != nil {
			parent.siblings[pos].nextLeaf = rightNode
		}
	}

	return parent
}

func (n *treeNode) insert(val int) {
	for _, v := range n.values {
		if v == val {
			return
		}
	}

	i := len(n.values) - 1
	for i >= 0 {
		if n.values[i] == 0 {
			i--
			continue
		}

		if n.values[i] > val {
			n.values[i+1] = n.values[i]
			i--
			continue
		}

		break
	}

	n.values[i+1] = val
	n.elementsCount++
}

func (n *treeNode) delete(val int) bool {
	if n.leaf {
		i := n.keyIndex(val)
		if i < 0 {
			return false
		}

		n.deleteValue(i)

		return true
	}

	var nextNodeIndex int
	for nextNodeIndex < n.elementsCount {
		if n.value(nextNodeIndex) > val {
			break
		}

		if n.value(nextNodeIndex) == 0 {
			break
		}

		nextNodeIndex++
	}

	nextNode := n.sibling(nextNodeIndex)
	if nextNode == nil {
		return false
	}

	// dive in
	result := nextNode.delete(val)

	// if key not found
	if !result {
		return false
	}

	// check leafs
	if nextNode.leaf {

		// check leaf for underflow
		if nextNode.elementsCount >= nextNode.degree-1 {
			return true
		}

		// borrow element from the right
		if rightSibling := n.sibling(nextNodeIndex + 1); rightSibling != nil && rightSibling.hasFreeElement() {
			nextNode.insert(rightSibling.value(0))
			rightSibling.deleteValue(0)
			n.values[nextNodeIndex] = rightSibling.value(0)

			return true
		}

		// borrow element from the left
		if leftSibling := n.sibling(nextNodeIndex - 1); leftSibling != nil && leftSibling.hasFreeElement() {
			key := leftSibling.popElement()
			nextNode.prependElement(key)
			n.values[nextNodeIndex-1] = key

			return true
		}

		// merging when not enough elements
		// merge left to right
		if rightSibling := n.sibling(nextNodeIndex + 1); rightSibling != nil {
			rightSibling.mergeLeft(nextNode)
			n.deleteSibling(nextNodeIndex)

			return true
		}

		// merge right to left
		if leftSibling := n.sibling(nextNodeIndex - 1); leftSibling != nil {
			leftSibling.mergeRight(nextNode)
			n.deleteSibling(nextNodeIndex)

			return true
		}
	}

	// check internal nodes
	// if treeNode has deleted key, replace it with the copy of next lowest from right sibling
	if i := nextNode.keyIndex(val); i >= 0 {
		nextNode.values[i] = nextNode.sibling(i + 1).value(0)
	}

	// check internal node for underflow
	if nextNode.elementsCount >= nextNode.degree-1 {
		return true
	}

	// borrow element from the right
	if rightSibling := n.sibling(nextNodeIndex + 1); rightSibling != nil && rightSibling.hasFreeElement() {
		children := rightSibling.sibling(0)
		key := rightSibling.value(0)

		rightSibling.deleteSibling(0)
		rightSibling.deleteValue(0)

		nextNode.values[nextNode.elementsCount] = n.value(nextNodeIndex)
		n.values[nextNodeIndex] = key
		nextNode.siblings[nextNode.elementsCount+1] = children
		nextNode.elementsCount++

		return true

	}

	// borrow element from the left
	if leftSibling := n.sibling(nextNodeIndex - 1); leftSibling != nil && leftSibling.hasFreeElement() {
		children := leftSibling.popChildren()
		key := leftSibling.popElement()

		nextNode.prependElement(n.values[nextNodeIndex-1])
		n.values[nextNodeIndex-1] = key
		nextNode.prependChildren(children)

		return true
	}

	// merging when not enough elements
	// merge left to right
	if rightSibling := n.sibling(nextNodeIndex + 1); rightSibling != nil {
		rightSibling.mergeLeft(nextNode)
		n.deleteSibling(nextNodeIndex)

		return true
	}

	// merge right to left
	if leftSIbling := n.sibling(nextNodeIndex - 1); leftSIbling != nil {
		leftSIbling.mergeRight(nextNode)
		n.deleteSibling(nextNodeIndex)

		return true
	}

	return result
}
