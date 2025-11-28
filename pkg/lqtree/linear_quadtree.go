package lqtree

import "github.com/kjkrol/gokq/pkg/pow2grid"

// LinearQuadTree is a linear (array-backed) quadtree over a 2D power-of-two
// grid. It implements pow2grid.Index[T] by storing objects at integer
// coordinates and providing point lookups, AABB range queries, and bulk
// insert/move/remove operations.
type LinearQuadTree[T any] struct {
	cells    []*T
	maxCoord uint32
	depth    pow2grid.Resolution
	count    uint64
}

// Compile-time guard: ensures LinearQuadTree implements SpatialIndex.
var _ pow2grid.Index[any] = (*LinearQuadTree[any])(nil)

func NewLinearQuadTree[T any](resolution pow2grid.Resolution) *LinearQuadTree[T] {
	cellCount := resolution.Cells()
	return &LinearQuadTree[T]{
		cells:    make([]*T, cellCount),
		maxCoord: resolution.MaxCoord(),
		depth:    resolution,
		count:    0,
	}
}

func (qt *LinearQuadTree[T]) BulkInsert(entries []pow2grid.Entry[T]) {
	for _, entry := range entries {
		qt.SingleBulkInsert(entry)
	}
}

func (qt *LinearQuadTree[T]) SingleBulkInsert(entry pow2grid.Entry[T]) {
	if !qt.inBounds(entry.X, entry.Y) || entry.Value == nil {
		return
	}
	qt.setCell(pow2grid.NewMortonCode(entry.X, entry.Y), entry.Value)
}

func (qt *LinearQuadTree[T]) BulkRemove(entities []pow2grid.Entry[T]) {
	for _, entry := range entities {
		qt.SignleBulkRemove(entry)
	}
}

func (qt *LinearQuadTree[T]) SignleBulkRemove(entry pow2grid.Entry[T]) {
	pos := entry.Pos
	if !qt.inBounds(pos.X, pos.Y) {
		return
	}
	code := pow2grid.NewMortonCode(pos.X, pos.Y)
	if entry.Value == nil || qt.cells[code] == entry.Value {
		qt.setCell(code, nil)
	}
}

func (qt *LinearQuadTree[T]) BulkMove(moves pow2grid.EntriesMove[T]) {
	qt.BulkRemove(moves.Old)
	qt.BulkInsert(moves.New)
}

// Get – O(1)
func (qt *LinearQuadTree[T]) Get(x, y uint32) (*T, bool) {
	if x > qt.maxCoord || y > qt.maxCoord {
		return nil, false
	}
	code := pow2grid.NewMortonCode(x, y)
	val := qt.cells[code]
	if val == nil {
		return nil, false
	}
	return val, true
}

func (qt *LinearQuadTree[T]) Count() uint64 {
	return qt.count
}

func (qt *LinearQuadTree[T]) Bounds() pow2grid.AABB {
	return pow2grid.AABB{
		Min: pow2grid.Pos{X: 0, Y: 0},
		Max: pow2grid.Pos{X: qt.maxCoord, Y: qt.maxCoord},
	}
}

func (qt *LinearQuadTree[T]) QueryRange(aabb pow2grid.AABB, out []*T) int {
	if len(out) == 0 {
		return 0
	}

	clear(out)

	if qt.count == 0 {
		return 0
	}

	if !qt.clampToBound(&aabb) {
		return 0
	}

	limit := len(out)
	written := 0
	pow2grid.MortonCodeAreaConsume(aabb, func(idx int, code pow2grid.MortonCode) {
		if written >= limit {
			return
		}
		if val := qt.cells[code]; val != nil {
			out[written] = val
			written++
		}
	})

	return written
}

func (qt *LinearQuadTree[T]) inBounds(x, y uint32) bool {
	return x <= qt.maxCoord && y <= qt.maxCoord
}

func (qt *LinearQuadTree[T]) clampToBound(aabb *pow2grid.AABB) bool {
	// Clamp AABB to tree bounds
	if aabb.Min.X > qt.maxCoord || aabb.Min.Y > qt.maxCoord {
		return false
	}
	if aabb.Max.X > qt.maxCoord {
		aabb.Max.X = qt.maxCoord
	}
	if aabb.Max.Y > qt.maxCoord {
		aabb.Max.Y = qt.maxCoord
	}
	return true
}

func (qt *LinearQuadTree[T]) setCell(code pow2grid.MortonCode, value *T) {
	idx := int(code)
	prev := qt.cells[idx]

	switch {
	case prev == nil && value != nil:
		qt.count++
	case prev != nil && value == nil:
		qt.count--
	}

	qt.cells[idx] = value
}
