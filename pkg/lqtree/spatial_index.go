package lqtree

type Pos struct {
	X, Y uint32
}

type Entry[T any] struct {
	Pos
	Value *T
}

type EntryMove[T any] struct {
	Old   Pos
	New   Pos
	Value *T
}

type AABB struct {
	Min Pos
	Max Pos
}

type SpatialIndex[T any] interface {
	// BulkInsert – insert many objects at once.
	BulkInsert(entries []Entry[T])

	// BulkRemove – remove whatever is stored at the given positions.
	BulkRemove(positions []Pos)

	// BulkMove – move objects (typically same Value, different XY).
	BulkMove(moves []EntryMove[T])

	// Get – single lookup at position (x,y).
	Get(x, y uint32) (*T, bool)

	// QueryRange – all objects within the AABB.
	QueryRange(aabb AABB) []*T

	// Count – number of objects in the structure.
	Count() uint64

	// Bounds – global bounds of the handled space.
	Bounds() AABB
}
