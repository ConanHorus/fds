package contracts

// TreeReader is an interface for reading from a tree data structure.
type TreeReader[T any] interface {
	// Contains reports whether the tree contains the specified value.
	//
	// Parameters:
	//   - value: The value to search for.
	//
	// Returns:
	//   - ok: true if the value is found, false otherwise.
	Contains(value T) (ok bool)

	// Min returns the minimum value in the tree.
	//
	// Returns:
	//   - value: The minimum value, or the zero value if the tree is empty.
	//   - ok: true if the tree is non-empty, false otherwise.
	Min() (value T, ok bool)

	// Max returns the maximum value in the tree.
	//
	// Returns:
	//   - value: The maximum value, or the zero value if the tree is empty.
	//   - ok: true if the tree is non-empty, false otherwise.
	Max() (value T, ok bool)

	// Length returns the number of elements in the tree.
	//
	// Returns:
	//   - length: The number of elements currently in the tree.
	Length() (length int)

	// IsEmpty reports whether the tree contains no elements.
	//
	// Returns:
	//   - empty: true if the tree has no elements, false otherwise.
	IsEmpty() (empty bool)
}

// TreeWriter is an interface for writing to a tree data structure.
type TreeWriter[T any] interface {
	// Insert adds a value to the tree.
	//
	// Parameters:
	//   - value: The value to insert.
	//
	// Returns:
	//   - ok: true if the value was inserted, false if it already exists.
	Insert(value T) (ok bool)
}

// TreeDeleter is an interface for deleting from a tree data structure.
type TreeDeleter[T any] interface {
	// Delete removes a value from the tree.
	//
	// Parameters:
	//   - value: The value to remove.
	//
	// Returns:
	//   - ok: true if the value was found and removed, false otherwise.
	Delete(value T) (ok bool)

	// Clear removes all elements from the tree.
	Clear()
}
