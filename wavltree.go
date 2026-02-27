package wavlgo

import (
	"fmt"
)

type keytype interface {
	LessThan(interface{}) bool
	Equal(interface{}) bool
}

type Iterator interface {
	IsNil() bool
	Next() Iterator
	Key() keytype
	Value() valuetype
	Set(v valuetype)
}

type valuetype interface{}

type node struct {
	left, right, parent *node
	key                 keytype
	value               valuetype
	rank                int
}

// Tree is a struct of red-black tree
type Tree struct {
	root *node
	size int
}

// NewTree returns a new rbtree
func NewTree() *Tree {
	return &Tree{}
}

// Find finds the node and return its value
func (t *Tree) Find(key keytype) interface{} {
	n := t.root.findnode(key)
	if n != nil {
		return n.value
	}
	return nil
}

// FindIt finds the node and return it as a iterator
func (t *Tree) FindIt(key keytype) Iterator {
	return t.root.findnode(key)
}

// Empty checks whether the rbtree is empty
func (t *Tree) Empty() bool {
	if t.root == nil {
		return true
	}
	return false
}

// Iterator creates the rbtree's iterator that points to the minmum node
func (t *Tree) Iterator() Iterator {
	return t.root.minimum()
}

// Size returns the size of the rbtree
func (t *Tree) Size() int {
	return t.size
}

// Clear destroys the rbtree
func (t *Tree) Clear() {
	t.root = nil
	t.size = 0
}

// Insert inserts the key-value pair into the rbtree
func (t *Tree) Insert(key keytype, value valuetype) {
	t.root = t.root.insert(&node{
		key:   key,
		value: value,
	}, nil)
	t.size++
}

// Delete deletes the node by key
func (t *Tree) Delete(key keytype) {
	n := t.root.findnode(key)
	if n == nil {
		return
	}

	p := n.parent
	var r *node
	var infix bool
	for _, c := range []*node{n.left, n.right} {
		if c != nil {
			if r != nil {
				infix = true
			}
			r = c
		}
	}

	if infix {
		//	r = n.successor()
		//	c := r.right
		//	n.left.parent = r
		//	n.right.parent = r
		//	//px.
	}

	if r != nil {
		r.parent = p
	}
	if p != nil {
		if p.right == n {
			p.right = r
		} else {
			p.left = r
		}
	} else {
		t.root = r
	}
	t.size--
}

// Preorder prints the tree in pre order
func (t *Tree) Preorder() {
	fmt.Println("preorder begin!")
	if t.root != nil {
		t.root.preorder()
	}
	fmt.Println("preorder end!")
}

// findnode finds the node by key and return it,if not exists return nil
func (n *node) findnode(key keytype) *node {
	if n == nil {
		return nil
	}
	if n.key.Equal(key) {
		return n
	}
	if key.LessThan(n.key) {
		return n.left.findnode(key)
	}
	return n.right.findnode(key)
}

// transplant transplants the subtree u and v
func (t *Tree) transplant(u, v *node) {}

// Next returns the node's successor as an iterator
func (n *node) Next() Iterator {
	return n.successor()
}

func (n *node) preorder() {
	fmt.Printf("(%v %v)", n.key, n.value)
	if n.left != nil {
		fmt.Printf("%v's left child is ", n.key)
		n.left.preorder()
	}
	if n.right != nil {
		fmt.Printf("%v's right child is ", n.key)
		n.right.preorder()
	}
}

func (n *node) insert(c *node, p *node) *node {
	if c == nil {
		return n
	}
	if n == nil {
		c.parent = p
		return c
	}
	if c.key.LessThan(n.key) {
		n.left = n.left.insert(c, n)
	} else {
		n.right = n.right.insert(c, n)
	}
	r := n.rebalance()
	return r
}

// successor returns the successor of the node
func (x *node) successor() *node {
	if x.right != nil {
		return x.right.minimum()
	}
	y := x.parent
	for y != nil && x == y.right {
		x = y
		y = x.parent
	}
	return y
}

// minimum finds the minimum node of subtree n.
func (n *node) minimum() *node {
	if n == nil {
		return nil
	}
	for n.left != nil {
		n = n.left
	}
	return n
}

// maximum finds the maximum node of subtree n.
func (n *node) maximum() *node {
	if n == nil {
		return nil
	}
	for n.right != nil {
		n = n.right
	}
	return n
}

func (n *node) Rank() int {
	if n == nil {
		return -1
	}
	return n.rank
}

func (n *node) IsNil() bool {
	return n == nil
}

func (n *node) Key() keytype {
	return n.key
}

func (n *node) Value() valuetype {
	return n.value
}

func (n *node) Set(v valuetype) {
	n.value = v
}

func (n *node) rebalance() *node {
	if n == nil {
		return n
	}

	nl := n.left.rebalance()
	nr := n.right.rebalance()
	n.refreshRank()

	for {
		if nr.Rank() > nl.Rank()+1 {
			nrr := n.right.right
			nrl := n.right.left
			if nrl.Rank() > nrr.Rank() {
				n.right = nrl.rotateRight(nr)
				nr = n.right
			} else {
				nl = n
				n = nr.rotateLeft(n)
				nr = n.right
			}
		} else if nl.Rank() > nr.Rank()+1 {
			nlr := n.left.right
			nll := n.left.left
			if nlr.Rank() > nll.Rank() {
				n.left = nlr.rotateRight(nl)
				nl = n.left
			} else {
				nr = n
				n = nl.rotateLeft(n)
				nl = n.left
			}
		} else {
			n.refreshRank()
			break
		}
	}
	return n
}

func (n *node) refreshRank() {
	n.rank = max(n.left.Rank(), n.right.Rank()) + 1
}

func (n *node) rotateRight(p *node) *node {
	p.left = n.right
	if n.right != nil {
		n.right.parent = p
	}
	n.parent = p.parent
	n.right = p
	p.parent = n
	p.refreshRank()
	n.refreshRank()
	return n
}

func (n *node) rotateLeft(p *node) *node {
	p.right = n.left
	if n.left != nil {
		n.left.parent = p
	}
	n.parent = p.parent
	n.left = p
	p.parent = n
	p.refreshRank()
	n.refreshRank()
	return n
}
