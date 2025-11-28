package pow2grid

type Pos struct {
	X, Y uint32
}

type Entry[T any] struct {
	Pos
	Value *T
}

type EntriesMove[T any] struct {
	Old []Entry[T]
	New []Entry[T]
}

func NewEntriesMove[T any](capHint int) EntriesMove[T] {
	return EntriesMove[T]{
		Old: make([]Entry[T], 0, capHint),
		New: make([]Entry[T], 0, capHint),
	}
}

func (u *EntriesMove[T]) Append(value *T, oldPos, newPos Pos) {
	if value == nil {
		return
	}

	u.Old = append(u.Old, Entry[T]{
		Pos:   oldPos,
		Value: value,
	})

	u.New = append(u.New, Entry[T]{
		Pos:   newPos,
		Value: value,
	})
}

type AABB struct {
	Min Pos
	Max Pos
}

// Index is a discrete spatial index over a 2D power-of-two grid.
// It stores objects at integer coordinates and supports point lookups,
// range queries (AABB) and bulk operations (insert, remove, move).
type Index[T any] interface {
	// BulkInsert – insert many objects at once.
	BulkInsert(entries []Entry[T])

	// BulkRemove – remove whatever is stored at the given positions.
	BulkRemove(entries []Entry[T])

	// BulkMove – update objects (typically same Value, different XY).
	BulkMove(moves EntriesMove[T])

	// Get – single lookup at position (x,y).
	Get(x, y uint32) (*T, bool)

	// QueryRange – all objects within the AABB.
	QueryRange(aabb AABB, out []*T) int

	// Count – number of objects in the structure.
	Count() uint64

	// Bounds – global bounds of the handled space.
	Bounds() AABB
}
