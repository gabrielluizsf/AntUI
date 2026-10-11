package template

// split wraps the ordered items into lines along the main axis, packing until
// the next item would cross the edge. Items too big for a line own it alone.
func (b *flexBatch) split(order []int, avail, gap int, wrap bool) [][]int {
	lines := [][]int{}
	cur := []int{}
	used := 0
	for _, idx := range order {
		it := &b.items[idx]
		om := it.base + it.ml + it.mr
		if wrap && len(cur) > 0 && used+om+gap > avail {
			lines = append(lines, cur)
			cur = []int{}
			used = 0
		}
		if len(cur) > 0 {
			used += gap
		}
		used += om
		cur = append(cur, idx)
	}
	return append(lines, cur)
}
