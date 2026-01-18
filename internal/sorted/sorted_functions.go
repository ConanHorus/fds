package sorted

// Sortable is a constraint that permits any integer type.
type Sortable interface {
	~uint8 | ~int8 | ~uint16 | ~int16 | ~uint32 | ~int32 | ~uint64 | ~int64 | ~uint | ~int | ~uintptr
}

// BinarySearch performs a binary search on a sorted slice of sortable
// elements to find the target value's index. The first occurence index is
// returned if duplicates exist. If the target is not found, the index where it
// can be inserted to maintain sorted order is returned.
//
// Parameters:
//   - slice: A sorted slice of elements of type T.
//   - target: The value to search for within the slice.
//
// Returns:
//   - index: The index of the first occurrence of the target value if found;
//     otherwise, the index where the target can be inserted to maintain sorted order.
//   - found: A boolean indicating whether the target value was found in the
//     slice.
func BinarySearch[T Sortable](slice []T, target T) (index int, found bool) {
	if len(slice) == 0 {
		return 0, false
	}

	low := 0
	high := len(slice)
	for low < high {
		mid := ((high - low) / 2) + low
		if slice[mid] < target {
			low = mid + 1
			continue
		}

		high = mid
	}

	if low < len(slice) && slice[low] == target {
		return low, true
	}

	return low, false
}

// GallopingSearch performs a galloping search on a sorted slice of sortable
// elements to find the target value.
//
// For small slices (length less than a defined threshold), it falls back to
// binary search for efficiency. Galloping search is particularly effective
// for large datasets where the target is likely to be found in exponentially
// growing intervals.
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
func GallopingSearch[T Sortable](slice []T, target T) (index int, found bool) {
	const threshold = 64
	if len(slice) < threshold {
		return BinarySearch(slice, target)
	}

	_ = slice[len(slice)-1] // Bounds check elimination hint

	spanSize := threshold
	prevIndex := 0
	for {
		currentIndex := prevIndex + spanSize
		if currentIndex < 0 || currentIndex >= len(slice) {
			currentIndex = len(slice) - 1
			index, found := BinarySearch(slice[prevIndex:currentIndex+1], target)
			return prevIndex + index, found
		}

		if slice[currentIndex] < target {
			prevIndex = currentIndex
			if spanSize > len(slice)/2 {
				spanSize = len(slice) - prevIndex
				continue
			}

			spanSize = spanSize << 1
			continue
		}

		index, found := BinarySearch(slice[prevIndex:currentIndex+1], target)
		return prevIndex + index, found
	}
}

// Insert inserts a value into a sorted slice of sortable elements while
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
func Insert[T Sortable](slice []T, value T, allowDuplicates bool) []T {
	index, found := GallopingSearch(slice, value)
	if found && !allowDuplicates {
		return slice
	}

	slice = append(slice, *new(T)) // Add zero value to grow slice
	copy(slice[index+1:], slice[index:len(slice)-1])
	slice[index] = value
	return slice
}
