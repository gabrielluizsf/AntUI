package template

import "sort"

func (b *gridBatch) gridOrder() []int {
	order := make([]int, len(b.items))
	for i := range order {
		order[i] = i
	}
	sort.SliceStable(order, func(i, j int) bool {
		return b.items[order[i]].order < b.items[order[j]].order
	})
	return order
}
