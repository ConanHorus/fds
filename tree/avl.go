package tree

import "cmp"

type avlNode[T cmp.Ordered] struct {
	value  T
	left   *avlNode[T]
	right  *avlNode[T]
	height int
}

// AVL is a self-balancing binary search tree that maintains balance by ensuring
// the height difference between left and right subtrees is at most 1.
//
// All operations are O(log n) in both average and worst case.
//
// The type parameter T must satisfy cmp.Ordered.
type AVL[T cmp.Ordered] struct {
	root   *avlNode[T]
	length int
}

// NewAVL creates a new empty AVL tree.
//
// Returns:
//   - A pointer to the newly created AVL tree.
func NewAVL[T cmp.Ordered]() *AVL[T] {
	return &AVL[T]{}
}

// Contains reports whether the tree contains the specified value.
//
// Parameters:
//   - value: The value to search for.
//
// Returns:
//   - ok: true if the value is found, false otherwise.
func (this *AVL[T]) Contains(value T) (ok bool) {
	return this.find(this.root, value) != nil
}

// Min returns the minimum value in the tree.
//
// Returns:
//   - value: The minimum value, or the zero value if the tree is empty.
//   - ok: true if the tree is non-empty, false otherwise.
func (this *AVL[T]) Min() (value T, ok bool) {
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
func (this *AVL[T]) Max() (value T, ok bool) {
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
func (this *AVL[T]) Length() (length int) {
	return this.length
}

// IsEmpty reports whether the tree contains no elements.
//
// Returns:
//   - empty: true if the tree has no elements, false otherwise.
func (this *AVL[T]) IsEmpty() (empty bool) {
	return this.length == 0
}

// Insert adds a value to the tree.
//
// Parameters:
//   - value: The value to insert.
//
// Returns:
//   - ok: true if the value was inserted, false if it already exists.
func (this *AVL[T]) Insert(value T) (ok bool) {
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
func (this *AVL[T]) Delete(value T) (ok bool) {
	this.root, ok = this.delete(this.root, value)
	if ok {
		this.length--
	}

	return ok
}

// Clear removes all elements from the tree.
func (this *AVL[T]) Clear() {
	this.root = nil
	this.length = 0
}

// --- private methods --- //

func (this *AVL[T]) find(node *avlNode[T], value T) *avlNode[T] {
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

func (this *AVL[T]) height(node *avlNode[T]) int {
	if node == nil {
		return -1
	}

	return node.height
}

func (this *AVL[T]) balanceFactor(node *avlNode[T]) int {
	if node == nil {
		return 0
	}

	return this.height(node.left) - this.height(node.right)
}

func (this *AVL[T]) updateHeight(node *avlNode[T]) {
	node.height = 1 + max(this.height(node.left), this.height(node.right))
}

func (this *AVL[T]) rotateRight(node *avlNode[T]) *avlNode[T] {
	left := node.left
	node.left = left.right
	left.right = node
	this.updateHeight(node)
	this.updateHeight(left)
	return left
}

func (this *AVL[T]) rotateLeft(node *avlNode[T]) *avlNode[T] {
	right := node.right
	node.right = right.left
	right.left = node
	this.updateHeight(node)
	this.updateHeight(right)
	return right
}

func (this *AVL[T]) rebalance(node *avlNode[T]) *avlNode[T] {
	this.updateHeight(node)
	balance := this.balanceFactor(node)

	if balance > 1 {
		if this.balanceFactor(node.left) < 0 {
			node.left = this.rotateLeft(node.left)
		}

		return this.rotateRight(node)
	}

	if balance < -1 {
		if this.balanceFactor(node.right) > 0 {
			node.right = this.rotateRight(node.right)
		}

		return this.rotateLeft(node)
	}

	return node
}

func (this *AVL[T]) insert(node *avlNode[T], value T) (*avlNode[T], bool) {
	if node == nil {
		return &avlNode[T]{value: value, height: 0}, true
	}

	var ok bool
	if value < node.value {
		node.left, ok = this.insert(node.left, value)
	} else if value > node.value {
		node.right, ok = this.insert(node.right, value)
	} else {
		return node, false
	}

	return this.rebalance(node), ok
}

func (this *AVL[T]) delete(node *avlNode[T], value T) (*avlNode[T], bool) {
	if node == nil {
		return nil, false
	}

	var ok bool
	if value < node.value {
		node.left, ok = this.delete(node.left, value)
	} else if value > node.value {
		node.right, ok = this.delete(node.right, value)
	} else {
		ok = true
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
	}

	if node == nil {
		return nil, ok
	}

	return this.rebalance(node), ok
}
