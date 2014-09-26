package gotrie

import (
	"fmt"
)

type Node struct {
	id        uint64
	isLeaf    bool
	splitPos  uint8
	value     uint64
	leftRight [2]int // left & right children
	next      int    // for leaf
}

type Index struct {
	nodes []Node
	Debug bool
}

func NewIndex(N int) *Index {
	var index = Index{}

	// create a root-indexing node.  It's leftRight[0] will
	// store the actual root node position
	index.nodes = make([]Node, 0, N)
	index.nodes = append(index.nodes, Node{})
	return &index
}

func (index *Index) mkLeafNode(id uint64, value uint64) int {
	var node = Node{
		id:       id,
		isLeaf:   true,
		splitPos: 64,
		value:    value,
	}
	index.nodes = append(index.nodes, node)
	return len(index.nodes) - 1
}
func (index *Index) mkNode(
	val uint64,
	splitPos uint8,
	left, right int) int {
	index.nodes = append(index.nodes, Node{})
	var i = len(index.nodes) - 1
	if i == 1 {
		// mark this as the root
		index.nodes[0].leftRight[0] = i
	}
	index.nodes[i].value = val
	index.nodes[i].splitPos = splitPos
	index.nodes[i].leftRight[0] = left
	index.nodes[i].leftRight[1] = right
	return i
}

func (index *Index) Add(id uint64, value uint64) {
	var i = index.mkLeafNode(id, value)
	index.insertInto(0, 0, i)
}

func (index *Index) insertInto(parent, branch, leaf int) {
	if index.nodes[parent].leftRight[branch] == 0 {
		index.nodes[parent].leftRight[branch] = leaf
		// index.d("Found empty slot.")
		return
	}

	// get the current now
	var (
		leaf_value uint64 = index.nodes[leaf].value
		now        int
		now_value  uint64
		s          uint8
	)

	for {
		now = index.nodes[parent].leftRight[branch]
		now_value = index.nodes[now].value
		s = LeadZeros(leaf_value ^ now_value)

		if s < index.nodes[now].splitPos {
			// create new node
			var bit = TestBit_Int(leaf_value, s)
			var new_node = index.mkNode(now_value, s, 0, 0)
			index.nodes[new_node].leftRight[bit] = leaf
			index.nodes[new_node].leftRight[1-bit] = now
			// update the parent link
			index.nodes[parent].leftRight[branch] = new_node
			break
		} else {
			if index.nodes[now].isLeaf {
				// add leaf to the linked list
				index.nodes[leaf].next = index.nodes[now].next
				index.nodes[now].next = leaf
				break
			} else {
				// insert under now-node
				var bit = TestBit_Int(leaf_value, index.nodes[now].splitPos)
				parent = now
				branch = bit
			}
		}
	}

	return
}

func (index *Index) Print() {
	fmt.Printf("Compressed trie with %d nodes\n", len(index.nodes))
	for i := range index.nodes {
		fmt.Printf("%.2d: ", i)
		index.PrintNode(i)
	}
}

func (index *Index) Graphviz() {
	fmt.Println("digraph G {")
	for i, n := range index.nodes {
		// left branch
		if n.leftRight[0] > 0 {
			fmt.Printf("n%d -> n%d [label=0];\n", i, n.leftRight[0])
		}
		// right branch
		if n.leftRight[1] > 0 {
			fmt.Printf("n%d -> n%d [label=1];\n", i, n.leftRight[1])
		}
		var val = Uint64_string(n.value)
		if n.isLeaf {
			val = "[" + val + "]"
		} else {

		}
		fmt.Printf("n%d [label=\"%d(%d)=%s\"];\n", i, i, n.splitPos, val)
	}
	fmt.Println("}")
}

func (index *Index) PrintNode(n int) {
	var node = index.nodes[n]
	if node.isLeaf {
		fmt.Printf("LEAF=%x(split:%d)[%.2d, %.2d]\n", node.value, node.splitPos, node.leftRight[0], node.leftRight[1])
	} else {
		fmt.Printf("(split:%d)[%.2d, %.2d]\n", node.splitPos, node.leftRight[0], node.leftRight[1])
	}
}

func (index *Index) Len() int {
	return len(index.nodes) - 1
}
func (index *Index) d(s string, a ...interface{}) {
	if index.Debug {
		fmt.Printf(s+"\n", a...)
	}
}

type Searcher struct {
	index      *Index
	queue      []int
	head, tail int
	O          int
}

func NewSearcher(index *Index) *Searcher {
	var searcher = Searcher{}
	searcher.queue = make([]int, len(index.nodes))
	searcher.index = index

	return &searcher
}
func (this *Searcher) Push(pos int) {
	this.queue[this.tail] = pos
	this.tail += 1
}
func (this *Searcher) Pop() int {
	if this.QueueLen() > 0 {
		x := this.queue[this.head]
		this.head += 1
		return x
	} else {
		panic("Popping an empty stack")
	}
}
func (this *Searcher) QueueLen() int {
	return this.tail - this.head
}

func (this *Searcher) Search(value uint64, r uint8) (count int) {
	var (
		diff uint8
	)
	this.head = 0
	this.tail = 0

	if len(this.index.nodes) < 2 {
		return
	}

	rootPos := this.index.nodes[0].leftRight[0]
	this.Push(rootPos)
	for this.QueueLen() > 0 {
		node := &this.index.nodes[this.Pop()]
		diff = PopCountPartial(value^node.value, node.splitPos)

		if diff > r {
			continue
		}
		if this.O > 1 {
			if diff+64-node.splitPos <= r {
				count += 1 // todo this should be the count under node
				continue
			}
		}

		if node.isLeaf {
			count += 1
		} else {
			if node.leftRight[0] > 0 {
				this.Push(node.leftRight[0])
			}
			if node.leftRight[1] > 0 {
				this.Push(node.leftRight[1])
			}
		}
	}

	return
}
