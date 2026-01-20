package ll

import (
	"fmt"
	"iter"
	"strings"

	_ll "github.com/ConanHorus/fds/ll/internal"
	"github.com/ConanHorus/fds/mm"
)

const (
	nilIndex = -1

	minCapacity = 16
)

// LinkedList is a generic doubly-linked list implementation optimized for
// performance and memory efficiency.
//
// It provides O(1) insertion and removal at both ends, O(n) indexed access,
// and O(n) search operations. The list automatically grows as needed.
//
// The type parameter T must be comparable to support operations like Contains
// and IndexOf.
//
// Incremental Untangling:
//
// The list performs automatic incremental untangling during iteration. Each
// time the list is traversed (via All, ForEach, Get, etc.), at most one
// out-of-order node pair is swapped to improve memory layout. Over time, this
// brings nodes closer to sequential memory order, improving cache locality.
//
// Crystallization:
//
// For read-heavy workloads, call Crystalize() to fully untangle the list in
// one O(n²) operation. Once crystallized, the list skips incremental
// untangling during reads (since it's already optimized). Any subsequent write
// operation automatically marks the list as non-crystallized, resuming
// incremental untangling.
//
// Tradeoffs:
//   - Memory Efficiency vs. Speed: Uses an internal buddy allocator for memory
//     management, reducing allocation overhead at the cost of slight memory
//     fragmentation.
//   - Cache Locality: Nodes may not be contiguous in memory after many
//     insertions and deletions. Incremental untangling gradually improves
//     layout during normal use. For immediate cache-friendly access, use
//     Crystalize() or ToSlice().
//   - Read/Write Side Effects: Read operations may modify internal node
//     positions (though not values) unless the list is crystallized. Concurrent
//     access requires external synchronization.
type LinkedList[T comparable] struct {
	headIndex int
	tailIndex int
	length    int

	nodes  []_ll.Node
	values []T

	allocator mm.BuddyAllocator

	crystallized bool
}

// NewLinkedList creates a new empty LinkedList.
//
// Returns:
//   - A pointer to the newly created LinkedList.
func NewLinkedList[T comparable]() *LinkedList[T] {
	return newLinkedList[T](minCapacity)
}

// NewLinkedListFromSlice creates a new LinkedList populated with elements from
// the provided slice.
//
// Elements are added in order, so the first element of the slice becomes the
// head of the list.
//
// Parameters:
//   - slice: The slice of elements to populate the list with. May be nil or
//     empty.
//
// Returns:
//   - A pointer to the newly created LinkedList containing all slice elements.
func NewLinkedListFromSlice[T comparable](slice []T) *LinkedList[T] {
	list := newLinkedList[T](max(len(slice)+len(slice)/4+1, minCapacity))
	for _, v := range slice {
		list.PushBack(v)
	}

	return list
}

func newLinkedList[T comparable](capacity int) *LinkedList[T] {
	allocator := *mm.NewBuddyAllocator(mm.WithInitialCapacity(uint64(capacity)), mm.WithMaxOrder(3))
	nodes := make([]_ll.Node, allocator.Capacity())
	for i := 0; i < len(nodes); i++ {
		nodes[i].Index = -1
	}

	return &LinkedList[T]{
		headIndex: -1,
		tailIndex: -1,
		length:    0,

		nodes:  nodes,
		values: make([]T, allocator.Capacity()),

		allocator: allocator,
	}
}

// All returns an iterator over all values in the list from head to tail.
//
// This method supports Go 1.23+ range-over-function syntax and can be used
// with early break.
//
// Side Effect: Iteration may swap one out-of-order node pair to improve memory
// layout. Values and their order are unchanged, but internal memory positions
// may shift.
//
// Examples:
//   - for value := range list.All() { ... }
//   - for value := range list.All() { if done { break } }
//
// Returns:
//   - An iter.Seq[T] that yields each value in order.
func (this *LinkedList[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, innerIndex := range this.allIndexesForward(!this.crystallized) {
			if !yield(this.values[innerIndex]) {
				return
			}
		}
	}
}

// AllIndexed returns an iterator over all index-value pairs in the list from
// head to tail.
//
// This method supports Go 1.23+ range-over-function syntax and can be used
// with early break.
//
// Examples:
//   - for index, value := range list.AllIndexed() { ... }
//
// Returns:
//   - An iter.Seq2[int, T] that yields each index and value in order.
func (this *LinkedList[T]) AllIndexed() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for outerIndex, innerIndex := range this.allIndexesForward(!this.crystallized) {
			if !yield(outerIndex, this.values[innerIndex]) {
				return
			}
		}
	}
}

// Clear removes all elements from the list.
//
// After calling Clear, the list will be empty (Length() == 0, IsEmpty() == true).
// The underlying memory is reset but not deallocated, allowing for efficient
// reuse.
func (this *LinkedList[T]) Clear() {
	this.headIndex = nilIndex
	this.tailIndex = nilIndex
	this.length = 0
	this.allocator.ClearAll()
}

// Contains reports whether the list contains the specified value.
//
// Uses the == operator for comparison, which is why T must be comparable.
//
// Parameters:
//   - value: The value to search for.
//
// Returns:
//   - ok: true if the value is found, false otherwise.
func (this *LinkedList[T]) Contains(value T) (ok bool) {
	for _, innerIndex := range this.allIndexesForward(!this.crystallized) {
		if this.values[innerIndex] == value {
			return true
		}
	}

	return false
}

// Crystalize optimizes the internal structure of the list for read operations.
//
// This method fully untangles the list, rearranging nodes so they are stored
// in sequential memory order matching their logical order. After crystallization:
//   - Read operations benefit from optimal cache locality
//   - Incremental untangling is disabled (already optimized)
//   - The list is marked as \"crystallized\"
//
// Any subsequent write operation (PushBack, PushFront, InsertAt, RemoveAt, Set,
// Clear) automatically marks the list as non-crystallized, resuming incremental
// untangling during reads.
//
// Complexity: O(n²) worst case, where n is the number of elements. For lists
// that are already mostly ordered, performance approaches O(n).
//
// Use Cases:
//   - Call after bulk insertions, before a read-heavy phase
//   - Call before serialization or benchmarking
//   - Not needed if incremental untangling is sufficient
func (this *LinkedList[T]) Crystalize() {
	insertAt := 0
	at := this.headIndex
	for at != nilIndex {
		if at == insertAt {
			at = this.nodes[insertAt].NextIndex
			insertAt++
			continue
		}

		this.swapNodes(insertAt, at)
		at = this.nodes[insertAt].NextIndex
		insertAt++
	}

	this.crystallized = true
}

// ForEach iterates over all values in the list, calling the delegate function
// for each value.
//
// Iteration proceeds from head to tail and stops early if the delegate returns
// false.
//
// Parameters:
//   - delegate: A function called for each value. Return true to continue
//     iteration, false to stop.
func (this *LinkedList[T]) ForEach(delegate func(value T) (ok bool)) {
	for _, innerIndex := range this.allIndexesForward(!this.crystallized) {
		if !delegate(this.values[innerIndex]) {
			return
		}
	}
}

// ForEachIndexed iterates over all index-value pairs in the list, calling the
// delegate function for each pair.
//
// Iteration proceeds from head to tail and stops early if the delegate returns
// false.
//
// Parameters:
//   - delegate: A function called for each index-value pair. Return true to
//     continue iteration, false to stop.
func (this *LinkedList[T]) ForEachIndexed(delegate func(index int, value T) (ok bool)) {
	for outerIndex, innerIndex := range this.allIndexesForward(!this.crystallized) {
		if !delegate(outerIndex, this.values[innerIndex]) {
			return
		}
	}
}

// Get retrieves the value at the specified index.
//
// This is an O(n) operation as it requires traversing the list from the head.
//
// Parameters:
//   - index: The zero-based index of the element to retrieve. Must be in the
//     range [0, Length()).
//
// Returns:
//   - value: The value at the specified index, or the zero value if not found.
//   - ok: true if the index is valid, false otherwise.
func (this *LinkedList[T]) Get(index int) (value T, ok bool) {
	innerIndex := this.fromOuterIndexToInner(index, !this.crystallized)
	if innerIndex == -1 {
		return value, false
	}

	return this.values[innerIndex], true
}

// InsertAt inserts a value at the specified index.
//
// All elements at and after the index are shifted toward the tail. Inserting
// at index 0 is equivalent to PushFront. Inserting at index Length() is
// equivalent to PushBack.
//
// Parameters:
//   - index: The zero-based index at which to insert. Must be in the range
//     [0, Length()].
//   - value: The value to insert.
//
// Returns:
//   - ok: true if the insertion succeeded, false if the index is out of range.
func (this *LinkedList[T]) InsertAt(index int, value T) (ok bool) {
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

	prevInnerIndex := this.fromOuterIndexToInner(index-1, false)
	if prevInnerIndex == -1 {
		return false
	}

	nextInnerIndex := this.nodes[prevInnerIndex].NextIndex
	newInnerIndex := this.allocateNode(value)
	this.nodes[newInnerIndex].PreviousIndex = prevInnerIndex
	this.nodes[newInnerIndex].NextIndex = nextInnerIndex
	this.nodes[prevInnerIndex].NextIndex = newInnerIndex
	this.nodes[nextInnerIndex].PreviousIndex = newInnerIndex
	return true
}

// IndexOf returns the index of the first occurrence of the specified value.
//
// Uses the == operator for comparison, which is why T must be comparable.
//
// Parameters:
//   - value: The value to search for.
//
// Returns:
//   - index: The zero-based index of the first occurrence, or -1 if not found.
func (this *LinkedList[T]) IndexOf(value T) (index int) {
	for outerIndex, innerIndex := range this.allIndexesForward(!this.crystallized) {
		if this.values[innerIndex] == value {
			return outerIndex
		}
	}

	return -1
}

// IsEmpty reports whether the list contains no elements.
//
// Returns:
//   - empty: true if the list has no elements, false otherwise.
func (this *LinkedList[T]) IsEmpty() (empty bool) {
	return this.length == 0
}

// Length returns the number of elements in the list.
//
// Returns:
//   - length: The number of elements currently in the list.
func (this *LinkedList[T]) Length() (length int) {
	return this.length
}

// PopBack removes and returns the last element of the list.
//
// Returns:
//   - value: The removed value, or the zero value if the list is empty.
//   - ok: true if an element was removed, false if the list was empty.
func (this *LinkedList[T]) PopBack() (value T, ok bool) {
	return this.RemoveAt(this.length - 1)
}

// PopFront removes and returns the first element of the list.
//
// Returns:
//   - value: The removed value, or the zero value if the list is empty.
//   - ok: true if an element was removed, false if the list was empty.
func (this *LinkedList[T]) PopFront() (value T, ok bool) {
	return this.RemoveAt(0)
}

// PushBack adds a value to the end of the list.
//
// This is an O(1) operation.
//
// Parameters:
//   - value: The value to add.
func (this *LinkedList[T]) PushBack(value T) {
	this.Set(this.length, value)
}

// PushFront adds a value to the beginning of the list.
//
// This is an O(1) operation.
//
// Parameters:
//   - value: The value to add.
func (this *LinkedList[T]) PushFront(value T) {
	this.Set(-1, value)
}

// RemoveAt removes and returns the element at the specified index.
//
// All elements after the index are shifted toward the head.
//
// Parameters:
//   - index: The zero-based index of the element to remove. Must be in the
//     range [0, Length()).
//
// Returns:
//   - value: The removed value, or the zero value if the index is invalid.
//   - ok: true if an element was removed, false if the index is out of range.
func (this *LinkedList[T]) RemoveAt(index int) (value T, ok bool) {
	this.crystallized = false
	innerIndex := this.fromOuterIndexToInner(index, false)
	if innerIndex == -1 {
		return value, false
	}

	node := this.nodes[innerIndex]
	value = this.values[innerIndex]
	prevIndex := node.PreviousIndex
	nextIndex := node.NextIndex

	if prevIndex != nilIndex {
		this.nodes[prevIndex].NextIndex = nextIndex
	} else {
		this.headIndex = nextIndex
	}

	if nextIndex != nilIndex {
		this.nodes[nextIndex].PreviousIndex = prevIndex
	} else {
		this.tailIndex = prevIndex
	}

	this.nodes[innerIndex].Index = -1
	this.values[innerIndex] = *new(T)
	this.allocator.Free(uint64(innerIndex), 1)
	this.length--
	return value, true
}

// Set replaces the value at the specified index.
//
// Special index values:
//   - index == -1: Equivalent to PushFront (adds to beginning)
//   - index == Length(): Equivalent to PushBack (adds to end)
//
// For all other indices, the existing value at that position is replaced.
//
// Parameters:
//   - index: The zero-based index of the element to set. Must be -1, in the
//     range [0, Length()), or exactly Length().
//   - value: The value to set.
//
// Returns:
//   - ok: true if the operation succeeded, false if the index is out of range.
func (this *LinkedList[T]) Set(index int, value T) (ok bool) {
	if index == -1 {
		index := this.allocateNode(value)
		if this.tailIndex == nilIndex {
			this.headIndex = index
			this.tailIndex = index
			return true
		}

		oldHeadIndex := this.headIndex
		this.nodes[oldHeadIndex].PreviousIndex = index
		this.nodes[index].NextIndex = oldHeadIndex
		this.headIndex = index
		return true
	}

	if index == this.length {
		index := this.allocateNode(value)
		if this.headIndex == nilIndex {
			this.headIndex = index
			this.tailIndex = index
			return true
		}

		oldTailIndex := this.tailIndex
		this.nodes[oldTailIndex].NextIndex = index
		this.nodes[index].PreviousIndex = oldTailIndex
		this.tailIndex = index
		return true
	}

	innerIndex := this.fromOuterIndexToInner(index, false)
	if innerIndex == -1 {
		return false
	}

	this.values[innerIndex] = value
	return true
}

// String returns a string representation of the list.
//
// The format is "[elem1, elem2, elem3]" with elements separated by ", ".
// An empty list returns "[]".
//
// Implements fmt.Stringer.
//
// Returns:
//   - str: The string representation of the list.
func (this *LinkedList[T]) String() (str string) {
	builder := strings.Builder{}
	builder.WriteString("[")
	first := true
	for value := range this.All() {
		if !first {
			builder.WriteString(", ")
		}

		first = false
		fmt.Fprintf(&builder, "%v", value)
	}

	builder.WriteString("]")
	return builder.String()
}

// ToSlice returns a new slice containing all elements in the list.
//
// The returned slice is a copy; modifications to it do not affect the list.
// Elements appear in the same order as in the list (head to tail).
//
// Returns:
//   - slice: A new slice containing all list elements.
func (this *LinkedList[T]) ToSlice() (slice []T) {
	slice = make([]T, 0, this.length)
	for value := range this.All() {
		slice = append(slice, value)
	}

	return slice
}

// --- private methods --- //

func (this *LinkedList[T]) allocateNode(value T) (index int) {
	this.crystallized = false
	allocationIndex, grew, ok := this.allocator.Allocate(1)
	if !ok {
		panic("allocation failed")
	}

	if grew {
		oldCapacity := len(this.nodes)
		newCapacity := this.allocator.Capacity()
		newNodes := make([]_ll.Node, newCapacity)
		newValues := make([]T, newCapacity)
		copy(newNodes, this.nodes)
		copy(newValues, this.values)
		this.nodes = newNodes
		this.values = newValues

		for i := oldCapacity; i < len(this.nodes); i++ {
			this.nodes[i].Index = -1
		}
	}

	index = int(allocationIndex)
	this.values[index] = value
	this.nodes[index] = _ll.Node{
		Index:         index,
		NextIndex:     nilIndex,
		PreviousIndex: nilIndex,
	}

	this.length++
	return index
}

func (this *LinkedList[T]) fromOuterIndexToInner(index int, untangle bool) (out int) {
	if index < 0 || index >= this.length {
		return -1
	}

	if this.crystallized {
		return index
	}

	mid := this.length / 2
	if index <= mid {
		for outerIndex, innerIndex := range this.allIndexesForward(untangle) {
			if outerIndex == index {
				return innerIndex
			}
		}
	} else {
		for outerIndex, innerIndex := range this.allIndexesReverse(untangle) {
			if outerIndex == index {
				return innerIndex
			}
		}
	}

	return -1
}

func (this *LinkedList[T]) allIndexesForward(untangle bool) iter.Seq2[int, int] {
	if this.crystallized {
		return func(yield func(int, int) bool) {
			for i := 0; i < this.length; i++ {
				if !yield(i, i) {
					return
				}
			}
		}
	}

	return func(yield func(int, int) bool) {
		outerIndex := 0
		innerIndex := this.headIndex
		if innerIndex == nilIndex {
			return
		}

		for innerIndex != nilIndex {
			if untangle && outerIndex > 0 {
				previousIndex := this.nodes[innerIndex].PreviousIndex
				if previousIndex > innerIndex {
					this.swapInUseNodes(previousIndex, innerIndex)
					innerIndex = previousIndex
					untangle = false
				}
			}

			if !yield(outerIndex, innerIndex) {
				return
			}

			innerIndex = this.nodes[innerIndex].NextIndex
			outerIndex++
		}
	}
}

func (this *LinkedList[T]) allIndexesReverse(untangle bool) iter.Seq2[int, int] {
	if this.crystallized {
		return func(yield func(int, int) bool) {
			for i := this.length - 1; i >= 0; i-- {
				if !yield(i, i) {
					return
				}
			}
		}
	}

	return func(yield func(int, int) bool) {
		outerIndex := this.length - 1
		innerIndex := this.tailIndex
		if innerIndex == nilIndex {
			return
		}

		for innerIndex != nilIndex {
			if untangle && outerIndex > 0 {
				previousIndex := this.nodes[innerIndex].PreviousIndex
				if previousIndex < innerIndex {
					this.swapInUseNodes(previousIndex, innerIndex)
					innerIndex = previousIndex
					untangle = false
				}
			}

			if !yield(outerIndex, innerIndex) {
				return
			}

			innerIndex = this.nodes[innerIndex].PreviousIndex
			outerIndex--
		}
	}
}

func (this *LinkedList[T]) moveNodeToLocation(fromIndex int, toIndex int) {
	this.crystallized = false

	this.allocator.Free(uint64(fromIndex), 1)
	if !this.allocator.AllocateAt(uint64(toIndex), 1) {
		panic("failed to allocate at specific location during compaction")
	}

	this.nodes[toIndex] = this.nodes[fromIndex]
	this.nodes[toIndex].Index = toIndex
	this.values[toIndex] = this.values[fromIndex]

	if this.nodes[toIndex].PreviousIndex != nilIndex {
		this.nodes[this.nodes[toIndex].PreviousIndex].NextIndex = toIndex
	} else {
		this.headIndex = toIndex
	}

	if this.nodes[toIndex].NextIndex != nilIndex {
		this.nodes[this.nodes[toIndex].NextIndex].PreviousIndex = toIndex
	} else {
		this.tailIndex = toIndex
	}

	this.nodes[fromIndex].Index = -1
	this.values[fromIndex] = *new(T)
}

func (this *LinkedList[T]) swapNodes(indexA int, indexB int) {
	this.crystallized = false
	nodeA := &this.nodes[indexA]
	nodeB := &this.nodes[indexB]

	if nodeA.Index != -1 && nodeB.Index != -1 {
		this.swapInUseNodes(indexA, indexB)
		return
	}

	if nodeA.Index == -1 && nodeB.Index != -1 {
		this.moveNodeToLocation(indexB, indexA)
	} else if nodeA.Index != -1 && nodeB.Index == -1 {
		this.moveNodeToLocation(indexA, indexB)
	}
}

func (this *LinkedList[T]) swapInUseNodes(indexA int, indexB int) {
	this.crystallized = false
	beforePointer := this.nodes[indexA].PreviousIndex
	if beforePointer == nilIndex {
		this.headIndex = indexB
	} else {
		this.nodes[beforePointer].NextIndex = indexB
	}

	afterPointer := this.nodes[indexB].NextIndex
	if afterPointer == nilIndex {
		this.tailIndex = indexA
	} else {
		this.nodes[afterPointer].PreviousIndex = indexA
	}

	nodeA := &this.nodes[indexA]
	nodeB := &this.nodes[indexB]

	nodeA.NextIndex = afterPointer
	nodeA.PreviousIndex = indexB
	nodeB.NextIndex = indexA
	nodeB.PreviousIndex = beforePointer

	this.values[indexA], this.values[indexB] = this.values[indexB], this.values[indexA]
}
