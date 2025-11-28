package lqtree

import "github.com/kjkrol/gokq/pkg/pow2grid"

// LinearQuadTree is a linear (array-backed) quadtree over a 2D power-of-two
// grid. It implements pow2grid.Index[T] by storing objects at integer
// coordinates and providing point lookups, AABB range queries, and bulk
// insert/move/remove operations.
type LinearQuadTree[T any] struct {
	cells []*T
	maxXY uint32
	depth uint8
	count uint64
}

// Compile-time guard: ensures LinearQuadTree implements SpatialIndex.
var _ pow2grid.Index[any] = (*LinearQuadTree[any])(nil)

func NewLinearQuadTree[T any](resolution pow2grid.Resolution) *LinearQuadTree[T] {
	maxCoord := resolution.Side()
	depth := uint8(resolution)
	cellCount := resolution.Cells()

	return &LinearQuadTree[T]{
		cells: make([]*T, cellCount),
		maxXY: maxCoord,
		depth: depth,
		count: 0,
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
	old := moves.Old
	new := moves.New

	for i := range old {
		pos := old[i].Pos
		if qt.inBounds(pos.X, pos.Y) {
			code := pow2grid.NewMortonCode(pos.X, pos.Y)
			qt.setCell(code, nil)
		}
	}

	for i := range new {
		pos := new[i].Pos
		if !qt.inBounds(pos.X, pos.Y) || new[i].Value == nil {
			continue
		}
		code := pow2grid.NewMortonCode(pos.X, pos.Y)
		qt.setCell(code, new[i].Value)
	}
}

// Get – O(1)
func (qt *LinearQuadTree[T]) Get(x, y uint32) (*T, bool) {
	if x > qt.maxXY || y > qt.maxXY {
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
		Max: pow2grid.Pos{X: qt.maxXY, Y: qt.maxXY},
	}
}

func (qt *LinearQuadTree[T]) QueryRange(aabb pow2grid.AABB, out []*T) []*T {
	if qt.count == 0 {
		return out
	}

	// Clamp AABB to tree bounds
	if aabb.Min.X > qt.maxXY || aabb.Min.Y > qt.maxXY {
		return out
	}
	if aabb.Max.X > qt.maxXY {
		aabb.Max.X = qt.maxXY
	}
	if aabb.Max.Y > qt.maxXY {
		aabb.Max.Y = qt.maxXY
	}

	// Start: level 0, whole world [0..2^depth-1] x [0..2^depth-1], prefix=0
	results := out
	qt.queryNode(0, 0, 0, 0, aabb, &results)
	return results
}

// queryNode traverses a quadtree node:
// level – how many levels below the root (0..depth)
// prefix – shared MortonCode prefix for the subtree of this node
// region: [x0..x1] x [y0..y1], size derived from (depth - level)
func (qt *LinearQuadTree[T]) queryNode(
	x0, y0 uint32,
	level uint8,
	prefix pow2grid.MortonCode,
	aabb pow2grid.AABB,
	out *[]*T,
) {
	// how many bits remain downward (subtree depth)
	sizeBits := qt.depth - level
	size := uint32(1) << sizeBits
	x1 := x0 + size - 1
	y1 := y0 + size - 1

	// 1. No intersection with AABB → stop
	if x1 < aabb.Min.X || x0 > aabb.Max.X ||
		y1 < aabb.Min.Y || y0 > aabb.Max.Y {
		return
	}

	// 2. Region fully inside AABB → scan the entire Morton range for this subtree.
	if x0 >= aabb.Min.X && x1 <= aabb.Max.X &&
		y0 >= aabb.Min.Y && y1 <= aabb.Max.Y {

		remainBits := sizeBits                  // tyle poziomów poniżej
		shift := remainBits * 2                 // 2 bity na poziom
		start := prefix << shift                // wspólny prefiks
		span := pow2grid.MortonCode(1) << shift // liczba kodów w poddrzewie
		end := start + span - 1

		for code := start; code <= end; code++ {
			obj := qt.cells[int(code)]
			if obj != nil {
				*out = append(*out, obj)
			}
		}
		return
	}

	// 3. Leaf (1x1 cell) + partial overlap → check the single cell
	if sizeBits == 0 {
		code := prefix // pełny MortonCode
		if x0 >= aabb.Min.X && x0 <= aabb.Max.X &&
			y0 >= aabb.Min.Y && y0 <= aabb.Max.Y {

			obj := qt.cells[int(code)]
			if obj != nil {
				*out = append(*out, obj)
			}
		}
		return
	}

	// 4. Partial overlap, non-leaf → split into 4 quadrants
	halfBits := sizeBits - 1
	half := uint32(1) << halfBits

	basePrefix := prefix << 2
	nwPrefix := basePrefix     // (dx=0, dy=0)
	nePrefix := basePrefix | 1 // (dx=1, dy=0)
	swPrefix := basePrefix | 2 // (dx=0, dy=1)
	sePrefix := basePrefix | 3 // (dx=1, dy=1)

	nextLevel := level + 1

	// NW
	qt.queryNode(
		x0,
		y0,
		nextLevel,
		nwPrefix,
		aabb,
		out,
	)
	// NE
	qt.queryNode(
		x0+half,
		y0,
		nextLevel,
		nePrefix,
		aabb,
		out,
	)
	// SW
	qt.queryNode(
		x0,
		y0+half,
		nextLevel,
		swPrefix,
		aabb,
		out,
	)
	// SE
	qt.queryNode(
		x0+half,
		y0+half,
		nextLevel,
		sePrefix,
		aabb,
		out,
	)
}

func (qt *LinearQuadTree[T]) inBounds(x, y uint32) bool {
	return x <= qt.maxXY && y <= qt.maxXY
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
