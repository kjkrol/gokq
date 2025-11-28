package lqtree

import "github.com/kjkrol/gokq/pkg/pow2grid"

type Node[T any] struct {
	Code  pow2grid.MortonCode // prefiks (2*Level bitów)
	Level uint8
	Count int
	Items []*T // dla tego szkicu: tylko liście (Level == MaxLevel) trzymają wartości
}

type Tree[T any] struct {
	MaxLevel uint8
	nodes    map[uint8]map[pow2grid.MortonCode]*Node[T] // level -> (prefix -> node)
}

func NewTree[T any](maxLevel uint8) *Tree[T] {
	return &Tree[T]{
		MaxLevel: maxLevel,
		nodes:    make(map[uint8]map[pow2grid.MortonCode]*Node[T]),
	}
}

func NodePrefixFromFull(full pow2grid.MortonCode, level, maxLevel uint8) pow2grid.MortonCode {
	shift := 2 * (maxLevel - level)
	return pow2grid.MortonCode(uint64(full) >> shift)
}

func ChildPrefix(parent pow2grid.MortonCode, childID uint8) pow2grid.MortonCode {
	return pow2grid.MortonCode((uint64(parent) << 2) | uint64(childID&0b11))
}

func NodeAABB(prefix pow2grid.MortonCode, level, maxLevel uint8) pow2grid.AABB {
	x, y := prefix.Decode() // x,y ∈ [0, 2^level)

	scale := uint32(1) << (maxLevel - level)

	minX := x * scale
	minY := y * scale
	maxX := (x+1)*scale - 1
	maxY := (y+1)*scale - 1

	return pow2grid.AABB{
		Min: pow2grid.Pos{X: minX, Y: minY},
		Max: pow2grid.Pos{X: maxX, Y: maxY},
	}
}

// --- helpers na mapach ---

func (t *Tree[T]) getNode(level uint8, code pow2grid.MortonCode) *Node[T] {
	m := t.nodes[level]
	if m == nil {
		return nil
	}
	return m[code]
}

func (t *Tree[T]) getOrCreateNode(level uint8, code pow2grid.MortonCode) *Node[T] {
	m := t.nodes[level]
	if m == nil {
		m = make(map[pow2grid.MortonCode]*Node[T])
		t.nodes[level] = m
	}
	n := m[code]
	if n == nil {
		n = &Node[T]{
			Code:  code,
			Level: level,
		}
		m[code] = n
	}
	return n
}

// --- pojedyncze operacje na punktach (x,y) ---

// insertPoint wstawia *T pod pos.
// Zwraca true, jeśli to NOWY wpis (zwiększa Count), false jeśli nadpisaliśmy istniejący.
func (t *Tree[T]) insertPoint(pos pow2grid.Pos, value *T) bool {
	full := pow2grid.NewMortonCode(pos.X, pos.Y)
	leafLevel := t.MaxLevel
	leafPrefix := NodePrefixFromFull(full, leafLevel, t.MaxLevel)

	leaf := t.getNode(leafLevel, leafPrefix)
	if leaf != nil && len(leaf.Items) > 0 {
		// w tym szkicu: jedno value na liść – nadpisujemy bez zmiany Count
		leaf.Items[0] = value
		return false
	}

	// nowy wpis – dodajemy ścieżkę od korzenia do liścia,
	// przy każdym węźle zwiększając Count
	for level := uint8(0); level <= t.MaxLevel; level++ {
		prefix := NodePrefixFromFull(full, level, t.MaxLevel)
		n := t.getOrCreateNode(level, prefix)
		n.Count++
		if level == t.MaxLevel {
			n.Items = []*T{value}
		}
	}

	return true
}

// removePoint usuwa wpis pod pos, jeżeli istnieje.
// Zwraca usuniętą wartość i bool, czy coś faktycznie usunięto.
func (t *Tree[T]) removePoint(pos pow2grid.Pos) (*T, bool) {
	full := pow2grid.NewMortonCode(pos.X, pos.Y)
	leafLevel := t.MaxLevel
	leafPrefix := NodePrefixFromFull(full, leafLevel, t.MaxLevel)

	leaf := t.getNode(leafLevel, leafPrefix)
	if leaf == nil || len(leaf.Items) == 0 {
		return nil, false
	}

	removed := leaf.Items[0]
	leaf.Items = leaf.Items[:0]

	// schodzimy z liścia do korzenia, dekrementując Count
	for level := leafLevel; ; level-- {
		prefix := NodePrefixFromFull(full, level, t.MaxLevel)
		m := t.nodes[level]
		if m == nil {
			break
		}
		n := m[prefix]
		if n == nil {
			break
		}
		n.Count--
		if n.Count <= 0 {
			delete(m, prefix)
			if len(m) == 0 {
				delete(t.nodes, level)
			}
		}
		if level == 0 {
			break
		}
	}

	return removed, true
}

// getPoint – pojedynczy lookup.
func (t *Tree[T]) getPoint(pos pow2grid.Pos) (*T, bool) {
	full := pow2grid.NewMortonCode(pos.X, pos.Y)
	leafPrefix := NodePrefixFromFull(full, t.MaxLevel, t.MaxLevel)
	leaf := t.getNode(t.MaxLevel, leafPrefix)
	if leaf == nil || len(leaf.Items) == 0 {
		return nil, false
	}
	return leaf.Items[0], true
}

// --- QueryRange na drzewie ---

func intersects(a, b pow2grid.AABB) bool {
	return a.Min.X <= b.Max.X && a.Max.X >= b.Min.X &&
		a.Min.Y <= b.Max.Y && a.Max.Y >= b.Min.Y
}

func (t *Tree[T]) QueryRange(aabb pow2grid.AABB, out []*T) int {
	if len(out) == 0 {
		return 0
	}

	// zaczynamy od korzenia
	level0 := uint8(0)
	m := t.nodes[level0]
	if m == nil {
		return 0
	}
	root, ok := m[0]
	if !ok {
		return 0
	}

	return t.queryNode(root, aabb, out)
}

func (t *Tree[T]) queryNode(node *Node[T], aabb pow2grid.AABB, out []*T) int {
	nodeBox := NodeAABB(node.Code, node.Level, t.MaxLevel)
	if !intersects(nodeBox, aabb) {
		return 0
	}

	// liść (poziom MaxLevel) – z jego AABB wynika, że kratka leży w zakresie
	if node.Level == t.MaxLevel {
		written := 0
		for _, v := range node.Items {
			if v == nil {
				continue
			}
			if written >= len(out) {
				break
			}
			out[written] = v
			written++
		}
		return written
	}

	// węzeł wewnętrzny – schodzimy do dzieci
	childLevel := node.Level + 1
	m := t.nodes[childLevel]
	if m == nil {
		return 0
	}

	written := 0
	for childID := uint8(0); childID < 4 && written < len(out); childID++ {
		childCode := ChildPrefix(node.Code, childID)
		childNode, ok := m[childCode]
		if !ok || childNode.Count == 0 {
			continue
		}
		w := t.queryNode(childNode, aabb, out[written:])
		written += w
	}
	return written
}
