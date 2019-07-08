package bplustree

import (
	"log"
	"math/rand"
	"reflect"
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
	node := newNode(2, false)

	//rand.Seed(time.Now().Unix())
	rand.Seed(1)

	vals := []*treeNode{
		&treeNode{},
		&treeNode{},
		&treeNode{degree: 2},
	}

	for i := 0; i < len(vals); i++ {
		node.siblings[i] = vals[i]
	}

	node.elementsCount = 2

	c := node.popChildren()

	if c == nil {
		log.Fatalln("Should pop existing child")
	}

	if node.siblings[node.elementsCount] != nil {
		log.Fatalln("Should delete popped siblings")
	}

	if c.degree != vals[2].degree {
		log.Fatalln("Wrong siblings popped")
	}
}

func TestPrependChildren(t *testing.T) {
	node := newNode(3, false)

	//rand.Seed(time.Now().Unix())
	rand.Seed(1)

	vals := []*treeNode{
		&treeNode{},
		&treeNode{},
		&treeNode{},
	}

	for i := 0; i < len(vals); i++ {
		node.siblings[i] = vals[i]
	}

	node.elementsCount = 2

	node.prependChildren(&treeNode{degree: 10})

	if node.siblings[0].degree != 10 {
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

func Test_treeNode_splitNode(t *testing.T) {

	type want struct {
		result        *treeNode
		siblingsCount int
		valuesCount   int
	}

	type fields struct {
		degree int
		values []int
		leaf   bool
	}

	type args struct {
		degree   int
		pos      int
		values   []int
		children []*treeNode
		count    int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			"Should split root leaf with degree 2",
			fields{degree: 2, values: []int{15, 5, 32}, leaf: true},
			args{degree: 2, pos: -1, children: nil, values: nil},
			want{result: &treeNode{values: []int{15, 0, 0}, elementsCount: 1}, siblingsCount: 2, valuesCount: 1},
		},

		{
			"Should split root leaf wit degree 7",
			fields{degree: 7, values: []int{90, 5, 32, 11, 23, 53, 12, 34, 3, 14, 10, 50, 60}, leaf: true},
			args{degree: 7, pos: -1, values: nil, children: nil},
			want{result: &treeNode{values: []int{23, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, elementsCount: 1}, siblingsCount: 2, valuesCount: 1},
		},

		{
			"Should split internal node",
			fields{degree: 7, values: []int{90, 500, 32, 110, 23, 53, 120, 34, 100, 140, 100, 5100, 600}, leaf: true},
			args{
				count:  1,
				values: []int{23, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				children: []*treeNode{
					nil,
					{degree: 7, values: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 0}, elementsCount: 12},
					nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				},
				degree: 7,
				pos:    1,
			},
			want{result: &treeNode{degree: 7, values: []int{23, 110, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, elementsCount: 2}, siblingsCount: 3, valuesCount: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := newNode(tt.fields.degree, tt.fields.leaf)
			for _, v := range tt.fields.values {
				n.insert(v)
			}

			var parent *treeNode
			if tt.args.pos > -1 {
				parent = &treeNode{}
				parent.degree = tt.args.degree
				parent.values = tt.args.values
				parent.siblings = tt.args.children
				parent.elementsCount = tt.args.count

				parent.siblings[0] = n
			}

			got := n.splitNode(parent, tt.args.degree, tt.args.pos)
			if !reflect.DeepEqual(got.values, tt.want) {
				for i := range got.values {
					if got.values[i] != tt.want.result.values[i] {
						t.Errorf("treeNode.splitNode().valules = %v, want %v", got.values, tt.want.result.values)
					}
				}

			}

			if got.elementsCount != tt.want.result.elementsCount {
				t.Errorf("treeNode.elementsCount = %v, want %v", got.elementsCount, tt.want.result.elementsCount)
			}

			siblings := 0
			for _, s := range got.siblings {
				if s != nil {
					siblings++
				}
			}

			if siblings != tt.want.siblingsCount {
				t.Errorf("treeNode.sibling count = %v, want %v", siblings, tt.want.siblingsCount)
			}

			values := 0
			for _, s := range got.values {
				if s != 0 {
					values++
				}
			}

			if values != tt.want.valuesCount {
				t.Errorf("treeNode.sibling count = %v, want %v", values, tt.want.valuesCount)
			}

		})
	}
}
