package qtree

import (
	"sort"

	"github.com/kjkrol/gokg/pkg/geometry"
)

func sortItems[T geometry.SupportedNumeric](items []Item[T]) {
	sort.Slice(items, func(i, j int) bool {
		ai, aj := items[i].Bound(), items[j].Bound()
		first, _ := geometry.SortBoxesBy(
			ai, aj,
			func(box geometry.BoundingBox[T]) T { return box.TopLeft.Y },
			func(box geometry.BoundingBox[T]) T { return box.TopLeft.X },
			func(box geometry.BoundingBox[T]) T { return box.BottomRight.Y },
			func(box geometry.BoundingBox[T]) T { return box.BottomRight.X },
		)
		return first.Equals(ai)
	})
}
