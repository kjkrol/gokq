package grid

import "github.com/kjkrol/gokq/pkg/pow2grid"

type Grid[T any] struct {
	cells    []*T
	maxCoord uint32
	count    uint64
}

// Compile-time guard: ensures LinearQuadTree implements SpatialIndex.
var _ pow2grid.Index[any] = (*Grid[any])(nil)

func NewGrid[T any](resolution pow2grid.Resolution) *Grid[T] {
	cellCount := resolution.Cells()
	return &Grid[T]{
		cells:    make([]*T, cellCount),
		maxCoord: resolution.MaxCoord(),
		count:    0,
	}
}

func (g *Grid[T]) BulkInsert(entries []pow2grid.Entry[T]) {
	for _, entry := range entries {
		g.SingleBulkInsert(entry)
	}
}

func (g *Grid[T]) SingleBulkInsert(entry pow2grid.Entry[T]) {
	if !g.inBounds(entry.X, entry.Y) || entry.Value == nil {
		return
	}
	g.setCell(pow2grid.NewMortonCode(entry.X, entry.Y), entry.Value)
}

func (g *Grid[T]) BulkRemove(entities []pow2grid.Entry[T]) {
	for _, entry := range entities {
		g.SignleBulkRemove(entry)
	}
}

func (g *Grid[T]) SignleBulkRemove(entry pow2grid.Entry[T]) {
	pos := entry.Pos
	if !g.inBounds(pos.X, pos.Y) {
		return
	}
	code := pow2grid.NewMortonCode(pos.X, pos.Y)
	if entry.Value == nil || g.cells[code] == entry.Value {
		g.setCell(code, nil)
	}
}

func (g *Grid[T]) BulkMove(moves pow2grid.EntriesMove[T]) {
	g.BulkRemove(moves.Old)
	g.BulkInsert(moves.New)
}

// Get – O(1)
func (g *Grid[T]) Get(x, y uint32) (*T, bool) {
	if x > g.maxCoord || y > g.maxCoord {
		return nil, false
	}
	code := pow2grid.NewMortonCode(x, y)
	val := g.cells[code]
	if val == nil {
		return nil, false
	}
	return val, true
}

func (g *Grid[T]) Count() uint64 {
	return g.count
}

func (g *Grid[T]) Bounds() pow2grid.AABB {
	return pow2grid.AABB{
		Min: pow2grid.Pos{X: 0, Y: 0},
		Max: pow2grid.Pos{X: g.maxCoord, Y: g.maxCoord},
	}
}

func (g *Grid[T]) QueryRange(aabb pow2grid.AABB, out []*T) int {
	if len(out) == 0 {
		return 0
	}

	clear(out)

	if g.count == 0 {
		return 0
	}

	if !g.clampToBound(&aabb) {
		return 0
	}

	limit := len(out)
	written := 0
	pow2grid.MortonCodeAreaConsume(aabb, func(idx int, code pow2grid.MortonCode) {
		if written >= limit {
			return
		}
		if val := g.cells[code]; val != nil {
			out[written] = val
			written++
		}
	})

	return written
}

func (g *Grid[T]) inBounds(x, y uint32) bool {
	return x <= g.maxCoord && y <= g.maxCoord
}

func (g *Grid[T]) clampToBound(aabb *pow2grid.AABB) bool {
	// Clamp AABB to tree bounds
	if aabb.Min.X > g.maxCoord || aabb.Min.Y > g.maxCoord {
		return false
	}
	if aabb.Max.X > g.maxCoord {
		aabb.Max.X = g.maxCoord
	}
	if aabb.Max.Y > g.maxCoord {
		aabb.Max.Y = g.maxCoord
	}
	return true
}

func (g *Grid[T]) setCell(code pow2grid.MortonCode, value *T) {
	idx := int(code)
	prev := g.cells[idx]

	switch {
	case prev == nil && value != nil:
		g.count++
	case prev != nil && value == nil:
		g.count--
	}

	g.cells[idx] = value
}
