package mm

import (
	"testing"

	"github.com/smarty/assertions"
	"github.com/smarty/assertions/should"
	"github.com/smarty/benchy"
	"github.com/smarty/benchy/options"
	"github.com/smarty/benchy/providers"
)

func TestBuddyAllocator(t *testing.T) {
	and := assertions.New(t)
	allocator := NewBuddyAllocator(
		WithInitialCapacity(minimumCapacity),
		WithMaxOrder(minimumMaxOrder),
	)

	and.So(allocator, should.NotBeNil)
	and.So(allocator.state.Capacity, should.Equal, minimumCapacity)
	and.So(allocator.state.MaxOrder, should.Equal, minimumMaxOrder)

	for range minimumCapacity {
		_, grew, ok := allocator.Allocate(1)
		and.So(grew, should.BeFalse)
		and.So(ok, should.BeTrue)
	}

	_, grew, ok := allocator.Allocate(1)
	and.So(grew, should.BeTrue)
	and.So(ok, should.BeTrue)
	and.So(allocator.Used(), should.Equal, minimumCapacity+1)

	for i := range minimumCapacity + 1 {
		ok := allocator.Free(uint64(i), 1)
		and.So(ok, should.BeTrue)
	}

	and.So(allocator.Used(), should.Equal, 0)
}

func TestBuddyAllocator_TooMuchAllocation(t *testing.T) {
	and := assertions.New(t)
	allocator := NewBuddyAllocator(
		WithInitialCapacity(minimumCapacity),
		WithMaxOrder(minimumMaxOrder),
	)

	_, _, ok := allocator.Allocate(1 << (minimumMaxOrder + 1))
	and.So(ok, should.BeFalse)
}

func TestBuddyAllocator_DoubleFree(t *testing.T) {
	and := assertions.New(t)
	allocator := NewBuddyAllocator(
		WithInitialCapacity(minimumCapacity),
		WithMaxOrder(minimumMaxOrder),
	)

	index, _, ok := allocator.Allocate(1)
	and.So(ok, should.BeTrue)

	ok = allocator.Free(index, 1)
	and.So(ok, should.BeTrue)

	ok = allocator.Free(index, 1)
	and.So(ok, should.BeFalse)
}

func TestBuddyAllocator_FreeMisaligned(t *testing.T) {
	and := assertions.New(t)
	allocator := NewBuddyAllocator(
		WithInitialCapacity(minimumCapacity),
		WithMaxOrder(minimumMaxOrder),
	)

	index, _, ok := allocator.Allocate(2)
	and.So(ok, should.BeTrue)

	ok = allocator.Free(index+1, 1)
	and.So(ok, should.BeFalse)
}

func TestBuddyAllocator_FreeIndexOutOfBounds(t *testing.T) {
	and := assertions.New(t)
	allocator := NewBuddyAllocator(
		WithInitialCapacity(minimumCapacity),
		WithMaxOrder(minimumMaxOrder),
	)

	ok := allocator.Free(minimumCapacity+1, 1)
	and.So(ok, should.BeFalse)
}

// --- benchmarks --- //

func BenchmarkBuddyAllocator_Allocate(b *testing.B) {
	allocator := NewBuddyAllocator()

	provider := providers.New1(func(int) {}).
		Add(1).
		Add(2).
		Add(4).
		Add(2).
		Add(1).
		Add(1).
		Add(2).
		Add(2).
		Add(4).
		Add(8)

	benchy.New(b, options.Medium).
		RegisterBenchmark("Allocate", provider.WrapBenchmarkFunc(func(int) {
			allocator.ClearAll()
			_, _, ok := allocator.Allocate(64)
			if !ok {
				b.Errorf("Allocation failed")
			}
		})).
		Run()
}

func BenchmarkBuddyAllocator_Free(b *testing.B) {
	allocator := NewBuddyAllocator(
		WithInitialCapacity(1<<12),
		WithMaxOrder(10),
	)

	index := uint64(0)
	benchy.New(b, options.Medium).
		RegisterBenchmark("Free", func() {
			if allocator.Used() == 0 {
				index = 0
				allocator.ClearAll()
				allocator.freeLists[len(allocator.freeLists)-1] = allocator.freeLists[len(allocator.freeLists)-1][:0]
				allocator.used = allocator.state.Capacity
			}

			allocator.Free(index, 2)
			index += 2
		}).
		Run()
}
