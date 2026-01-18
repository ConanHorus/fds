# Memory Manager Package (`mm`)

The `mm` package provides high-performance memory allocation primitives for use by other data structures in the FDS library.

## Overview

`BuddyAllocator` is the primary allocator, implementing the buddy allocation algorithm. It manages memory blocks efficiently by:

- Dividing memory into power-of-two sized blocks
- Merging adjacent free blocks ("buddies") to reduce fragmentation
- Growing automatically when more memory is needed

## Quick Start

```go
import "github.com/ConanHorus/fds/mm"

// Create with default settings
allocator := mm.NewBuddyAllocator()

// Create with options
allocator := mm.NewBuddyAllocator(
    mm.WithInitialCapacity(1024),
    mm.WithMaxOrder(10),
)

// Allocate memory
index, grew, ok := allocator.Allocate(16)
if ok {
    // Use indices [index, index+16)
}

// Free memory
ok := allocator.Free(index, 16)

// Query state
capacity := allocator.Capacity()
used := allocator.Used()
efficiency := allocator.Efficiency()
```

## BuddyAllocator

### Constructor

```go
func NewBuddyAllocator(options ...BuddyAllocatorOption) *BuddyAllocator
```

Creates a new allocator. Options can customize initial capacity and maximum block order.

### Options

| Option | Description |
|--------|-------------|
| `WithInitialCapacity(capacity uint64)` | Sets the initial memory capacity (minimum 8) |
| `WithMaxOrder(order uint8)` | Sets the maximum block order (minimum 3, meaning max block size of 8) |

### Methods

| Method | Description |
|--------|-------------|
| `Allocate(size uint64) (index, grew, ok)` | Allocates a block of at least `size` units |
| `Free(index, size uint64) bool` | Frees a previously allocated block |
| `Capacity() uint64` | Returns total capacity |
| `Used() uint64` | Returns amount of memory in use |
| `Efficiency() Percent` | Returns used/capacity ratio |
| `ClearAll()` | Resets all allocations |

## API Reference

### Allocate

```go
func (this *BuddyAllocator) Allocate(size uint64) (index uint64, grew bool, ok bool)
```

Allocates a memory block of at least the specified size.

**Parameters:**
- `size`: The size to allocate (must be > 0)

**Returns:**
- `index`: Starting index of the allocated block
- `grew`: Whether the allocator had to grow its capacity
- `ok`: Whether allocation succeeded

**Behavior:**
- Size is rounded up to the nearest power of two
- If no suitable block exists, the allocator grows automatically
- If `WithMaxOrder` was set and growth would exceed it, allocation fails

**Examples:**
- `Allocate(1)` → allocates 1 unit (2⁰)
- `Allocate(5)` → allocates 8 units (2³)
- `Allocate(0)` → returns `ok = false`

### Free

```go
func (this *BuddyAllocator) Free(index uint64, size uint64) (ok bool)
```

Frees a previously allocated block.

**Parameters:**
- `index`: Starting index of the block to free
- `size`: Size of the block (must match the allocation)

**Returns:**
- `ok`: Whether the free succeeded (false indicates double-free or invalid parameters)

**Behavior:**
- Automatically merges with buddy blocks when possible
- Returns false for invalid operations (double-free, overlapping regions)

### ClearAll

```go
func (this *BuddyAllocator) ClearAll()
```

Resets the allocator, marking all memory as free. This is more efficient than freeing individual blocks.

## Buddy Allocation Algorithm

The buddy system divides memory into blocks whose sizes are powers of two. When allocating:

1. Find the smallest power-of-two block that fits the request
2. If no block of that size exists, split a larger block
3. Continue splitting until a suitable block is available

When freeing:

1. Check if the block's "buddy" (adjacent block of same size) is also free
2. If so, merge them into a larger block
3. Repeat until no more merging is possible

### Block Orders

The "order" of a block is its power of two:
- Order 0 = 2⁰ = 1 unit
- Order 3 = 2³ = 8 units
- Order 10 = 2¹⁰ = 1024 units

`WithMaxOrder(n)` limits the maximum block size to 2ⁿ units.

## Performance Characteristics

| Operation | Complexity |
|-----------|------------|
| Allocate | O(log n) average, O(n) worst case when growing |
| Free | O(log n) for buddy merging |
| Capacity | O(1) |
| Used | O(1) |
| ClearAll | O(max_order) |

## Tradeoffs

**Advantages:**
- Fast allocation and deallocation
- Automatic coalescing reduces fragmentation
- Predictable block sizes

**Disadvantages:**
- Internal fragmentation (e.g., allocating 5 units uses 8)
- External fragmentation can still occur with varied allocation patterns
- Memory overhead for tracking free lists

## Usage in FDS

The `BuddyAllocator` is used internally by:
- `LinkedList` — manages node storage
- (Future data structures)

Users of these data structures don't interact with the allocator directly; it's an implementation detail hidden behind clean APIs.

## Example: Custom Data Structure

```go
type MyContainer[T any] struct {
    data      []T
    allocator mm.BuddyAllocator
}

func NewMyContainer[T any]() *MyContainer[T] {
    alloc := mm.NewBuddyAllocator(mm.WithInitialCapacity(64))
    return &MyContainer[T]{
        data:      make([]T, alloc.Capacity()),
        allocator: *alloc,
    }
}

func (c *MyContainer[T]) Add(value T) int {
    index, grew, ok := c.allocator.Allocate(1)
    if !ok {
        panic("allocation failed")
    }
    if grew {
        newData := make([]T, c.allocator.Capacity())
        copy(newData, c.data)
        c.data = newData
    }
    c.data[index] = value
    return int(index)
}
```
