# FDS Documentation Map

Fast Data Structures (FDS) is a Go library providing high-performance, memory-efficient data structures.

## Package Overview

| Package | Directory | Description |
|---------|-----------|-------------|
| `ll` | [ll/](ll/ll.md) | Linked list implementations |
| `mm` | [mm/](mm/mm.md) | Memory management primitives |

## Linked Lists (`ll`)

The `ll` package provides generic linked list implementations optimized for performance.

**Contents:**
- `LinkedList[T comparable]` — Doubly-linked list with O(1) end operations

See [ll/ll.md](ll/ll.md) for full documentation.

## Memory Managers (`mm`)

The `mm` package provides memory allocation primitives used internally by other data structures.

**Contents:**
- `BuddyAllocator` — Power-of-two block allocator with automatic coalescing

See [mm/mm.md](mm/mm.md) for full documentation.

## Additional Resources

- [Coding Style Guide](coding_style.md) — Conventions used in this codebase
