package tree

import "cmp"

type bstNode[T cmp.Ordered] struct {
	value T
	left  *bstNode[T]
	right *bstNode[T]
}

// BST is a basic binary search tree with no self-balancing.
//
// It provides O(log n) average-case operations but degrades to O(n) in the
// worst case when elements are inserted in sorted order.
//
// The type parameter T must satisfy cmp.Ordered.
type BST[T cmp.Ordered] struct {
	root   *bstNode[T]
	length int
}

// NewBST creates a new empty binary search tree.
//
// Returns:
//   - A pointer to the newly created BST.
func NewBST[T cmp.Ordered]() *BST[T] {
	return &BST[T]{}
}

// Contains reports whether the tree contains the specified value.
//
// Parameters:
//   - value: The value to search for.
//
// Returns:
//   - ok: true if the value is found, false otherwise.
func (this *BST[T]) Contains(value T) (ok bool) {
	return this.find(this.root, value) != nil
}

// Min returns the minimum value in the tree.
//
// Returns:
//   - value: The minimum value, or the zero value if the tree is empty.
//   - ok: true if the tree is non-empty, false otherwise.
func (this *BST[T]) Min() (value T, ok bool) {
	if this.root == nil {
		return value, false
	}

	node := this.root
	for node.left != nil {
		node = node.left
	}

	return node.value, true
}

// Max returns the maximum value in the tree.
//
// Returns:
//   - value: The maximum value, or the zero value if the tree is empty.
//   - ok: true if the tree is non-empty, false otherwise.
func (this *BST[T]) Max() (value T, ok bool) {
	if this.root == nil {
		return value, false
	}

	node := this.root
	for node.right != nil {
		node = node.right
	}

	return node.value, true
}

// Length returns the number of elements in the tree.
//
// Returns:
//   - length: The number of elements currently in the tree.
func (this *BST[T]) Length() (length int) {
	return this.length
}

// IsEmpty reports whether the tree contains no elements.
//
// Returns:
//   - empty: true if the tree has no elements, false otherwise.
func (this *BST[T]) IsEmpty() (empty bool) {
	return this.length == 0
}

// Insert adds a value to the tree.
//
// Parameters:
//   - value: The value to insert.
//
// Returns:
//   - ok: true if the value was inserted, false if it already exists.
func (this *BST[T]) Insert(value T) (ok bool) {
	this.root, ok = this.insert(this.root, value)
	if ok {
		this.length++
	}

	return ok
}

// Delete removes a value from the tree.
//
// Parameters:
//   - value: The value to remove.
//
// Returns:
//   - ok: true if the value was found and removed, false otherwise.
func (this *BST[T]) Delete(value T) (ok bool) {
	this.root, ok = this.delete(this.root, value)
	if ok {
		this.length--
	}

	return ok
}

// Clear removes all elements from the tree.
func (this *BST[T]) Clear() {
	this.root = nil
	this.length = 0
}

// --- private methods --- //

func (this *BST[T]) find(node *bstNode[T], value T) *bstNode[T] {
	if node == nil {
		return nil
	}

	if value < node.value {
		return this.find(node.left, value)
	}

	if value > node.value {
		return this.find(node.right, value)
	}

	return node
}

func (this *BST[T]) insert(node *bstNode[T], value T) (*bstNode[T], bool) {
	if node == nil {
		return &bstNode[T]{value: value}, true
	}

	if value < node.value {
		node.left, _ = this.insert(node.left, value)
		return node, true
	}

	if value > node.value {
		node.right, _ = this.insert(node.right, value)
		return node, true
	}

	return node, false
}

func (this *BST[T]) delete(node *bstNode[T], value T) (*bstNode[T], bool) {
	if node == nil {
		return nil, false
	}

	var ok bool
	if value < node.value {
		node.left, ok = this.delete(node.left, value)
		return node, ok
	}

	if value > node.value {
		node.right, ok = this.delete(node.right, value)
		return node, ok
	}

	if node.left == nil {
		return node.right, true
	}

	if node.right == nil {
		return node.left, true
	}

	successor := node.right
	for successor.left != nil {
		successor = successor.left
	}

	node.value = successor.value
	node.right, _ = this.delete(node.right, successor.value)
	return node, true
}
