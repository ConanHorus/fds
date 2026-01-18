# Coding Style Guide

This document describes the coding conventions used throughout the FDS (Fast Data Structures) repository.

## General Principles

- **Performance First**: All decisions prioritize speed and memory efficiency.
- **Clean Abstractions**: Implementation details (e.g., buddy allocators) are hidden from users.
- **Explicit Over Implicit**: Named return values, explicit type constraints, clear error handling.

## Naming Conventions

### Types and Structs

- **Exported types**: PascalCase with descriptive names (`LinkedList`, `BuddyAllocator`)
- **Generic constraints**: Use semantic constraints (`comparable`, `any`) based on actual requirements
- **Options structs**: Suffix with `Options` (`BuddyAllocatorOptions`)

### Functions and Methods

- **Constructors**: Prefix with `New` (`NewLinkedList`, `NewBuddyAllocator`)
- **Factory from data**: `New<Type>From<Source>` pattern (`NewLinkedListFromSlice`)
- **Option functions**: Prefix with `With` (`WithInitialCapacity`, `WithMaxOrder`)
- **Boolean queries**: Use verbs or adjectives (`IsEmpty`, `Contains`)
- **Accessors**: Simple noun names (`Length`, `Capacity`)
- **Mutators**: Action verbs (`PushBack`, `PopFront`, `InsertAt`, `RemoveAt`)

### Variables and Constants

- **Private constants**: camelCase, grouped in `const` blocks
- **Sentinel values**: Named constants with clear meaning (`nilIndex = -1`)
- **Capacities**: Named constants (`minCapacity = 16`)

### Receivers

- **Use `this`** for method receivers (not single letters like `l` or `ba`)

```go
func (this *LinkedList[T]) Length() int {
    return this.length
}
```

## Documentation Style

### Godoc Comments

Every exported symbol must have a godoc comment following this structure:

1. **First line**: Brief description starting with the symbol name
2. **Description**: Detailed explanation if needed
3. **Tradeoffs**: (For complex types/constructors) Document performance tradeoffs
4. **Examples**: (When helpful) Show usage with `Examples:` block
5. **Parameters**: Document each parameter with `Parameters:` block
6. **Returns**: Document return values with `Returns:` block

```go
// Allocate allocates a memory block of the specified size.
//
// If necessary, the allocator will grow its memory to accommodate the request,
// unless a maximum order has been set and the request exceeds it, in which case
// it will fail.
//
// Examples:
//   - Allocate(1) allocates a 1-unit block
//   - Allocate(5) allocates an 8-unit block (rounded up to next power of 2)
//   - Allocate(0) returns failure
//
// Parameters:
//   - size: The size of the memory block to allocate (must be > 0)
//
// Returns:
//   - index: The starting index of the allocated memory block
//   - grew: A boolean indicating whether the allocator had to grow its memory
//   - ok: A boolean indicating whether the allocation was successful
func (this *BuddyAllocator) Allocate(size uint64) (index uint64, grew bool, ok bool) {
```

### Named Return Values

Always use named return values for clarity:

```go
func (this *LinkedList[T]) Get(index int) (value T, ok bool)
func (this *LinkedList[T]) IsEmpty() (empty bool)
func (this *LinkedList[T]) Length() (length int)
```

## Code Organization

### File Structure

1. Package declaration
2. Imports (grouped: stdlib, external, internal)
3. Constants
4. Types (struct definitions)
5. Constructors
6. Exported methods (alphabetical)
7. `// --- private methods --- //` separator
8. Private methods

### Method Ordering

- Constructors first
- Then exported methods alphabetically
- Private methods at the bottom, separated by a comment

### Import Grouping

```go
import (
    "fmt"
    "strings"

    "github.com/external/package"

    "github.com/ConanHorus/fds/internal/pkg"
)
```

### Internal Packages

Use internal packages for implementation details that shouldn't be exposed:

```go
import _ll "github.com/ConanHorus/fds/ll/internal"
```

Prefix internal package aliases with underscore for clarity.

## Error Handling

### Return Patterns

- **Boolean ok**: Use `(value T, ok bool)` for operations that may fail
- **No errors**: Prefer `ok bool` over `error` for simple success/failure
- **Panic for invariants**: Panic only for programming errors (e.g., allocation failure)

```go
func (this *LinkedList[T]) Get(index int) (value T, ok bool) {
    if index < 0 || index >= this.length {
        return value, false  // Return zero value and false
    }
    return this.values[index], true
}
```

## Testing Style

### Framework

Use the `github.com/smarty/assertions` package:

```go
and := assertions.New(t)
and.So(actual, should.Equal, expected)
and.So(slice, should.Resemble, expectedSlice)
```

### Test Organization

- Group tests by category using `// --- Category --- //` comments
- Use table-driven tests with `map[string]struct{}` pattern
- Name test cases descriptively
- Use `t.Parallel()` for independent tests

### Test Naming

```go
func TestTypeName_MethodName(t *testing.T)
func TestMethodName(t *testing.T)           // When type is obvious
func TestMethodName_Scenario(t *testing.T)  // For specific scenarios
```

## Iteration

### Go 1.23+ Range Functions

Implement `iter.Seq[T]` and `iter.Seq2[K, V]` for iteration:

```go
func (this *LinkedList[T]) All() iter.Seq[T] {
    return func(yield func(T) bool) {
        for _, innerIndex := range this.allIndexes() {
            if !yield(this.values[innerIndex]) {
                return
            }
        }
    }
}
```

### Callback Patterns

Also provide callback-style iteration with early-exit support:

```go
func (this *LinkedList[T]) ForEach(delegate func(value T) (ok bool))
func (this *LinkedList[T]) ForEachIndexed(delegate func(index int, value T) (ok bool))
```

## Memory Management

### Buddy Allocator Integration

Data structures use `BuddyAllocator` internally but hide this from users:

- Users interact with a clean API (`PushBack`, `Get`, etc.)
- Growth is automatic and transparent
- Internal indices are mapped to user-facing indices

### Capacity Planning

- Define sensible minimum capacities (`minCapacity = 16`)
- Pre-allocate with headroom for known sizes (`len(slice) + len(slice)/4 + 1`)
