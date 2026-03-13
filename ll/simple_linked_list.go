package ll

import (
	"fmt"
	"iter"
	"strings"
)

// SimpleLinkedList is a basic pointer-based doubly-linked list implementation.
//
// This is a straightforward, textbook-style linked list using heap-allocated
// nodes connected by pointers. It is provided for benchmarking comparisons
// against the optimized LinkedList implementation.
//
// The API mirrors LinkedList for direct performance comparison.
type SimpleLinkedList[T comparable] struct {
	head   *simpleNode[T]
	tail   *simpleNode[T]
	length int
}

type simpleNode[T comparable] struct {
	value T
	next  *simpleNode[T]
	prev  *simpleNode[T]
}

// NewSimpleLinkedList creates a new empty SimpleLinkedList.
func NewSimpleLinkedList[T comparable]() *SimpleLinkedList[T] {
	return &SimpleLinkedList[T]{
		head:   nil,
		tail:   nil,
		length: 0,
	}
}

// NewSimpleLinkedListFromSlice creates a new SimpleLinkedList from a slice.
func NewSimpleLinkedListFromSlice[T comparable](slice []T) *SimpleLinkedList[T] {
	list := NewSimpleLinkedList[T]()
	for _, v := range slice {
		list.PushBack(v)
	}
	return list
}

// All returns an iterator over all values in the list.
func (this *SimpleLinkedList[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		current := this.head
		for current != nil {
			if !yield(current.value) {
				return
			}
			current = current.next
		}
	}
}

// AllIndexed returns an iterator over all index-value pairs.
func (this *SimpleLinkedList[T]) AllIndexed() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		index := 0
		current := this.head
		for current != nil {
			if !yield(index, current.value) {
				return
			}
			current = current.next
			index++
		}
	}
}

// Clear removes all elements from the list.
func (this *SimpleLinkedList[T]) Clear() {
	this.head = nil
	this.tail = nil
	this.length = 0
}

// Contains reports whether the list contains the specified value.
func (this *SimpleLinkedList[T]) Contains(value T) bool {
	current := this.head
	for current != nil {
		if current.value == value {
			return true
		}
		current = current.next
	}
	return false
}

// Crystalize is a no-op for SimpleLinkedList (provided for API compatibility).
func (this *SimpleLinkedList[T]) Crystalize() {
	// No-op: pointer-based lists don't benefit from crystallization
}

// ForEach iterates over all values, calling the delegate for each.
func (this *SimpleLinkedList[T]) ForEach(delegate func(value T) bool) {
	current := this.head
	for current != nil {
		if !delegate(current.value) {
			return
		}
		current = current.next
	}
}

// ForEachIndexed iterates over all index-value pairs.
func (this *SimpleLinkedList[T]) ForEachIndexed(delegate func(index int, value T) bool) {
	index := 0
	current := this.head
	for current != nil {
		if !delegate(index, current.value) {
			return
		}
		current = current.next
		index++
	}
}

// Get retrieves the value at the specified index.
func (this *SimpleLinkedList[T]) Get(index int) (value T, ok bool) {
	node := this.nodeAt(index)
	if node == nil {
		return value, false
	}
	return node.value, true
}

// IndexOf returns the index of the first occurrence of the value, or -1.
func (this *SimpleLinkedList[T]) IndexOf(value T) int {
	index := 0
	current := this.head
	for current != nil {
		if current.value == value {
			return index
		}
		current = current.next
		index++
	}
	return -1
}

// InsertAt inserts a value at the specified index.
func (this *SimpleLinkedList[T]) InsertAt(index int, value T) bool {
	if index < 0 || index > this.length {
		return false
	}

	if index == 0 {
		this.PushFront(value)
		return true
	}

	if index == this.length {
		this.PushBack(value)
		return true
	}

	prevNode := this.nodeAt(index - 1)
	if prevNode == nil {
		return false
	}

	newNode := &simpleNode[T]{
		value: value,
		next:  prevNode.next,
		prev:  prevNode,
	}

	if prevNode.next != nil {
		prevNode.next.prev = newNode
	}
	prevNode.next = newNode
	this.length++
	return true
}

// IsEmpty reports whether the list contains no elements.
func (this *SimpleLinkedList[T]) IsEmpty() bool {
	return this.length == 0
}

// Length returns the number of elements in the list.
func (this *SimpleLinkedList[T]) Length() int {
	return this.length
}

// PopBack removes and returns the last element.
func (this *SimpleLinkedList[T]) PopBack() (value T, ok bool) {
	if this.tail == nil {
		return value, false
	}

	value = this.tail.value
	if this.tail.prev != nil {
		this.tail.prev.next = nil
		this.tail = this.tail.prev
	} else {
		this.head = nil
		this.tail = nil
	}
	this.length--
	return value, true
}

// PopFront removes and returns the first element.
func (this *SimpleLinkedList[T]) PopFront() (value T, ok bool) {
	if this.head == nil {
		return value, false
	}

	value = this.head.value
	if this.head.next != nil {
		this.head.next.prev = nil
		this.head = this.head.next
	} else {
		this.head = nil
		this.tail = nil
	}
	this.length--
	return value, true
}

// PushBack adds a value to the end of the list.
func (this *SimpleLinkedList[T]) PushBack(value T) {
	newNode := &simpleNode[T]{
		value: value,
		next:  nil,
		prev:  this.tail,
	}

	if this.tail != nil {
		this.tail.next = newNode
	} else {
		this.head = newNode
	}
	this.tail = newNode
	this.length++
}

// PushFront adds a value to the beginning of the list.
func (this *SimpleLinkedList[T]) PushFront(value T) {
	newNode := &simpleNode[T]{
		value: value,
		next:  this.head,
		prev:  nil,
	}

	if this.head != nil {
		this.head.prev = newNode
	} else {
		this.tail = newNode
	}
	this.head = newNode
	this.length++
}

// RemoveAt removes and returns the element at the specified index.
func (this *SimpleLinkedList[T]) RemoveAt(index int) (value T, ok bool) {
	node := this.nodeAt(index)
	if node == nil {
		return value, false
	}

	value = node.value

	if node.prev != nil {
		node.prev.next = node.next
	} else {
		this.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		this.tail = node.prev
	}

	this.length--
	return value, true
}

// Set replaces the value at the specified index.
func (this *SimpleLinkedList[T]) Set(index int, value T) bool {
	if index == -1 {
		this.PushFront(value)
		return true
	}

	if index == this.length {
		this.PushBack(value)
		return true
	}

	node := this.nodeAt(index)
	if node == nil {
		return false
	}

	node.value = value
	return true
}

// String returns a string representation of the list.
func (this *SimpleLinkedList[T]) String() string {
	builder := strings.Builder{}
	builder.WriteString("[")
	first := true
	current := this.head
	for current != nil {
		if !first {
			builder.WriteString(", ")
		}
		first = false
		fmt.Fprintf(&builder, "%v", current.value)
		current = current.next
	}
	builder.WriteString("]")
	return builder.String()
}

// ToSlice returns a new slice containing all elements.
func (this *SimpleLinkedList[T]) ToSlice() []T {
	slice := make([]T, 0, this.length)
	current := this.head
	for current != nil {
		slice = append(slice, current.value)
		current = current.next
	}
	return slice
}

// --- private methods --- //

func (this *SimpleLinkedList[T]) nodeAt(index int) *simpleNode[T] {
	if index < 0 || index >= this.length {
		return nil
	}

	current := this.head
	for i := 0; i < index; i++ {
		current = current.next
	}
	return current
}
