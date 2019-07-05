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
	children []*treeNode
}

func newNode(degree int, leaf bool) *treeNode {
	return &treeNode{
		values:           make([]int, degree*2-1),
		children:         make([]*treeNode, degree*2),
		leaf:             leaf,
		degree:           degree,
		maxElementsCount: degree*2 - 1,
	}
}

func (n treeNode) value(i int) int {
	return n.values[i]
}

func (n treeNode) lastKey() int {
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
		n.prependElement(n.children[0].values[0])
		for i := n.elementsCount; i > 0; i-- {
			n.children[i+node.elementsCount] = n.children[i-1]
		}

		for i := 0; i < node.elementsCount+1; i++ {
			n.children[i] = node.children[i]
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
			n.children[n.elementsCount+1+i] = node.children[i]
		}

		n.values[n.elementsCount] = node.children[0].values[0]
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

func (n *treeNode) deleteKey(i int) {
	var j int
	for j = i; j < n.elementsCount-1; j++ {
		n.values[j] = n.values[j+1]
	}
	n.values[j] = 0
	n.elementsCount--
}

func (n *treeNode) deleteChildren(i int) {
	var j int
	for j = i; j <= n.elementsCount; j++ {
		n.children[j] = n.children[j+1]
	}
	n.children[j] = nil

	//n.values = append(n.values[:i], n.values[i+1:]...)
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
	temp := n.children[n.elementsCount]
	n.children[n.elementsCount] = nil
	return temp
}

func (n *treeNode) prependChildren(c *treeNode) {
	for i := n.elementsCount; i >= 0; i-- {
		n.children[i+1] = n.children[i]
	}
	n.children[0] = c
}

func (n treeNode) sibling(i int) *treeNode {
	if i > n.elementsCount {
		return nil
	}

	if i < 0 {
		return nil
	}

	d := n.children[i]
	if d == nil {
		return nil
	}

	return d
}

func (n *treeNode) deleteSibling(pos int) {

	for pos < n.elementsCount {
		n.children[pos] = n.children[pos+1]
		if pos < n.elementsCount-1 {
			n.values[pos] = n.values[pos+1]
		}
		pos++
	}

	n.children[pos] = nil
	n.values[pos-1] = 0
	n.elementsCount--
}

func (n treeNode) traverse(lvl int) string {

	s := "|"

	for _, v := range n.values {
		s += fmt.Sprint(v, "|")
	}
	s += fmt.Sprint(" (", n.elementsCount, ")")

	var c string
	for _, v := range n.children {
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

func (n *treeNode) splitLeaf(parent *treeNode, degree int, pos int) *treeNode {
	// if splitting leaf Root
	if parent == nil {
		parent = newNode(degree, false)
	}
	rightLeaf := newNode(degree, true)

	mid := degree - 1

	// @TODO implement linked list
	//rightLeaf.prevLeaf = n
	//rightLeaf.nextLeaf = n.nextLeaf
	//n.nextLeaf = rightLeaf

	midVal := n.values[mid]
	n.values[mid] = 0
	n.elementsCount--

	shift := func() {
		for i := mid + 1; i < len(n.values); i++ {
			rightLeaf.values[i-mid] = n.values[i]
			n.values[i] = 0

			rightLeaf.elementsCount++
			n.elementsCount--
		}

		rightLeaf.values[0] = midVal
		rightLeaf.elementsCount++
	}

	// split for root
	if pos == -1 {
		shift()

		parent.values[parent.elementsCount] = midVal
		parent.children[parent.elementsCount] = n
		parent.children[parent.elementsCount+1] = rightLeaf
		parent.elementsCount++

		return parent
	}

	shift()
	for i := parent.elementsCount; i > pos; i-- {
		parent.values[i] = parent.values[i-1]
	}
	parent.values[pos] = midVal

	for i := parent.elementsCount; i > pos; i-- {
		parent.children[i+1] = parent.children[i]
	}
	parent.children[pos+1] = rightLeaf
	parent.elementsCount++

	return parent
}

func (n *treeNode) splitNode(parent *treeNode, degree int, pos int) *treeNode {
	// if splitting leaf Root
	if parent == nil {
		parent = newNode(degree, false)
	}
	rightNode := newNode(degree, false)

	mid := degree - 1

	midVal := n.values[mid]
	n.values[mid] = 0
	n.elementsCount--

	shift := func() {
		for i := mid + 1; i <= n.elementsCount+1; i++ {
			rightNode.children[i-mid-1] = n.children[i]
			n.children[i] = nil
		}

		to := n.elementsCount

		for i := mid + 1; i <= to; i++ {
			rightNode.values[i-mid-1] = n.values[i]
			n.values[i] = 0

			rightNode.elementsCount++
			n.elementsCount--
		}
	}

	// split for root
	if pos == -1 {
		shift()

		parent.values[parent.elementsCount] = midVal
		parent.children[parent.elementsCount] = n
		parent.children[parent.elementsCount+1] = rightNode
		parent.elementsCount++

		return parent
	}

	shift()
	for i := parent.elementsCount; i > pos; i-- {
		parent.values[i] = parent.values[i-1]
	}
	parent.values[pos] = midVal

	for i := parent.elementsCount; i > pos; i-- {
		parent.children[i+1] = parent.children[i]
	}
	parent.children[pos+1] = rightNode
	parent.elementsCount++

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

		n.deleteKey(i)

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
			rightSibling.deleteKey(0)
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

		rightSibling.deleteChildren(0)
		rightSibling.deleteKey(0)

		nextNode.values[nextNode.elementsCount] = n.value(nextNodeIndex)
		n.values[nextNodeIndex] = key
		nextNode.children[nextNode.elementsCount+1] = children
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