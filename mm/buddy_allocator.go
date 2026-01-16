package mm

import (
	"math/bits"

	"github.com/ConanHorus/fds/contracts"
	"github.com/ConanHorus/fds/internal/sorted"
)

// BuddyAllocator is a memory allocator that uses the buddy allocation algorithm
// to manage memory blocks efficiently.
//
// It divides memory into blocks of sizes that are powers of two and merges
// adjacent free blocks (buddies) to minimize fragmentation.
//
// The allocator supports dynamic growth of memory when needed, up to an
// optionally specified maximum order.
//
// Tradeoffs:
//   - Memory Efficiency vs. Speed: Buddy allocation is generally faster than
//     more complex allocation strategies, but it can lead to fragmentation
//     especially with varying allocation sizes. The merging of buddy blocks
//     helps mitigate this issue, but not completely.
type BuddyAllocator struct {
	state BuddyAllocatorOptions
	used  uint64

	freeLists [][]uint64 // various offsets of free blocks by order
}

// NewBuddyAllocator creates a new BuddyAllocator with the specified options.
//
// Tradeoffs:
//   - Memory Efficiency vs. Speed: Buddy allocation is generally faster than
//     more complex allocation strategies, but it can lead to fragmentation
//     especially with varying allocation sizes. The merging of buddy blocks
//     helps mitigate this issue, but not completely.
//
// Parameters:
//   - options: A variadic list of BuddyAllocatorOption functions to configure
//     the allocator.
//
// Returns:
//   - A pointer to the newly created BuddyAllocator.
func NewBuddyAllocator(options ...BuddyAllocatorOption) *BuddyAllocator {
	const (
		initialMaxOrder = 10
		initialSize     = 1 << 10 // default initial size of 1024 units
	)

	state := BuddyAllocatorOptions{
		Capacity:    initialSize,
		MaxOrder:    initialMaxOrder,
		MaxOrderSet: false,
	}

	for _, option := range options {
		option(&state)
	}

	if state.Capacity < 1<<state.MaxOrder {
		state.Capacity = 1 << state.MaxOrder
	}

	// Align capacity to be multiple of largest block size.
	state.Capacity = (state.Capacity + (1<<state.MaxOrder - 1)) & ^(1<<state.MaxOrder - 1)

	freeLists := make([][]uint64, state.MaxOrder+1)
	for index := uint64(0); index < state.Capacity; index += uint64(1) << state.MaxOrder {
		freeLists[state.MaxOrder] = append(freeLists[state.MaxOrder], index)
	}

	return &BuddyAllocator{
		state:     state,
		used:      0,
		freeLists: freeLists,
	}
}

// Allocate allocates a memory block of the specified size.
//
// If necessary, the allocator will grow its memory to accommodate the request,
// unless a maximum order has been set and the request exceeds it, in which case
// it will fail.
//
// Parameters:
//   - size: The size of the memory block to allocate.
//
// Returns:
//   - index: The starting index of the allocated memory block.
//   - grew: A boolean indicating whether the allocator had to grow its memory.
//   - ok: A boolean indicating whether the allocation was successful.
func (this *BuddyAllocator) Allocate(size uint64) (index uint64, grew bool, ok bool) {
	requiredOrder := orderOf(size)
	if requiredOrder > this.state.MaxOrder {
		if this.state.MaxOrderSet {
			return 0, false, false
		}

		for requiredOrder > this.state.MaxOrder {
			this.grow()
			grew = true
		}
	}

	for order := requiredOrder; order <= this.state.MaxOrder; order++ {
		if len(this.freeLists[order]) > 0 {
			index = this.freeLists[order][len(this.freeLists[order])-1]
			this.freeLists[order] = this.freeLists[order][:len(this.freeLists[order])-1]

			for currentOrder := order; currentOrder > requiredOrder; currentOrder-- {
				buddyIndex := index + (uint64(1) << (currentOrder - 1))
				this.freeLists[currentOrder-1] = sorted.InsertInt(this.freeLists[currentOrder-1], buddyIndex, false)
			}

			this.used += size
			return index, grew, true
		}
	}

	grew = true
	this.grow()
	index, _, ok = this.Allocate(size)
	return index, grew, ok
}

// Capacity returns the total capacity of the BuddyAllocator.
//
// Returns:
//   - The total capacity of the allocator.
func (this *BuddyAllocator) Capacity() uint64 {
	return this.state.Capacity
}

// ClearAll resets the BuddyAllocator, freeing all allocated memory blocks.
//
// After calling this method, all memory blocks are considered free, and the
// used memory count is reset to zero.
func (this *BuddyAllocator) ClearAll() {
	this.used = 0
	for order := uint8(0); order <= this.state.MaxOrder; order++ {
		this.freeLists[order] = this.freeLists[order][:0]
	}

	for index := uint64(0); index < this.state.Capacity; index += uint64(1) << this.state.MaxOrder {
		this.freeLists[this.state.MaxOrder] = append(this.freeLists[this.state.MaxOrder], index)
	}
}

// Efficiency returns the efficiency of the BuddyAllocator as a Percent.
//
// Efficiency is calculated as the ratio of used memory to total capacity.
//
// Returns:
//   - A Percent representing the efficiency of the allocator.
func (this *BuddyAllocator) Efficiency() contracts.Percent {
	return contracts.MakePercentFromDecimal(float64(this.used) / float64(this.state.Capacity))
}

// Free releases a previously allocated memory block back to the allocator.
//
// It attempts to merge the freed block with its buddy blocks to minimize
// fragmentation.
//
// Parameters:
//   - index: The starting index of the memory block to free.
//   - size: The size of the memory block to free.
//
// Returns:
//   - ok: A boolean indicating whether the free operation was successful. A
//     false signal on ok is an exceptional case indicating an invalid free
//     operation (e.g., double free).
func (this *BuddyAllocator) Free(index uint64, size uint64) (ok bool) {
	if index >= this.state.Capacity {
		return false
	}

	for currentOrder := uint8(0); currentOrder <= this.state.MaxOrder; currentOrder++ {
		if _, found := sorted.GallopingSearchInt(this.freeLists[currentOrder], index); found {
			return false
		}
	}

	requiredOrder := orderOf(size)
	this.used -= size

	currentIndex := index
	for currentOrder := requiredOrder; currentOrder < this.state.MaxOrder; currentOrder++ {
		if currentIndex%(uint64(1)<<(currentOrder+1)) != 0 {
			return false
		}

		buddyIndex := currentIndex ^ (uint64(1) << currentOrder)
		buddyFound := false

		if index, found := sorted.GallopingSearchInt(this.freeLists[currentOrder], buddyIndex); found {
			this.freeLists[currentOrder] = append(this.freeLists[currentOrder][:index], this.freeLists[currentOrder][index+1:]...)
			currentIndex = min(currentIndex, buddyIndex)
			buddyFound = true
		}

		if !buddyFound {
			this.freeLists[currentOrder] = sorted.InsertInt(this.freeLists[currentOrder], currentIndex, false)
			return true
		}
	}

	this.freeLists[this.state.MaxOrder] = sorted.InsertInt(this.freeLists[this.state.MaxOrder], currentIndex, false)
	return true
}

// Used returns the amount of memory currently used by the BuddyAllocator.
//
// Returns:
//   - The amount of used memory.
func (this *BuddyAllocator) Used() uint64 {
	return this.used
}

// --- private methods --- //

func (this *BuddyAllocator) grow() uint64 {
	newMemoryIndex := this.state.Capacity
	newMemorySize := max((uint64(1) << this.state.MaxOrder), 8)
	newMemoryOrder := orderOf(newMemorySize)
	newCapacity := this.state.Capacity + newMemorySize
	this.state.Capacity = newCapacity

	if !this.state.MaxOrderSet {
		this.state.MaxOrder = orderOf(newCapacity)
	}

	if uint8(len(this.freeLists)) <= this.state.MaxOrder {
		newFreeLists := make([][]uint64, this.state.MaxOrder+1)
		copy(newFreeLists, this.freeLists)
		this.freeLists = newFreeLists
	}

	this.freeLists[newMemoryOrder] = append(this.freeLists[newMemoryOrder], newMemoryIndex)
	return this.state.Capacity
}

// --- private functions --- //

func orderOf(size uint64) uint8 {
	if size == 0 {
		return 0
	}

	return uint8(bits.Len64(size) - 1)
}
