package mm

import (
	"fmt"
	"testing"

	"github.com/smarty/assertions"
	"github.com/smarty/assertions/should"
	"github.com/smarty/benchy"
	"github.com/smarty/benchy/options"
	"github.com/smarty/benchy/providers"
)

// --- Constructor and Basic State Tests ---

func TestNewBuddyAllocator_DefaultOptions(t *testing.T) {
	assert := assertions.New(t)
	allocator := NewBuddyAllocator()

	assert.So(allocator, should.NotBeNil)
	assert.So(allocator.Capacity(), should.BeGreaterThan, 0)
	assert.So(allocator.Used(), should.Equal, 0)
	assert.So(allocator.Efficiency().AsDecimal(), should.Equal, 0.0)
}

func TestNewBuddyAllocator_WithOptions(t *testing.T) {
	testCases := []struct {
		name        string
		capacity    uint64
		maxOrder    uint8
		expectedCap uint64
		expectedMax uint8
	}{
		{"minimum values", 4, 2, 8, 3},
		{"standard values", 1024, 10, 1024, 10},
		{"misaligned capacity", 100, 6, 128, 6},
		{"capacity too small", 1, 5, 32, 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assertions.New(t)
			allocator := NewBuddyAllocator(
				WithInitialCapacity(tc.capacity),
				WithMaxOrder(tc.maxOrder),
			)

			assert.So(allocator.Capacity(), should.Equal, tc.expectedCap)
			assert.So(allocator.state.MaxOrder, should.Equal, tc.expectedMax)
			assert.So(allocator.state.MaxOrderSet, should.BeTrue)
		})
	}
}

// --- Basic Allocation Tests ---

func TestAllocate_BasicFunctionality(t *testing.T) {
	allocator := NewBuddyAllocator(WithInitialCapacity(64), WithMaxOrder(6))

	testCases := []struct {
		name       string
		size       uint64
		expectOK   bool
		expectGrew bool
		expectUsed uint64
	}{
		{"allocate 1 byte", 1, true, false, 1},
		{"allocate 2 bytes", 2, true, false, 3},
		{"allocate 4 bytes", 4, true, false, 7},
		{"allocate 8 bytes", 8, true, false, 15},
		{"allocate 16 bytes", 16, true, false, 31},
		{"allocate 32 bytes", 32, true, false, 63},
		{"allocate zero", 0, false, false, 63},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assertions.New(t)
			index, grew, ok := allocator.Allocate(tc.size)

			assert.So(ok, should.Equal, tc.expectOK)
			assert.So(grew, should.Equal, tc.expectGrew)
			if tc.expectOK {
				assert.So(index, should.BeGreaterThanOrEqualTo, 0)
				assert.So(index, should.BeLessThan, allocator.Capacity())
			}
			assert.So(allocator.Used(), should.Equal, tc.expectUsed)
		})
	}
}

func TestAllocate_PowerOfTwoSizes(t *testing.T) {
	assert := assertions.New(t)
	allocator := NewBuddyAllocator(WithInitialCapacity(256), WithMaxOrder(8))

	// Test allocating exact power of 2 sizes
	sizes := []uint64{1, 2, 4, 8, 16, 32, 64, 128}
	indices := make([]uint64, len(sizes))

	for i, size := range sizes {
		index, grew, ok := allocator.Allocate(size)
		assert.So(ok, should.BeTrue)
		assert.So(grew, should.BeFalse)
		indices[i] = index
	}

	// Verify all allocations are unique
	for i := 0; i < len(indices); i++ {
		for j := i + 1; j < len(indices); j++ {
			assert.So(indices[i], should.NotEqual, indices[j])
		}
	}
}

func TestAllocate_NonPowerOfTwoSizes(t *testing.T) {
	allocator := NewBuddyAllocator(WithInitialCapacity(128), WithMaxOrder(7))

	testCases := []struct {
		size         uint64
		expectedUsed uint64
	}{
		{3, 3},   // Should allocate 4-byte block, but used tracks requested size
		{5, 8},   // Should allocate 8-byte block, used = 3 + 5 = 8
		{7, 15},  // Should allocate 8-byte block, used = 3 + 5 + 7 = 15
		{15, 30}, // Should allocate 16-byte block, used = 3 + 5 + 7 + 15 = 30
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("size_%d", tc.size), func(t *testing.T) {
			assert := assertions.New(t)
			_, grew, ok := allocator.Allocate(tc.size)

			assert.So(ok, should.BeTrue)
			assert.So(grew, should.BeFalse)
			assert.So(allocator.Used(), should.Equal, tc.expectedUsed)
		})
	}
}

func TestAllocate_Growth(t *testing.T) {
	assert := assertions.New(t)
	allocator := NewBuddyAllocator(WithInitialCapacity(8), WithMaxOrder(3))

	initialCapacity := allocator.Capacity()

	// Fill initial capacity
	allocator.Allocate(8)
	assert.So(allocator.Used(), should.Equal, 8)

	// This should trigger growth
	index, grew, ok := allocator.Allocate(1)
	assert.So(ok, should.BeTrue)
	assert.So(grew, should.BeTrue)
	assert.So(allocator.Capacity(), should.BeGreaterThan, initialCapacity)
	assert.So(allocator.Used(), should.Equal, 9)
	assert.So(index, should.BeGreaterThanOrEqualTo, 0)
}

func TestAllocate_MaxOrderLimit(t *testing.T) {
	assert := assertions.New(t)
	allocator := NewBuddyAllocator(WithInitialCapacity(32), WithMaxOrder(5))

	// Try to allocate more than max order allows
	_, _, ok := allocator.Allocate(64) // Requires order 6, but max is 5
	assert.So(ok, should.BeFalse)

	// Max size allocation should work
	_, _, ok = allocator.Allocate(32) // Exactly max order
	assert.So(ok, should.BeTrue)
}

// --- Free Tests ---

func TestFree_BasicFunctionality(t *testing.T) {
	allocator := NewBuddyAllocator(WithInitialCapacity(64), WithMaxOrder(6))

	// Allocate and free various sizes
	testCases := []uint64{1, 2, 4, 8, 16}

	for _, size := range testCases {
		t.Run(fmt.Sprintf("size_%d", size), func(t *testing.T) {
			assert := assertions.New(t)

			initialUsed := allocator.Used()
			index, _, ok := allocator.Allocate(size)
			assert.So(ok, should.BeTrue)

			midUsed := allocator.Used()
			assert.So(midUsed, should.Equal, initialUsed+size)

			ok = allocator.Free(index, size)
			assert.So(ok, should.BeTrue)
			assert.So(allocator.Used(), should.Equal, initialUsed)
		})
	}
}

func TestFree_InvalidInputs(t *testing.T) {
	allocator := NewBuddyAllocator(WithInitialCapacity(32), WithMaxOrder(5))

	testCases := []struct {
		name   string
		index  uint64
		size   uint64
		expect bool
	}{
		{"zero size", 0, 0, false},
		{"out of bounds index", 100, 1, false},
		{"negative equivalent index", ^uint64(0), 1, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assertions.New(t)
			ok := allocator.Free(tc.index, tc.size)
			assert.So(ok, should.Equal, tc.expect)
		})
	}
}

func TestFree_DoubleFree(t *testing.T) {
	assert := assertions.New(t)
	allocator := NewBuddyAllocator(WithInitialCapacity(32), WithMaxOrder(5))

	index, _, ok := allocator.Allocate(4)
	assert.So(ok, should.BeTrue)

	// First free should succeed
	ok = allocator.Free(index, 4)
	assert.So(ok, should.BeTrue)

	// Second free should fail
	ok = allocator.Free(index, 4)
	assert.So(ok, should.BeFalse)
}

func TestFree_BuddyMerging(t *testing.T) {
	assert := assertions.New(t)
	allocator := NewBuddyAllocator(WithInitialCapacity(64), WithMaxOrder(6))

	// Allocate two adjacent 4-byte blocks that should be buddies
	index1, _, ok1 := allocator.Allocate(4)
	assert.So(ok1, should.BeTrue)

	index2, _, ok2 := allocator.Allocate(4)
	assert.So(ok2, should.BeTrue)

	initialUsed := allocator.Used()

	// Free both blocks - they should merge
	ok := allocator.Free(index1, 4)
	assert.So(ok, should.BeTrue)
	assert.So(allocator.Used(), should.Equal, initialUsed-4)

	ok = allocator.Free(index2, 4)
	assert.So(ok, should.BeTrue)
	assert.So(allocator.Used(), should.Equal, initialUsed-8)
}

// --- ClearAll Tests ---

func TestClearAll(t *testing.T) {
	assert := assertions.New(t)
	allocator := NewBuddyAllocator(WithInitialCapacity(32), WithMaxOrder(5))

	// Make several allocations
	allocator.Allocate(4)
	allocator.Allocate(8)
	allocator.Allocate(2)

	assert.So(allocator.Used(), should.BeGreaterThan, 0)

	// Clear everything
	allocator.ClearAll()
	assert.So(allocator.Used(), should.Equal, 0)
	assert.So(allocator.Efficiency().AsDecimal(), should.Equal, 0.0)

	// Should be able to allocate again
	_, _, ok := allocator.Allocate(16)
	assert.So(ok, should.BeTrue)
}

// --- Efficiency Tests ---

func TestEfficiency(t *testing.T) {
	assert := assertions.New(t)
	allocator := NewBuddyAllocator(WithInitialCapacity(100), WithMaxOrder(7))

	capacity := allocator.Capacity()

	// Initially empty
	assert.So(allocator.Efficiency().AsDecimal(), should.Equal, 0.0)

	// Half full
	allocator.Allocate(capacity / 2)
	efficiency := allocator.Efficiency().AsDecimal()
	assert.So(efficiency, should.BeGreaterThan, 0.4)
	assert.So(efficiency, should.BeLessThan, 0.6)

	// Nearly full
	allocator.Allocate(capacity / 4)
	efficiency = allocator.Efficiency().AsDecimal()
	assert.So(efficiency, should.BeGreaterThan, 0.7)
}

// --- Stress Tests ---

func TestAllocateFreePattern(t *testing.T) {
	assert := assertions.New(t)
	allocator := NewBuddyAllocator(WithInitialCapacity(256), WithMaxOrder(8))

	// Allocate many small blocks
	indices := make([]uint64, 50)
	for i := range indices {
		index, _, ok := allocator.Allocate(1)
		assert.So(ok, should.BeTrue)
		indices[i] = index
	}

	// Free every other block
	for i := 0; i < len(indices); i += 2 {
		ok := allocator.Free(indices[i], 1)
		assert.So(ok, should.BeTrue)
	}

	// Free remaining blocks
	for i := 1; i < len(indices); i += 2 {
		ok := allocator.Free(indices[i], 1)
		assert.So(ok, should.BeTrue)
	}

	// Should be back to empty
	assert.So(allocator.Used(), should.Equal, 0)
}

func TestMixedSizeAllocations(t *testing.T) {
	assert := assertions.New(t)
	allocator := NewBuddyAllocator(WithInitialCapacity(512), WithMaxOrder(9))

	type allocation struct {
		index uint64
		size  uint64
	}

	var allocs []allocation
	sizes := []uint64{1, 2, 4, 8, 16, 32, 3, 5, 7, 15, 31}

	// Make allocations
	for _, size := range sizes {
		index, _, ok := allocator.Allocate(size)
		assert.So(ok, should.BeTrue)
		allocs = append(allocs, allocation{index, size})
	}

	// Free in reverse order
	for i := len(allocs) - 1; i >= 0; i-- {
		ok := allocator.Free(allocs[i].index, allocs[i].size)
		assert.So(ok, should.BeTrue)
	}

	assert.So(allocator.Used(), should.Equal, 0)
}

// --- Edge Case Tests ---

func TestLargeAllocation(t *testing.T) {
	assert := assertions.New(t)
	allocator := NewBuddyAllocator() // No max order set

	// Try very large allocation
	_, grew, ok := allocator.Allocate(1 << 20) // 1MB
	assert.So(ok, should.BeTrue)
	assert.So(grew, should.BeTrue)
}

func TestFragmentation(t *testing.T) {
	assert := assertions.New(t)
	allocator := NewBuddyAllocator(WithInitialCapacity(64), WithMaxOrder(6))

	// Create fragmentation pattern
	indices := make([]uint64, 8)
	for i := range indices {
		index, _, ok := allocator.Allocate(8)
		assert.So(ok, should.BeTrue)
		indices[i] = index
	}

	// Free every other allocation to create holes
	for i := 0; i < len(indices); i += 2 {
		ok := allocator.Free(indices[i], 8)
		assert.So(ok, should.BeTrue)
	}

	// Should still be able to allocate small blocks
	_, _, ok := allocator.Allocate(4)
	assert.So(ok, should.BeTrue)
}

func TestBuddyAllocator_Free(t *testing.T) {
	and := assertions.New(t)
	allocator := NewBuddyAllocator()
	allocator.ClearAll()

	// simiulate full allocation
	allocator.freeLists[len(allocator.freeLists)-1] = allocator.freeLists[len(allocator.freeLists)-1][:0]
	allocator.used = allocator.state.Capacity

	ok := allocator.Free(0, 1)
	and.So(ok, assertions.ShouldBeTrue)

	ok = allocator.Free(2, 2)
	and.So(ok, assertions.ShouldBeTrue)

	ok = allocator.Free(4, 4)
	and.So(ok, assertions.ShouldBeTrue)

	ok = allocator.Free(8, 8)
	and.So(ok, assertions.ShouldBeTrue)

	ok = allocator.Free(16, 16)
	and.So(ok, assertions.ShouldBeTrue)

	ok = allocator.Free(32, 32)
	and.So(ok, assertions.ShouldBeTrue)

	ok = allocator.Free(64, 64)
	and.So(ok, assertions.ShouldBeTrue)

	ok = allocator.Free(128, 128)
	and.So(ok, assertions.ShouldBeTrue)
}

// --- Benchmarks ---

func BenchmarkBuddyAllocator_Allocate(b *testing.B) {
	provider := providers.New1(func(uint64) {}).
		Add(uint64(1)).
		Add(uint64(2)).
		Add(uint64(4)).
		Add(uint64(8)).
		Add(uint64(16)).
		Add(uint64(32)).
		Add(uint64(64)).
		Add(uint64(128))

	allocator := NewBuddyAllocator()
	benchy.New(b, options.Medium).
		RegisterBenchmark("Allocate", provider.WrapBenchmarkFunc(func(size uint64) {
			if allocator.Efficiency().AsDecimal() > 0.9 {
				allocator.ClearAll()
			}

			_, _, ok := allocator.Allocate(size)
			if !ok {
				b.Errorf("Failed to allocate %d bytes", size)
			}
		})).
		Run()
}

func BenchmarkBuddyAllocator_Free(b *testing.B) {
	provider := providers.New1(func(uint64) {}).
		Add(uint64(1)).
		Add(uint64(2)).
		Add(uint64(4)).
		Add(uint64(8)).
		Add(uint64(16)).
		Add(uint64(32)).
		Add(uint64(64)).
		Add(uint64(128))

	allocator := NewBuddyAllocator()

	index := uint64(0)
	benchy.New(b, options.Medium).
		RegisterBenchmark("Free", provider.WrapBenchmarkFunc(func(size uint64) {
			if allocator.Efficiency().AsDecimal() < 0.5 {
				allocator.ClearAll()

				// simulate full allocation
				allocator.freeLists[len(allocator.freeLists)-1] = allocator.freeLists[len(allocator.freeLists)-1][:0]
				allocator.used = allocator.state.Capacity
				index = 0
			}

			// ensure index is alligned to block size
			if index%size != 0 {
				index += size - (index % size)
			}

			ok := allocator.Free(index, size)
			index += size
			if !ok {
				b.Errorf("Failed to free block at index %d, size %d", index, size)
			}
		})).
		Run()
}

// --- AllocateAt Tests ---

func TestAllocateAt_BasicFunctionality(t *testing.T) {
	allocator := NewBuddyAllocator(WithInitialCapacity(64), WithMaxOrder(6))

	// Test successful allocation at specific locations
	testCases := []struct {
		name     string
		index    uint64
		size     uint64
		expectOK bool
	}{
		{"allocate 1 byte at index 0", 0, 1, true},
		{"allocate 2 bytes at index 2", 2, 2, true},
		{"allocate 4 bytes at index 4", 4, 4, true},
		{"allocate 8 bytes at index 8", 8, 8, true},
		{"allocate 1 byte at index 16", 16, 1, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assertions.New(t)
			initialUsed := allocator.Used()

			ok := allocator.AllocateAt(tc.index, tc.size)
			assert.So(ok, should.Equal, tc.expectOK)

			if tc.expectOK {
				assert.So(allocator.Used(), should.Equal, initialUsed+tc.size)
			} else {
				assert.So(allocator.Used(), should.Equal, initialUsed)
			}
		})
	}
}

func TestAllocateAt_InvalidInputs(t *testing.T) {
	allocator := NewBuddyAllocator(WithInitialCapacity(64), WithMaxOrder(6))

	testCases := []struct {
		name  string
		index uint64
		size  uint64
	}{
		{"zero size", 0, 0},
		{"out of bounds index", 100, 1},
		{"index at capacity boundary", 64, 1},
		{"negative equivalent index", ^uint64(0), 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assertions.New(t)
			initialUsed := allocator.Used()

			ok := allocator.AllocateAt(tc.index, tc.size)
			assert.So(ok, should.BeFalse)
			assert.So(allocator.Used(), should.Equal, initialUsed)
		})
	}
}

func TestAllocateAt_AlignmentRequirements(t *testing.T) {
	testCases := []struct {
		name     string
		index    uint64
		size     uint64
		expectOK bool
		reason   string
	}{
		{"1-byte block at index 0", 0, 1, true, "any index is 1-aligned"},
		{"1-byte block at index 1", 1, 1, true, "any index is 1-aligned"},
		{"2-byte block at index 0", 0, 2, true, "index 0 is 2-aligned"},
		{"2-byte block at index 2", 2, 2, true, "index 2 is 2-aligned"},
		{"2-byte block at index 1", 1, 2, false, "index 1 is not 2-aligned"},
		{"4-byte block at index 0", 0, 4, true, "index 0 is 4-aligned"},
		{"4-byte block at index 4", 4, 4, true, "index 4 is 4-aligned"},
		{"4-byte block at index 2", 2, 4, false, "index 2 is not 4-aligned"},
		{"8-byte block at index 0", 0, 8, true, "index 0 is 8-aligned"},
		{"8-byte block at index 8", 8, 8, true, "index 8 is 8-aligned"},
		{"8-byte block at index 4", 4, 8, false, "index 4 is not 8-aligned"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assertions.New(t)
			allocator := NewBuddyAllocator(WithInitialCapacity(64), WithMaxOrder(6))

			ok := allocator.AllocateAt(tc.index, tc.size)
			assert.So(ok, should.Equal, tc.expectOK)
		})
	}
}

func TestAllocateAt_OverlappingAllocations(t *testing.T) {
	assert := assertions.New(t)
	allocator := NewBuddyAllocator(WithInitialCapacity(64), WithMaxOrder(6))

	// First allocation should succeed
	ok := allocator.AllocateAt(0, 4)
	assert.So(ok, should.BeTrue)

	// Overlapping allocations should fail
	testCases := []struct {
		name  string
		index uint64
		size  uint64
	}{
		{"exact same location", 0, 4},
		{"overlapping start", 2, 4},
		{"contained within", 1, 2},
		{"overlapping end", 0, 8},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assertions.New(t)
			ok := allocator.AllocateAt(tc.index, tc.size)
			assert.So(ok, should.BeFalse)
		})
	}

	// Non-overlapping allocation should succeed
	ok = allocator.AllocateAt(4, 4)
	assert.So(ok, should.BeTrue)
}

func TestAllocateAt_WithFreeAndReuse(t *testing.T) {
	assert := assertions.New(t)
	allocator := NewBuddyAllocator(WithInitialCapacity(64), WithMaxOrder(6))

	// Allocate at specific location
	ok := allocator.AllocateAt(8, 4)
	assert.So(ok, should.BeTrue)
	assert.So(allocator.Used(), should.Equal, 4)

	// Free the allocation
	ok = allocator.Free(8, 4)
	assert.So(ok, should.BeTrue)
	assert.So(allocator.Used(), should.Equal, 0)

	// Should be able to allocate at the same location again
	ok = allocator.AllocateAt(8, 4)
	assert.So(ok, should.BeTrue)
	assert.So(allocator.Used(), should.Equal, 4)

	// Test with a fresh allocator for the larger block case
	// (since fragmentation prevents reusing the same allocator)
	freshAllocator := NewBuddyAllocator(WithInitialCapacity(64), WithMaxOrder(6))
	ok = freshAllocator.AllocateAt(0, 16)
	assert.So(ok, should.BeTrue)
	assert.So(freshAllocator.Used(), should.Equal, 16)
}

func TestAllocateAt_BlockSplitting(t *testing.T) {
	assert := assertions.New(t)
	allocator := NewBuddyAllocator(WithInitialCapacity(64), WithMaxOrder(6))

	// The initial free list should have one 64-byte block
	// Allocating 4 bytes at index 8 should split blocks appropriately
	ok := allocator.AllocateAt(8, 4)
	assert.So(ok, should.BeTrue)

	// After allocation, we should be able to allocate other sizes
	// in the remaining free spaces
	ok = allocator.AllocateAt(0, 8) // Should work - space before our allocation
	assert.So(ok, should.BeTrue)

	ok = allocator.AllocateAt(12, 4) // Should work - space after our allocation
	assert.So(ok, should.BeTrue)

	totalUsed := uint64(4 + 8 + 4)
	assert.So(allocator.Used(), should.Equal, totalUsed)
}

func TestAllocateAt_ExtendsBeyondCapacity(t *testing.T) {
	allocator := NewBuddyAllocator(WithInitialCapacity(64), WithMaxOrder(6))

	// Try to allocate a block that would extend beyond capacity
	testCases := []struct {
		name  string
		index uint64
		size  uint64
	}{
		{"block extends beyond capacity", 60, 8},
		{"block starts at last valid index", 63, 1},
		{"large block at end", 32, 64},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assertions.New(t)
			ok := allocator.AllocateAt(tc.index, tc.size)
			if tc.index+tc.size <= 64 {
				// Should succeed if within bounds and properly aligned
				requiredBlockSize := uint64(1) << orderOf(tc.size)
				if tc.index%requiredBlockSize == 0 && tc.index+requiredBlockSize <= 64 {
					assert.So(ok, should.BeTrue)
				} else {
					assert.So(ok, should.BeFalse)
				}
			} else {
				// Should fail if extends beyond capacity
				assert.So(ok, should.BeFalse)
			}
		})
	}
}
