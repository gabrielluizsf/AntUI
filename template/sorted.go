package template

import "sort"

// sorted returns the item indices in flex order: order ascending, ties keep
// the draw order.
func (b *flexBatch) sorted() []int {
	idx := make([]int, len(b.items))
	for i := range idx {
		idx[i] = i
	}
	sort.SliceStable(idx, func(i, j int) bool {
		return b.items[idx[i]].order < b.items[idx[j]].order
	})
	return idx
}
