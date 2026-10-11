package template

// lineCross measures one line's natural cross size: the tallest (or widest)
// item plus its cross margins.
func (b *flexBatch) lineCross(line []int) int {
	cross := 0
	for _, idx := range line {
		it := &b.items[idx]
		size, margin := it.nh, it.mt+it.mb
		if b.col() {
			size, margin = it.nw, it.ml+it.mr
		}
		cross = max(cross, size+margin)
	}
	return cross
}
