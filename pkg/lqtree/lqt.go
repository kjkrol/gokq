package lqtree

// -----------------------------------------
// Raw Linear QuadTree
// -----------------------------------------
type LinearQuadTree[T any] struct {
	cells []*T
	maxXY uint32
	depth uint8
	count uint64
}

// Compile-time guard: ensures LinearQuadTree implements SpatialIndex.
var _ SpatialIndex[any] = (*LinearQuadTree[any])(nil)

func NewLinearQuadTree[T any](size LQTSize) *LinearQuadTree[T] {
	maxCoord := size.Resolution()
	depth := uint8(size)
	cellCount := size.ArraySize()

	return &LinearQuadTree[T]{
		cells: make([]*T, cellCount),
		maxXY: maxCoord,
		depth: depth,
		count: 0,
	}
}

func (qt *LinearQuadTree[T]) BulkInsert(entries []Entry[T]) {
	for i := range entries {
		entry := entries[i]
		if !qt.inBounds(entry.X, entry.Y) || entry.Value == nil {
			continue
		}
		qt.setCell(NewMortonCode(entry.X, entry.Y), entry.Value)
	}
}

func (qt *LinearQuadTree[T]) BulkRemove(positions []Pos) {
	for i := range positions {
		pos := positions[i]
		if !qt.inBounds(pos.X, pos.Y) {
			continue
		}

		code := NewMortonCode(pos.X, pos.Y)
		if qt.cells[code] != nil {
			qt.setCell(code, nil)
		}
	}
}

func (qt *LinearQuadTree[T]) BulkMove(moves []EntryMove[T]) {
	for i := range moves {
		move := moves[i]

		oldInBounds := qt.inBounds(move.Old.X, move.Old.Y)
		newInBounds := qt.inBounds(move.New.X, move.New.Y)

		var value *T = move.Value

		if oldInBounds {
			oldCode := NewMortonCode(move.Old.X, move.Old.Y)
			if value == nil {
				value = qt.cells[oldCode]
			}
			if qt.cells[oldCode] != nil {
				qt.setCell(oldCode, nil)
			}
		}

		if newInBounds && value != nil {
			qt.setCell(NewMortonCode(move.New.X, move.New.Y), value)
		}
	}
}

// Get – O(1)
func (qt *LinearQuadTree[T]) Get(x, y uint32) (*T, bool) {
	if x > qt.maxXY || y > qt.maxXY {
		return nil, false
	}
	code := NewMortonCode(x, y)
	val := qt.cells[code]
	if val == nil {
		return nil, false
	}
	return val, true
}

func (qt *LinearQuadTree[T]) Count() uint64 {
	return qt.count
}

func (qt *LinearQuadTree[T]) Bounds() AABB {
	return AABB{
		Min: Pos{X: 0, Y: 0},
		Max: Pos{X: qt.maxXY, Y: qt.maxXY},
	}
}

func (qt *LinearQuadTree[T]) QueryRange(aabb AABB) []*T {
	if qt.count == 0 {
		return nil
	}

	// Clamp AABB to tree bounds
	if aabb.Min.X > qt.maxXY || aabb.Min.Y > qt.maxXY {
		return nil
	}
	if aabb.Max.X > qt.maxXY {
		aabb.Max.X = qt.maxXY
	}
	if aabb.Max.Y > qt.maxXY {
		aabb.Max.Y = qt.maxXY
	}

	results := make([]*T, 0)
	// Start: level 0, whole world [0..2^depth-1] x [0..2^depth-1], prefix=0
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
	prefix MortonCode,
	aabb AABB,
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

		remainBits := sizeBits         // tyle poziomów poniżej
		shift := remainBits * 2        // 2 bity na poziom
		start := prefix << shift       // wspólny prefiks
		span := MortonCode(1) << shift // liczba kodów w poddrzewie
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

func (qt *LinearQuadTree[T]) setCell(code MortonCode, value *T) {
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
