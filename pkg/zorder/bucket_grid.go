package zorder

import (
	"github.com/kjkrol/gokq/pkg/lqtree"
	"github.com/kjkrol/gokq/pkg/pow2grid"
)

type ChunkKey struct {
	X, Y uint32
}

// BucketGrid shards space into Z-order buckets (chunks) and stores points in per-bucket
// linear quadtrees. It implements pow2grid.Index for Morton-friendly, power-of-two grids.
type BucketGrid[T any] struct {
	side              uint32
	bucketsResolution pow2grid.Resolution
	bound             pow2grid.AABB
	buckets           map[ChunkKey]*lqtree.LinearQuadTree[T]
	count             uint64
}

var _ pow2grid.Index[any] = (*BucketGrid[any])(nil)

// Choose bucketstSize (and thus bucketMaxXY) to roughly match typical query AABB.
// Example: if most queries are <32x32, a 64–128 bucket keeps them in 1–4 buckets with moderate memory;
// QueryRange that spans many buckets will pay for many map lookups + per-bucket scans.
func NewBucketGrid[T any](
	bucketsResolution pow2grid.Resolution,
	bound pow2grid.AABB,
) *BucketGrid[T] {
	return &BucketGrid[T]{
		side:              bucketsResolution.Side(),
		bucketsResolution: bucketsResolution,
		bound:             bound,
		buckets:           make(map[ChunkKey]*lqtree.LinearQuadTree[T]),
	}
}

func (g *BucketGrid[T]) BulkInsert(entries []pow2grid.Entry[T]) {
	g.processEntriesOneByOne(
		entries,
		func(e pow2grid.Entry[T]) bool { return e.Value != nil && g.inBounds(e.Pos) },
		true,
		func(chunkKey ChunkKey, bucket *lqtree.LinearQuadTree[T], locals pow2grid.Entry[T]) {
			bucket.SingleBulkInsert(locals)
		},
	)
}

func (g *BucketGrid[T]) BulkRemove(entries []pow2grid.Entry[T]) {
	g.processEntriesOneByOne(
		entries,
		func(e pow2grid.Entry[T]) bool { return g.inBounds(e.Pos) },
		false,
		func(chunkKey ChunkKey, bucket *lqtree.LinearQuadTree[T], locals pow2grid.Entry[T]) {
			bucket.SignleBulkRemove(locals)
			if bucket.Count() == 0 {
				// delete(g.buckets, chunkKey)
			}
		},
	)
}

func (g *BucketGrid[T]) BulkMove(moves pow2grid.EntriesMove[T]) {
	if len(moves.Old) == 0 && len(moves.New) == 0 {
		return
	}

	g.BulkRemove(moves.Old)
	g.BulkInsert(moves.New)
}

func (g *BucketGrid[T]) Get(x, y uint32) (*T, bool) {
	pos := pow2grid.Pos{X: x, Y: y}
	if !g.inBounds(pos) {
		return nil, false
	}

	key, local := g.chunkKey(pos)
	bucket := g.buckets[key]
	if bucket == nil {
		return nil, false
	}
	return bucket.Get(local.X, local.Y)
}

func (g *BucketGrid[T]) QueryRange(aabb pow2grid.AABB, out []*T) int {
	if len(out) == 0 {
		return 0
	}
	clear(out)

	if len(g.buckets) == 0 {
		return 0
	}

	if !g.intersectsBound(aabb) {
		return 0
	}

	// clamp query to world bound
	aabb = g.clampToBound(aabb)

	minChunkX := aabb.Min.X / g.side
	maxChunkX := aabb.Max.X / g.side
	minChunkY := aabb.Min.Y / g.side
	maxChunkY := aabb.Max.Y / g.side

	outCounter := 0
	localAABB := pow2grid.AABB{Min: pow2grid.Pos{X: 0, Y: 0}, Max: pow2grid.Pos{X: 0, Y: 0}}

	for cx := minChunkX; cx <= maxChunkX; cx++ {
		for cy := minChunkY; cy <= maxChunkY; cy++ {
			chunkKey := ChunkKey{X: cx, Y: cy}

			bucket := g.buckets[chunkKey]
			if bucket == nil {
				continue
			}

			view := out[outCounter:]
			if len(view) == 0 {
				return outCounter
			}

			g.configureChunkAABB(chunkKey, aabb, &localAABB)

			written := bucket.QueryRange(localAABB, view)

			outCounter += written
			if outCounter >= len(out) {
				return outCounter
			}
		}
	}
	return outCounter
}

func (g *BucketGrid[T]) configureChunkAABB(chunkKey ChunkKey, oryginal pow2grid.AABB, chunkAABB *pow2grid.AABB) {
	chunkMinX := chunkKey.X * g.side
	chunkMinY := chunkKey.Y * g.side
	chunkMaxX := chunkMinX + g.side - 1
	chunkMaxY := chunkMinY + g.side - 1

	localMinX := max(oryginal.Min.X, chunkMinX)
	localMinY := max(oryginal.Min.Y, chunkMinY)

	localMaxX := min(oryginal.Max.X, chunkMaxX)
	localMaxY := min(oryginal.Max.Y, chunkMaxY)

	chunkAABB.Min = pow2grid.Pos{X: localMinX - chunkMinX, Y: localMinY - chunkMinY}
	chunkAABB.Max = pow2grid.Pos{X: localMaxX - chunkMinX, Y: localMaxY - chunkMinY}
}

func (g *BucketGrid[T]) Count() uint64 {
	return g.count
}

func (g *BucketGrid[T]) Bounds() pow2grid.AABB {
	return g.bound
}

func (g *BucketGrid[T]) ensureBucket(key ChunkKey) *lqtree.LinearQuadTree[T] {
	if bucket := g.buckets[key]; bucket != nil {
		return bucket
	}
	bucket := lqtree.NewLinearQuadTree[T](g.bucketsResolution)
	g.buckets[key] = bucket
	return bucket
}

func (g *BucketGrid[T]) chunkKey(pos pow2grid.Pos) (ChunkKey, pow2grid.Pos) {
	cx := pos.X / g.side
	cy := pos.Y / g.side

	local := pow2grid.Pos{
		X: pos.X - cx*g.side,
		Y: pos.Y - cy*g.side,
	}

	return ChunkKey{X: cx, Y: cy}, local
}

func (g *BucketGrid[T]) adjustCount(before, after uint64) {
	if after > before {
		g.count += after - before
		return
	}
	g.count -= before - after
}

func (g *BucketGrid[T]) processEntriesOneByOne(
	entries []pow2grid.Entry[T],
	keep func(pow2grid.Entry[T]) bool,
	ensureBucket bool,
	apply func(ChunkKey, *lqtree.LinearQuadTree[T], pow2grid.Entry[T]),
) {
	for _, e := range entries {
		if !keep(e) {
			continue
		}

		chunkKey, localPos := g.chunkKey(e.Pos)
		local := pow2grid.Entry[T]{Pos: localPos, Value: e.Value}

		var bucket *lqtree.LinearQuadTree[T]
		if ensureBucket {
			bucket = g.ensureBucket(chunkKey)
		} else {
			bucket = g.buckets[chunkKey]
			if bucket == nil {
				continue
			}
		}

		before := bucket.Count()
		apply(chunkKey, bucket, local)
		g.adjustCount(before, bucket.Count())
	}
}

func (g *BucketGrid[T]) inBounds(p pow2grid.Pos) bool {
	return p.X >= g.bound.Min.X && p.X <= g.bound.Max.X &&
		p.Y >= g.bound.Min.Y && p.Y <= g.bound.Max.Y
}

func (g *BucketGrid[T]) clampToBound(aabb pow2grid.AABB) pow2grid.AABB {
	if aabb.Min.X < g.bound.Min.X {
		aabb.Min.X = g.bound.Min.X
	}
	if aabb.Min.Y < g.bound.Min.Y {
		aabb.Min.Y = g.bound.Min.Y
	}
	if aabb.Max.X > g.bound.Max.X {
		aabb.Max.X = g.bound.Max.X
	}
	if aabb.Max.Y > g.bound.Max.Y {
		aabb.Max.Y = g.bound.Max.Y
	}
	return aabb
}

func (g *BucketGrid[T]) intersectsBound(aabb pow2grid.AABB) bool {
	return !(aabb.Max.X < g.bound.Min.X || aabb.Min.X > g.bound.Max.X ||
		aabb.Max.Y < g.bound.Min.Y || aabb.Min.Y > g.bound.Max.Y)
}
