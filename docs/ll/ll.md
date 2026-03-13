# LinkedList Package (`ll`)

The `ll` package provides a high-performance, generic doubly-linked list implementation optimized for speed and memory efficiency.

## Overview

`LinkedList[T comparable]` is a doubly-linked list that uses an internal buddy allocator for memory management. This design provides:

- **O(1)** insertion and removal at both ends (`PushFront`, `PushBack`, `PopFront`, `PopBack`)
- **O(n)** indexed access and modification (`Get`, `Set`, `InsertAt`, `RemoveAt`)
- **O(n)** search operations (`Contains`, `IndexOf`)
- **Automatic growth** with no user-visible memory management

## Quick Start

```go
import "github.com/ConanHorus/fds/ll"

// Create an empty list
list := ll.NewLinkedList[int]()

// Create from a slice
list := ll.NewLinkedListFromSlice([]int{1, 2, 3, 4, 5})

// Basic operations
list.PushBack(6)      // Add to end
list.PushFront(0)     // Add to beginning
value, ok := list.PopFront()  // Remove from beginning
value, ok := list.PopBack()   // Remove from end

// Iteration (Go 1.23+)
for value := range list.All() {
    fmt.Println(value)
}

for index, value := range list.AllIndexed() {
    fmt.Printf("%d: %v\n", index, value)
}
```

## API Reference

### Constructors

| Function | Description |
|----------|-------------|
| `NewLinkedList[T]()` | Creates a new empty list |
| `NewLinkedListFromSlice[T](slice []T)` | Creates a list from a slice |

### Query Methods

| Method | Complexity | Description |
|--------|------------|-------------|
| `Length() int` | O(1) | Returns the number of elements |
| `IsEmpty() bool` | O(1) | Reports if the list is empty |
| `Contains(value T) bool` | O(n) | Reports if value exists in the list |
| `IndexOf(value T) int` | O(n) | Returns index of first occurrence, or -1 |
| `Get(index int) (T, bool)` | O(n) | Returns value at index |

### Mutation Methods

| Method | Complexity | Description |
|--------|------------|-------------|
| `PushFront(value T)` | O(1) | Adds value to the beginning |
| `PushBack(value T)` | O(1) | Adds value to the end |
| `PopFront() (T, bool)` | O(1) | Removes and returns first element |
| `PopBack() (T, bool)` | O(1) | Removes and returns last element |
| `InsertAt(index int, value T) bool` | O(n) | Inserts at specified index |
| `RemoveAt(index int) (T, bool)` | O(n) | Removes element at index |
| `Set(index int, value T) bool` | O(n) | Replaces value at index |
| `Clear()` | O(1) | Removes all elements |

### Optimization Methods

| Method | Complexity | Description |
|--------|------------|-------------|
| `Crystalize()` | O(n²) | Fully untangles the list for optimal read performance |

### Iteration Methods

| Method | Description |
|--------|-------------|
| `All() iter.Seq[T]` | Iterator over all values (Go 1.23+) |
| `AllIndexed() iter.Seq2[int, T]` | Iterator over index-value pairs |
| `ForEach(func(T) bool)` | Callback iteration with early-exit |
| `ForEachIndexed(func(int, T) bool)` | Indexed callback iteration |

### Conversion Methods

| Method | Description |
|--------|-------------|
| `ToSlice() []T` | Returns a copy of elements as a slice |
| `String() string` | Returns string representation `[a, b, c]` |

## Type Constraint

The type parameter `T` must satisfy the `comparable` constraint. This is required for `Contains` and `IndexOf` operations which use the `==` operator.

```go
// Works with any comparable type
listInt := ll.NewLinkedList[int]()
listStr := ll.NewLinkedList[string]()

type Point struct { X, Y int }
listPt := ll.NewLinkedList[Point]()
```

## Performance Considerations

### Incremental Untangling

The list performs **automatic incremental untangling** during iteration. Each time the list is traversed (via `All()`, `ForEach()`, `Get()`, etc.), at most one out-of-order node pair is swapped to improve memory layout.

**How it works:**
- During iteration, if a node's predecessor has a higher memory index than the current node, they are swapped
- Only one swap occurs per iteration pass
- Over many operations, nodes gradually migrate toward sequential memory order

**Benefits:**
- Cache locality improves automatically over time
- No explicit compaction step required
- Amortized cost: O(1) extra work per iteration

**Tradeoffs:**
- Read operations have side effects (internal memory positions change)
- Not thread-safe without external synchronization, even for reads
- May take many iterations to fully untangle a heavily fragmented list

### Crystallization

For read-heavy workloads, call `Crystalize()` to fully untangle the list in one operation.

```go
// After bulk insertions
for i := 0; i < 10000; i++ {
    list.PushBack(i)
}

// Optimize for reads
list.Crystalize()

// Now iterations have optimal cache locality
for value := range list.All() {
    process(value)
}
```

**Behavior:**
- Rearranges all nodes so memory order matches logical order
- Marks the list as "crystallized"
- Subsequent reads skip incremental untangling (already optimized)
- Any write operation automatically de-crystallizes the list

**When to use:**
- After bulk insertions, before a read-heavy phase
- Before serialization or benchmarking
- When you need predictable iteration performance

**Complexity:** O(n²) worst case for heavily fragmented lists. For mostly-ordered lists, approaches O(n).

### Memory Layout

The list stores nodes in an internal array managed by a buddy allocator. After many insertions and deletions, nodes may become scattered in memory. Incremental untangling gradually improves this layout during normal use.

For immediate cache-friendly access:
- Use `Crystalize()` to fully optimize the list in-place
- Use `ToSlice()` to get a contiguous copy
- Consider whether a slice or other data structure better fits your access pattern

### Tradeoffs

| Aspect | LinkedList | Slice |
|--------|------------|-------|
| Front insertion | O(1) | O(n) |
| Back insertion | O(1) | O(1) amortized |
| Random access | O(n) | O(1) |
| Cache locality | Poor after churn | Excellent |
| Memory overhead | Higher (node metadata) | Lower |

## Examples

### Stack (LIFO)

```go
stack := ll.NewLinkedList[int]()
stack.PushBack(1)
stack.PushBack(2)
stack.PushBack(3)

for !stack.IsEmpty() {
    value, _ := stack.PopBack()
    fmt.Println(value)  // 3, 2, 1
}
```

### Queue (FIFO)

```go
queue := ll.NewLinkedList[string]()
queue.PushBack("first")
queue.PushBack("second")
queue.PushBack("third")

for !queue.IsEmpty() {
    value, _ := queue.PopFront()
    fmt.Println(value)  // first, second, third
}
```

### Filtering

```go
list := ll.NewLinkedListFromSlice([]int{1, 2, 3, 4, 5, 6})

// Remove all even numbers
i := 0
for i < list.Length() {
    if value, _ := list.Get(i); value%2 == 0 {
        list.RemoveAt(i)
    } else {
        i++
    }
}
// list: [1, 3, 5]
```
