package contracts

// Crystalizer is an interface for data structures that can be optimized for
// read operations by "crystalizing" their internal structure.
type Crystalizer interface {
	// Crystalize immediately optimizes the data structure for future read
	// operations. Future write operations may be slower, or not possible at all
	// after calling this method. Check the documentation of the specific data
	// structure for details.
	Crystalize()
}
