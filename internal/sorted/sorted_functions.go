package sorted

// Sortable is a constraint that permits any integer type.
type Sortable interface {
	~uint8 | ~int8 | ~uint16 | ~int16 | ~uint32 | ~int32 | ~uint64 | ~int64 | ~uint | ~int | ~uintptr
}

// BinarySearchInt performs a binary search on a sorted slice of sortable
// elements to find the target value.
//
// Parameters:
//   - slice: A sorted slice of elements of type T.
//   - target: The value to search for within the slice.
//
// Returns:
//   - index: The index of the target value if found; otherwise, the index where
//     the target can be inserted to maintain sorted order.
//   - found: A boolean indicating whether the target value was found in the
//     slice.
func BinarySearchInt[T Sortable](slice []T, target T) (index int, found bool) {
	low := 0
	high := len(slice)
	for low < high {
		mid := (low + high) / 2
		if slice[mid] < target {
			low = mid + 1
			continue
		}

		high = mid
	}

	if low < len(slice) && slice[low] == target {
		for low > 0 && slice[low-1] == target {
			low--
		}

		return low, true
	}

	return low, false
}

// GallopingSearchInt performs a galloping search on a sorted slice of sortable
// elements to find the target value.
//
// For small slices (length less than a defined threshold), it falls back to
// binary search for efficiency.
//
// Parameters:
//   - slice: A sorted slice of elements of type T.
//   - target: The value to search for within the slice.
//
// Returns:
//   - index: The index of the target value if found; otherwise, the index where
//     the target can be inserted to maintain sorted order.
//   - found: A boolean indicating whether the target value was found in the
//     slice.
func GallopingSearchInt[T Sortable](slice []T, target T) (index int, found bool) {
	const threshold = 64
	if len(slice) < threshold {
		return BinarySearchInt(slice, target)
	}

	_ = slice[len(slice)-1] // touch to avoid bounds checks

	spanSize := threshold
	prevIndex := 0
	for {
		currentIndex := prevIndex + spanSize
		if currentIndex < 0 || currentIndex >= len(slice) {
			currentIndex = len(slice) - 1
			return BinarySearchInt(slice[prevIndex:currentIndex+1], target)
		}

		if slice[currentIndex] < target {
			prevIndex = currentIndex
			spanSize = spanSize << 1
			continue
		}

		return BinarySearchInt(slice[prevIndex:currentIndex+1], target)
	}
}

// InsertInt inserts a value into a sorted slice of sortable elements while
// maintaining sorted order.
//
// If allowDuplicates is false and the value already exists in the slice, the
// slice is returned unchanged.
//
// Parameters:
//   - slice: A sorted slice of elements of type T.
//   - value: The value to insert into the slice.
//   - allowDuplicates: A boolean indicating whether duplicate values are
//     allowed.
//
// Returns:
//   - A new slice with the value inserted in sorted order.
func InsertInt[T Sortable](slice []T, value T, allowDuplicates bool) []T {
	index, found := GallopingSearchInt(slice, value)
	if found && !allowDuplicates {
		return slice
	}

	slice = append(slice, value) // make space
	copy(slice[index+1:], slice[index:])
	slice[index] = value
	return slice
}
