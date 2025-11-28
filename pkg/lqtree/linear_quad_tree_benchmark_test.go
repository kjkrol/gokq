package lqtree

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/kjkrol/gokq/pkg/pow2grid"
)

const (
	qtBenchResolution = pow2grid.Size1024x1024
	qtBenchMaxXY      = (uint32(1) << qtBenchResolution) - 1
)

func BenchmarkLinearQuadTreeQueryRange(b *testing.B) {
	treeResolution := qtBenchResolution
	aabbSize := pow2grid.Size16x16

	cases := []struct {
		name         string
		totalEntries int
		querySize    uint32
	}{
		// name: [entries]-[tree size]-[aabb size]
		{"1k-" + treeResolution.String() + "-" + aabbSize.String(), 1000, aabbSize.Side()},
		{"10k-" + treeResolution.String() + "-" + aabbSize.String(), 10000, aabbSize.Side()},
		{"100k-" + treeResolution.String() + "-" + aabbSize.String(), 100000, aabbSize.Side()},
		{"500k-" + treeResolution.String() + "-" + aabbSize.String(), 500000, aabbSize.Side()},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			benchmarkLinearQuadTreeQueryRange(b, tc.totalEntries, tc.querySize)
		})
	}
}

func benchmarkLinearQuadTreeQueryRange(b *testing.B, totalEntries int, querySize uint32) {
	src := rand.New(rand.NewSource(3))

	entries := make([]pow2grid.Entry[string], totalEntries)
	for i := range totalEntries {
		entries[i] = pow2grid.Entry[string]{
			Pos:   pow2grid.Pos{X: randCoord(src, qtBenchMaxXY), Y: randCoord(src, qtBenchMaxXY)},
			Value: strPtr(fmt.Sprintf("v%d", i)),
		}
	}

	tree := NewLinearQuadTree[string](qtBenchResolution)
	tree.BulkInsert(entries)

	query := pow2grid.AABB{
		Min: pow2grid.Pos{X: 100, Y: 100},
		Max: pow2grid.Pos{X: 100 + querySize, Y: 100 + querySize},
	}

	results := make([]*string, querySize*querySize)

	for i := 0; b.Loop(); i++ {
		_ = tree.QueryRange(query, results)
	}
}

func randCoord(r *rand.Rand, max uint32) uint32 {
	return uint32(r.Intn(int(max + 1)))
}
