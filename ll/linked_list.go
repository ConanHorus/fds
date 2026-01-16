package ll

import "github.com/ConanHorus/fds/mm"

// LinkedList is a generic linked list implementation that supports dynamic
// resizing using a BuddyAllocator for memory management.
//
// It provides efficient insertion and removal of elements at both ends of the
// list as well as at arbitrary positions.
//
// Tradeoffs:
//   - Memory Efficiency vs. Speed: Using a BuddyAllocator allows for efficient
//     memory management and reduces fragmentation, but may introduce some
//     overhead compared to using native slices or arrays.
//   - Flexibility vs. Complexity: This implementation provides flexibility in
//     terms of dynamic resizing and generic types, but adds complexity to the
//     codebase compared to simpler linked list implementations.
type LinkedList[T any] struct {
	head   uint64 // zero indicates nil value
	tail   uint64 // zero indicates nil value
	size   uint64
	nodes  []linkedListNode // zero is never used
	values []T              // zero is never used

	memory mm.BuddyAllocator
}

// LinkedListNode represents a node in the LinkedList.
type LinkedListNode[T any] struct {
	Value      T
	valueIndex uint64

	linkedList *LinkedList[T]
	next       uint64
	previous   uint64
}

type linkedListNode struct {
	next       uint64 // zero indicates nil value
	previous   uint64 // zero indicates nil value
	valueIndex uint64
}

// NewLinkedList creates a new LinkedList with an optional initial size.
//
// Tradeoffs:
//   - Memory Efficiency vs. Speed: Using a BuddyAllocator allows for efficient
//     memory management and reduces fragmentation, but may introduce some
//     overhead compared to using native slices or arrays.
//   - Flexibility vs. Complexity: This implementation provides flexibility in
//     terms of dynamic resizing and generic types, but adds complexity to the
//     codebase compared to simpler linked list implementations.
//
// Parameters:
//   - initialSize: An optional parameter specifying the initial size of the
//     linked list. If not provided, a default size is used.
//
// Returns:
//   - A pointer to the newly created LinkedList.
func NewLinkedList[T any](initialSize ...int) *LinkedList[T] {
	const defaultInitialSize = 16

	size := defaultInitialSize
	if len(initialSize) > 0 {
		size = max(initialSize[0], size)
	}

	return &LinkedList[T]{
		head:   0,
		tail:   0,
		size:   0,
		nodes:  make([]linkedListNode, size+1), // index 0 is never used
		values: make([]T, size+1),              // index 0 is never used
		memory: *mm.NewBuddyAllocator(
			mm.WithInitialCapacity(uint64(size)),
		),
	}
}

// Append adds a new value to the end of the linked list.
//
// Parameters:
//   - value: The value to append to the linked list.
func (this *LinkedList[T]) Append(value T) {
	newNodeIndex, grew, _ := this.memory.Allocate(1)
	newNodeIndex++
	if grew {
		newNodes := make([]linkedListNode, this.memory.Capacity()+1)
		newValues := make([]T, this.memory.Capacity()+1)
		copy(newNodes, this.nodes)
		copy(newValues, this.values)
		this.nodes = newNodes
		this.values = newValues
	}

	this.values[newNodeIndex] = value
	newNode := linkedListNode{
		next:       0,
		previous:   this.tail,
		valueIndex: newNodeIndex,
	}

	this.nodes[newNodeIndex] = newNode
	if this.tail != 0 {
		this.nodes[this.tail].next = newNodeIndex
	}

	this.tail = newNodeIndex
	if this.head == 0 {
		this.head = newNodeIndex
	}

	this.size++
}

// Head returns the first node in the linked list.
//
// Returns:
//   - A pointer to the first LinkedListNode, or nil if the list is empty.
func (this *LinkedList[T]) Head() *LinkedListNode[T] {
	if this.head == 0 {
		return nil
	}

	node := &this.nodes[this.head]
	return &LinkedListNode[T]{
		Value:      this.values[node.valueIndex],
		valueIndex: node.valueIndex,
		linkedList: this,
		next:       node.next,
		previous:   node.previous,
	}
}

// PopHead removes and returns the first node from the linked list.
//
// Returns:
//   - A pointer to the removed LinkedListNode, or nil if the list is empty.
func (this *LinkedList[T]) PopHead() *LinkedListNode[T] {
	if this.head == 0 {
		return nil
	}

	nodeIndex := this.head
	node := &this.nodes[nodeIndex]
	this.memory.Free(nodeIndex-1, 1)

	this.head = node.next
	if this.head != 0 {
		this.nodes[this.head].previous = 0
	} else {
		this.tail = 0
	}

	this.size--

	return &LinkedListNode[T]{
		Value:      this.values[node.valueIndex],
		linkedList: this,
		valueIndex: node.valueIndex,
		next:       node.next,
		previous:   node.previous,
	}
}

// PopTail removes and returns the last node from the linked list.
//
// Returns:
//   - A pointer to the removed LinkedListNode, or nil if the list is empty.
func (this *LinkedList[T]) PopTail() *LinkedListNode[T] {
	if this.tail == 0 {
		return nil
	}

	nodeIndex := this.tail
	node := &this.nodes[nodeIndex]
	this.memory.Free(nodeIndex-1, 1)

	this.tail = node.previous
	if this.tail != 0 {
		this.nodes[this.tail].next = 0
	} else {
		this.head = 0
	}

	this.size--

	return &LinkedListNode[T]{
		Value:      this.values[node.valueIndex],
		linkedList: this,
		valueIndex: node.valueIndex,
		next:       node.next,
		previous:   node.previous,
	}
}

// Prepend adds a new value to the beginning of the linked list.
//
// Parameters:
//   - value: The value to prepend to the linked list.
func (this *LinkedList[T]) Prepend(value T) {
	newNodeIndex, grew, _ := this.memory.Allocate(1)
	newNodeIndex++
	if grew {
		newNodes := make([]linkedListNode, this.memory.Capacity()+1)
		newValues := make([]T, this.memory.Capacity()+1)
		copy(newNodes, this.nodes)
		copy(newValues, this.values)
		this.nodes = newNodes
		this.values = newValues
	}

	this.values[newNodeIndex] = value
	newNode := linkedListNode{
		next:       this.head,
		previous:   0,
		valueIndex: newNodeIndex,
	}

	this.nodes[newNodeIndex] = newNode
	if this.head != 0 {
		this.nodes[this.head].previous = newNodeIndex
	}

	this.head = newNodeIndex
	if this.tail == 0 {
		this.tail = newNodeIndex
	}

	this.size++
}

// Remove removes a specified node from the linked list.
//
// Parameters:
//   - node: A pointer to the LinkedListNode to remove.
func (this *LinkedList[T]) Remove(node *LinkedListNode[T]) {
	if node == nil || node.linkedList != this {
		return
	}

	if node.previous != 0 {
		this.nodes[node.previous].next = node.next
	} else {
		this.head = node.next
	}

	if node.next != 0 {
		this.nodes[node.next].previous = node.previous
	} else {
		this.tail = node.previous
	}

	this.memory.Free(this.findNodeIndex(node), 1)
	this.size--
}

// Tail returns the last node in the linked list.
//
// Returns:
//   - A pointer to the last LinkedListNode, or nil if the list is empty.
func (this *LinkedList[T]) Tail() *LinkedListNode[T] {
	if this.tail == 0 {
		return nil
	}

	node := &this.nodes[this.tail]
	return &LinkedListNode[T]{
		Value:      this.values[node.valueIndex],
		linkedList: this,
		valueIndex: node.valueIndex,
		next:       node.next,
		previous:   node.previous,
	}
}

// Size returns the number of elements in the linked list.
//
// Returns:
//   - The size of the linked list as a uint64.
func (this *LinkedList[T]) Size() uint64 {
	return this.size
}

// InsertNext inserts a new value after the current node.
//
// Parameters:
//   - value: The value to insert after the current node.
func (this *LinkedListNode[T]) InsertNext(value T) {
	newNodeIndex, grew, _ := this.linkedList.memory.Allocate(1)
	newNodeIndex++
	if grew {
		newNodes := make([]linkedListNode, this.linkedList.memory.Capacity()+1)
		newValues := make([]T, this.linkedList.memory.Capacity()+1)
		copy(newNodes, this.linkedList.nodes)
		copy(newValues, this.linkedList.values)
		this.linkedList.nodes = newNodes
		this.linkedList.values = newValues
	}

	this.linkedList.values[newNodeIndex] = value
	thisIndex := this.linkedList.findNodeIndex(this)
	newNode := linkedListNode{
		next:       this.next,
		previous:   thisIndex,
		valueIndex: newNodeIndex,
	}

	this.linkedList.nodes[newNodeIndex] = newNode
	if this.next != 0 {
		this.linkedList.nodes[this.next].previous = newNodeIndex
	} else {
		this.linkedList.tail = newNodeIndex
	}

	this.linkedList.nodes[thisIndex].next = newNodeIndex
	this.next = newNodeIndex

	this.linkedList.size++
}

// InsertPrevious inserts a new value before the current node.
//
// Parameters:
//   - value: The value to insert before the current node.
func (this *LinkedListNode[T]) InsertPrevious(value T) {
	newNodeIndex, grew, _ := this.linkedList.memory.Allocate(1)
	newNodeIndex++
	if grew {
		newNodes := make([]linkedListNode, this.linkedList.memory.Capacity()+1)
		newValues := make([]T, this.linkedList.memory.Capacity()+1)
		copy(newNodes, this.linkedList.nodes)
		copy(newValues, this.linkedList.values)
		this.linkedList.nodes = newNodes
		this.linkedList.values = newValues
	}

	this.linkedList.values[newNodeIndex] = value
	thisIndex := this.linkedList.findNodeIndex(this)
	newNode := linkedListNode{
		next:       thisIndex,
		previous:   this.previous,
		valueIndex: newNodeIndex,
	}

	this.linkedList.nodes[newNodeIndex] = newNode
	if this.previous != 0 {
		this.linkedList.nodes[this.previous].next = newNodeIndex
	} else {
		this.linkedList.head = newNodeIndex
	}

	this.linkedList.nodes[thisIndex].previous = newNodeIndex
	this.previous = newNodeIndex

	this.linkedList.size++
}

// Next returns the next node in the linked list.
//
// Returns:
//   - A pointer to the next LinkedListNode, or nil if there is no next node.
func (this *LinkedListNode[T]) Next() *LinkedListNode[T] {
	if this.next == 0 {
		return nil
	}

	node := &this.linkedList.nodes[this.next]
	return &LinkedListNode[T]{
		Value:      this.linkedList.values[node.valueIndex],
		linkedList: this.linkedList,
		next:       node.next,
		previous:   node.previous,
	}
}

// Previous returns the previous node in the linked list.
//
// Returns:
//   - A pointer to the previous LinkedListNode, or nil if there is no previous
//     node.
func (this *LinkedListNode[T]) Previous() *LinkedListNode[T] {
	if this.previous == 0 {
		return nil
	}

	node := &this.linkedList.nodes[this.previous]
	return &LinkedListNode[T]{
		Value:      this.linkedList.values[node.valueIndex],
		linkedList: this.linkedList,
		next:       node.next,
		previous:   node.previous,
	}
}

// --- private methods --- //

func (this *LinkedList[T]) findNodeIndex(node *LinkedListNode[T]) uint64 {
	currentIndex := this.head
	for currentIndex != 0 {
		currentNode := &this.nodes[currentIndex]
		if currentNode.valueIndex == node.valueIndex {
			return currentIndex
		}

		currentIndex = currentNode.next
	}

	return 0
}
