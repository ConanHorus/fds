package mm

const (
	minimumMaxOrder = 3      // minimum max order of 8 units
	minimumCapacity = 1 << 3 // minimum capacity of 8 units
)

// BuddyAllocatorOptions holds configuration options for BuddyAllocator.
type BuddyAllocatorOptions struct {
	Capacity uint64

	// MaxOrder defines which multiple of 2 is being used to divide memory
	// blocks.
	//
	// For example, if MaxOrder is set to 10, the largest block size will be
	// 2^10 = 1024 units.
	//
	// If MaxOrderSet is false, the BuddyAllocator will determine the maximum
	// order every time memory growth is needed.
	MaxOrder    uint8
	MaxOrderSet bool
}

// BuddyAllocatorOption defines a function type for setting BuddyAllocator
// options.
type BuddyAllocatorOption func(*BuddyAllocatorOptions)

// WithInitialCapacity sets the initial capacity option for BuddyAllocator.
//
// Parameters:
//   - capacity: The initial capacity to set.
//
// Returns:
//   - A BuddyAllocatorOption function that sets the Capacity field.
func WithInitialCapacity(capacity uint64) BuddyAllocatorOption {
	return func(opts *BuddyAllocatorOptions) {
		opts.Capacity = max(capacity, minimumCapacity)
	}
}

// WithMaxOrder sets the maximum order option for BuddyAllocator.
//
// MaxOrder defines which multiple of 2 is being used to divide memory. For
// example, if MaxOrder is set to 10, the largest block size will be 2^10 = 1024
// units.
//
// Parameters:
//   - order: The maximum order to set.
//
// Returns:
//   - A BuddyAllocatorOption function that sets the MaxOrder field and marks
//     MaxOrderSet as true.
func WithMaxOrder(order uint8) BuddyAllocatorOption {
	return func(opts *BuddyAllocatorOptions) {
		opts.MaxOrder = max(order, minimumMaxOrder)
		opts.MaxOrderSet = true
	}
}
