package tree

import "cmp"

const (
	red   = true
	black = false
)

type color = bool

type rbNode[T cmp.Ordered] struct {
	value  T
	left   *rbNode[T]
	right  *rbNode[T]
	parent *rbNode[T]
	color  color
}

// RedBlack is a self-balancing binary search tree that uses node coloring to
// maintain balance. It guarantees O(log n) operations in all cases.
//
// Properties maintained:
//   - Every node is either red or black
//   - The root is always black
//   - No two consecutive red nodes (red node cannot have a red child)
//   - Every path from root to nil has the same number of black nodes
//
// The type parameter T must satisfy cmp.Ordered.
type RedBlack[T cmp.Ordered] struct {
	root   *rbNode[T]
	length int
}

// NewRedBlack creates a new empty red-black tree.
//
// Returns:
//   - A pointer to the newly created RedBlack tree.
func NewRedBlack[T cmp.Ordered]() *RedBlack[T] {
	return &RedBlack[T]{}
}

// Contains reports whether the tree contains the specified value.
//
// Parameters:
//   - value: The value to search for.
//
// Returns:
//   - ok: true if the value is found, false otherwise.
func (this *RedBlack[T]) Contains(value T) (ok bool) {
	return this.find(value) != nil
}

// Min returns the minimum value in the tree.
//
// Returns:
//   - value: The minimum value, or the zero value if the tree is empty.
//   - ok: true if the tree is non-empty, false otherwise.
func (this *RedBlack[T]) Min() (value T, ok bool) {
	if this.root == nil {
		return value, false
	}

	node := this.minimum(this.root)
	return node.value, true
}

// Max returns the maximum value in the tree.
//
// Returns:
//   - value: The maximum value, or the zero value if the tree is empty.
//   - ok: true if the tree is non-empty, false otherwise.
func (this *RedBlack[T]) Max() (value T, ok bool) {
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
func (this *RedBlack[T]) Length() (length int) {
	return this.length
}

// IsEmpty reports whether the tree contains no elements.
//
// Returns:
//   - empty: true if the tree has no elements, false otherwise.
func (this *RedBlack[T]) IsEmpty() (empty bool) {
	return this.length == 0
}

// Insert adds a value to the tree.
//
// Parameters:
//   - value: The value to insert.
//
// Returns:
//   - ok: true if the value was inserted, false if it already exists.
func (this *RedBlack[T]) Insert(value T) (ok bool) {
	node := &rbNode[T]{value: value, color: red}

	if this.root == nil {
		this.root = node
		this.root.color = black
		this.length++
		return true
	}

	current := this.root
	for {
		if value < current.value {
			if current.left == nil {
				current.left = node
				node.parent = current
				break
			}

			current = current.left
		} else if value > current.value {
			if current.right == nil {
				current.right = node
				node.parent = current
				break
			}

			current = current.right
		} else {
			return false
		}
	}

	this.fixInsert(node)
	this.length++
	return true
}

// Delete removes a value from the tree.
//
// Parameters:
//   - value: The value to remove.
//
// Returns:
//   - ok: true if the value was found and removed, false otherwise.
func (this *RedBlack[T]) Delete(value T) (ok bool) {
	node := this.find(value)
	if node == nil {
		return false
	}

	this.deleteNode(node)
	this.length--
	return true
}

// Clear removes all elements from the tree.
func (this *RedBlack[T]) Clear() {
	this.root = nil
	this.length = 0
}

// --- private methods --- //

func (this *RedBlack[T]) find(value T) *rbNode[T] {
	node := this.root
	for node != nil {
		if value < node.value {
			node = node.left
		} else if value > node.value {
			node = node.right
		} else {
			return node
		}
	}

	return nil
}

func (this *RedBlack[T]) minimum(node *rbNode[T]) *rbNode[T] {
	for node.left != nil {
		node = node.left
	}

	return node
}

func (this *RedBlack[T]) rotateLeft(node *rbNode[T]) {
	right := node.right
	node.right = right.left
	if right.left != nil {
		right.left.parent = node
	}

	right.parent = node.parent
	if node.parent == nil {
		this.root = right
	} else if node == node.parent.left {
		node.parent.left = right
	} else {
		node.parent.right = right
	}

	right.left = node
	node.parent = right
}

func (this *RedBlack[T]) rotateRight(node *rbNode[T]) {
	left := node.left
	node.left = left.right
	if left.right != nil {
		left.right.parent = node
	}

	left.parent = node.parent
	if node.parent == nil {
		this.root = left
	} else if node == node.parent.right {
		node.parent.right = left
	} else {
		node.parent.left = left
	}

	left.right = node
	node.parent = left
}

func (this *RedBlack[T]) fixInsert(node *rbNode[T]) {
	for node != this.root && node.parent.color == red {
		if node.parent == node.parent.parent.left {
			uncle := node.parent.parent.right
			if uncle != nil && uncle.color == red {
				node.parent.color = black
				uncle.color = black
				node.parent.parent.color = red
				node = node.parent.parent
			} else {
				if node == node.parent.right {
					node = node.parent
					this.rotateLeft(node)
				}

				node.parent.color = black
				node.parent.parent.color = red
				this.rotateRight(node.parent.parent)
			}
		} else {
			uncle := node.parent.parent.left
			if uncle != nil && uncle.color == red {
				node.parent.color = black
				uncle.color = black
				node.parent.parent.color = red
				node = node.parent.parent
			} else {
				if node == node.parent.left {
					node = node.parent
					this.rotateRight(node)
				}

				node.parent.color = black
				node.parent.parent.color = red
				this.rotateLeft(node.parent.parent)
			}
		}
	}

	this.root.color = black
}

func (this *RedBlack[T]) transplant(u *rbNode[T], v *rbNode[T]) {
	if u.parent == nil {
		this.root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}

	if v != nil {
		v.parent = u.parent
	}
}

func (this *RedBlack[T]) deleteNode(node *rbNode[T]) {
	var child *rbNode[T]
	var childParent *rbNode[T]
	originalColor := node.color

	if node.left == nil {
		child = node.right
		childParent = node.parent
		this.transplant(node, node.right)
	} else if node.right == nil {
		child = node.left
		childParent = node.parent
		this.transplant(node, node.left)
	} else {
		successor := this.minimum(node.right)
		originalColor = successor.color
		child = successor.right
		if successor.parent == node {
			childParent = successor
		} else {
			childParent = successor.parent
			this.transplant(successor, successor.right)
			successor.right = node.right
			successor.right.parent = successor
		}

		this.transplant(node, successor)
		successor.left = node.left
		successor.left.parent = successor
		successor.color = node.color
	}

	if originalColor == black {
		this.fixDelete(child, childParent)
	}
}

func (this *RedBlack[T]) fixDelete(node *rbNode[T], parent *rbNode[T]) {
	for node != this.root && (node == nil || node.color == black) {
		if node == parent.left {
			sibling := parent.right
			if sibling.color == red {
				sibling.color = black
				parent.color = red
				this.rotateLeft(parent)
				sibling = parent.right
			}

			if (sibling.left == nil || sibling.left.color == black) &&
				(sibling.right == nil || sibling.right.color == black) {
				sibling.color = red
				node = parent
				parent = node.parent
			} else {
				if sibling.right == nil || sibling.right.color == black {
					if sibling.left != nil {
						sibling.left.color = black
					}

					sibling.color = red
					this.rotateRight(sibling)
					sibling = parent.right
				}

				sibling.color = parent.color
				parent.color = black
				if sibling.right != nil {
					sibling.right.color = black
				}

				this.rotateLeft(parent)
				node = this.root
			}
		} else {
			sibling := parent.left
			if sibling.color == red {
				sibling.color = black
				parent.color = red
				this.rotateRight(parent)
				sibling = parent.left
			}

			if (sibling.right == nil || sibling.right.color == black) &&
				(sibling.left == nil || sibling.left.color == black) {
				sibling.color = red
				node = parent
				parent = node.parent
			} else {
				if sibling.left == nil || sibling.left.color == black {
					if sibling.right != nil {
						sibling.right.color = black
					}

					sibling.color = red
					this.rotateLeft(sibling)
					sibling = parent.left
				}

				sibling.color = parent.color
				parent.color = black
				if sibling.left != nil {
					sibling.left.color = black
				}

				this.rotateRight(parent)
				node = this.root
			}
		}
	}

	if node != nil {
		node.color = black
	}
}

func (this *RedBlack[T]) nodeColor(node *rbNode[T]) color {
	if node == nil {
		return black
	}

	return node.color
}
