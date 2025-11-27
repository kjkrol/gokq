package lqtree

type ChunkKey struct {
	X, Y uint32
}

type ZOrderBucketGrid[T any] struct {
	bucketMaxXY  uint32
	bucketstSize LQTSize
	bound        AABB
	buckets      map[ChunkKey]*LinearQuadTree[T]
	count        uint64
}

var _ SpatialIndex[any] = (*ZOrderBucketGrid[any])(nil)

// Choose bucketstSize (and thus bucketMaxXY) to roughly match typical query AABB.
// Example: if most queries are <32x32, a 64–128 bucket keeps them in 1–4 buckets with moderate memory;
// QueryRange that spans many buckets will pay for many map lookups + per-bucket scans.
func NewZOrderBucketGrid[T any](bucketstSize LQTSize, bound AABB) *ZOrderBucketGrid[T] {
	size := bucketstSize.Resolution() + 1
	return &ZOrderBucketGrid[T]{
		bucketMaxXY:  size,
		bucketstSize: bucketstSize,
		bound:        bound,
		buckets:      make(map[ChunkKey]*LinearQuadTree[T]),
	}
}

func (g *ZOrderBucketGrid[T]) BulkInsert(entries []Entry[T]) {
	if len(entries) == 0 {
		return
	}

	perBucket := make(map[ChunkKey][]Entry[T])

	for i := range entries {
		entry := entries[i]
		if entry.Value == nil || !g.inBounds(entry.Pos) {
			continue
		}

		key, local := g.chunkKey(entry.Pos)
		perBucket[key] = append(perBucket[key], Entry[T]{Pos: local, Value: entry.Value})
	}

	for key, bucketEntries := range perBucket {
		bucket := g.ensureBucket(key)
		before := bucket.Count()
		bucket.BulkInsert(bucketEntries)
		g.adjustCount(before, bucket.Count())
	}
}

func (g *ZOrderBucketGrid[T]) BulkRemove(positions []Pos) {
	if len(positions) == 0 {
		return
	}

	perBucket := make(map[ChunkKey][]Pos)

	for i := range positions {
		if !g.inBounds(positions[i]) {
			continue
		}

		key, local := g.chunkKey(positions[i])
		if g.buckets[key] == nil {
			continue
		}
		perBucket[key] = append(perBucket[key], local)
	}

	for key, bucketPositions := range perBucket {
		bucket := g.buckets[key]
		before := bucket.Count()
		bucket.BulkRemove(bucketPositions)
		g.adjustCount(before, bucket.Count())
	}
}

func (g *ZOrderBucketGrid[T]) BulkMove(moves []EntryMove[T]) {
	for i := range moves {
		move := moves[i]

		srcKey, srcLocal := g.chunkKey(move.Old)
		dstKey, dstLocal := g.chunkKey(move.New)

		value := move.Value
		srcBucket := g.buckets[srcKey]

		if value == nil && srcBucket != nil {
			if v, ok := srcBucket.Get(srcLocal.X, srcLocal.Y); ok {
				value = v
			}
		}

		if srcBucket != nil {
			before := srcBucket.Count()
			srcBucket.BulkRemove([]Pos{srcLocal})
			g.adjustCount(before, srcBucket.Count())
		}

		if value == nil {
			continue
		}

		if !g.inBounds(move.New) {
			continue
		}

		dstBucket := g.ensureBucket(dstKey)
		before := dstBucket.Count()
		dstBucket.BulkInsert([]Entry[T]{{Pos: dstLocal, Value: value}})
		g.adjustCount(before, dstBucket.Count())
	}
}

func (g *ZOrderBucketGrid[T]) Get(x, y uint32) (*T, bool) {
	pos := Pos{X: x, Y: y}
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

func (g *ZOrderBucketGrid[T]) QueryRange(aabb AABB) []*T {
	if len(g.buckets) == 0 {
		return nil
	}

	if !g.intersectsBound(aabb) {
		return nil
	}

	// clamp query to world bound
	aabb = g.clampToBound(aabb)

	minChunkX := aabb.Min.X / g.bucketMaxXY
	maxChunkX := aabb.Max.X / g.bucketMaxXY
	minChunkY := aabb.Min.Y / g.bucketMaxXY
	maxChunkY := aabb.Max.Y / g.bucketMaxXY

	results := make([]*T, 0)

	for cx := minChunkX; cx <= maxChunkX; cx++ {
		for cy := minChunkY; cy <= maxChunkY; cy++ {
			key := ChunkKey{X: cx, Y: cy}
			bucket := g.buckets[key]
			if bucket == nil {
				continue
			}

			chunkMinX := cx * g.bucketMaxXY
			chunkMinY := cy * g.bucketMaxXY
			chunkMaxX := chunkMinX + g.bucketMaxXY - 1
			chunkMaxY := chunkMinY + g.bucketMaxXY - 1

			localMinX := max(aabb.Min.X, chunkMinX)
			localMinY := max(aabb.Min.Y, chunkMinY)

			localMaxX := min(aabb.Max.X, chunkMaxX)
			localMaxY := min(aabb.Max.Y, chunkMaxY)

			localAABB := AABB{
				Min: Pos{X: localMinX - chunkMinX, Y: localMinY - chunkMinY},
				Max: Pos{X: localMaxX - chunkMinX, Y: localMaxY - chunkMinY},
			}

			results = append(results, bucket.QueryRange(localAABB)...)
		}
	}

	return results
}

func (g *ZOrderBucketGrid[T]) Count() uint64 {
	return g.count
}

func (g *ZOrderBucketGrid[T]) Bounds() AABB {
	return g.bound
}

func (g *ZOrderBucketGrid[T]) ensureBucket(key ChunkKey) *LinearQuadTree[T] {
	if bucket := g.buckets[key]; bucket != nil {
		return bucket
	}
	bucket := NewLinearQuadTree[T](g.bucketstSize)
	g.buckets[key] = bucket
	return bucket
}

func (g *ZOrderBucketGrid[T]) chunkKey(pos Pos) (ChunkKey, Pos) {
	cx := pos.X / g.bucketMaxXY
	cy := pos.Y / g.bucketMaxXY

	local := Pos{
		X: pos.X - cx*g.bucketMaxXY,
		Y: pos.Y - cy*g.bucketMaxXY,
	}

	return ChunkKey{X: cx, Y: cy}, local
}

func (g *ZOrderBucketGrid[T]) adjustCount(before, after uint64) {
	if after > before {
		g.count += after - before
		return
	}
	g.count -= before - after
}

func (g *ZOrderBucketGrid[T]) inBounds(p Pos) bool {
	return p.X >= g.bound.Min.X && p.X <= g.bound.Max.X &&
		p.Y >= g.bound.Min.Y && p.Y <= g.bound.Max.Y
}

func (g *ZOrderBucketGrid[T]) clampToBound(aabb AABB) AABB {
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

func (g *ZOrderBucketGrid[T]) intersectsBound(aabb AABB) bool {
	return !(aabb.Max.X < g.bound.Min.X || aabb.Min.X > g.bound.Max.X ||
		aabb.Max.Y < g.bound.Min.Y || aabb.Min.Y > g.bound.Max.Y)
}
