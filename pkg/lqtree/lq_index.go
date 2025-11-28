package lqtree

import "github.com/kjkrol/gokq/pkg/pow2grid"

// LQIndex to implementacja pow2grid.Index[T] na bazie Mortonowego quadtree.
type LQIndex[T any] struct {
	tree   *Tree[T]
	bounds pow2grid.AABB
	count  uint64
}

func NewLQIndex[T any](maxLevel uint8) *LQIndex[T] {
	size := uint32(1) << maxLevel
	bounds := pow2grid.AABB{
		Min: pow2grid.Pos{X: 0, Y: 0},
		Max: pow2grid.Pos{X: size - 1, Y: size - 1},
	}
	return &LQIndex[T]{
		tree:   NewTree[T](maxLevel),
		bounds: bounds,
	}
}

// --- pomocnicze AABB/Pos ---

func posInside(a pow2grid.AABB, p pow2grid.Pos) bool {
	return p.X >= a.Min.X && p.X <= a.Max.X &&
		p.Y >= a.Min.Y && p.Y <= a.Max.Y
}

func aabbIntersects(a, b pow2grid.AABB) bool {
	return a.Min.X <= b.Max.X && a.Max.X >= b.Min.X &&
		a.Min.Y <= b.Max.Y && a.Max.Y >= b.Min.Y
}

func clampAABB(a, bounds pow2grid.AABB) pow2grid.AABB {
	if a.Min.X < bounds.Min.X {
		a.Min.X = bounds.Min.X
	}
	if a.Min.Y < bounds.Min.Y {
		a.Min.Y = bounds.Min.Y
	}
	if a.Max.X > bounds.Max.X {
		a.Max.X = bounds.Max.X
	}
	if a.Max.Y > bounds.Max.Y {
		a.Max.Y = bounds.Max.Y
	}
	return a
}

// --- implementacja pow2grid.Index[T] ---

// BulkInsert – wstawia (nadpisuje) wartości w pozycjach z entries.
func (idx *LQIndex[T]) BulkInsert(entries []pow2grid.Entry[T]) {
	for i := range entries {
		e := entries[i]
		if e.Value == nil {
			continue
		}
		if !posInside(idx.bounds, e.Pos) {
			continue
		}

		added := idx.tree.insertPoint(e.Pos, e.Value)
		if added {
			idx.count++
		}
	}
}

// BulkRemove – usuwa cokolwiek jest pod podanymi pozycjami.
func (idx *LQIndex[T]) BulkRemove(entries []pow2grid.Entry[T]) {
	for i := range entries {
		e := entries[i]
		if !posInside(idx.bounds, e.Pos) {
			continue
		}
		if _, removed := idx.tree.removePoint(e.Pos); removed {
			if idx.count > 0 {
				idx.count--
			}
		}
	}
}

// BulkMove – usuwa spod Old[i].Pos i wstawia pod New[i].Pos.
// Zakładam, że Old[i].Value == New[i].Value – zgodnie z Twoim EntriesMove.
func (idx *LQIndex[T]) BulkMove(moves pow2grid.EntriesMove[T]) {
	n := len(moves.Old)
	if len(moves.New) < n {
		n = len(moves.New)
	}
	for i := 0; i < n; i++ {
		old := moves.Old[i]
		new := moves.New[i]
		if old.Value == nil {
			continue
		}

		// Usuń ze starego miejsca (jeśli w bounds).
		if posInside(idx.bounds, old.Pos) {
			if _, removed := idx.tree.removePoint(old.Pos); removed {
				if idx.count > 0 {
					idx.count--
				}
			}
		}

		// Wstaw do nowego miejsca (jeśli w bounds).
		if posInside(idx.bounds, new.Pos) {
			added := idx.tree.insertPoint(new.Pos, new.Value)
			if added {
				idx.count++
			}
		}
	}
}

// Get – pojedynczy lookup.
func (idx *LQIndex[T]) Get(x, y uint32) (*T, bool) {
	pos := pow2grid.Pos{X: x, Y: y}
	if !posInside(idx.bounds, pos) {
		return nil, false
	}
	return idx.tree.getPoint(pos)
}

// QueryRange – zwróć wszystkie obiekty w AABB, do len(out).
func (idx *LQIndex[T]) QueryRange(aabb pow2grid.AABB, out []*T) int {
	if len(out) == 0 {
		return 0
	}
	if !aabbIntersects(idx.bounds, aabb) {
		return 0
	}

	clamped := clampAABB(aabb, idx.bounds)
	return idx.tree.QueryRange(clamped, out)
}

func (idx *LQIndex[T]) Count() uint64 {
	return idx.count
}

func (idx *LQIndex[T]) Bounds() pow2grid.AABB {
	return idx.bounds
}
