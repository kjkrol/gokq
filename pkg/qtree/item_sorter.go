package qtree

import (
	"sort"

	"github.com/kjkrol/gokg/pkg/geom"
)

func sortItems[T geom.Numeric](items []Item[T]) {
	sort.Slice(items, func(i, j int) bool {
		ai, aj := items[i].Bound(), items[j].Bound()
		first, _ := geom.SortAABBsBy(
			ai, aj,
			func(box geom.AABB[T]) T { return box.TopLeft.Y },
			func(box geom.AABB[T]) T { return box.TopLeft.X },
			func(box geom.AABB[T]) T { return box.BottomRight.Y },
			func(box geom.AABB[T]) T { return box.BottomRight.X },
		)
		return first.Equals(ai)
	})
}
