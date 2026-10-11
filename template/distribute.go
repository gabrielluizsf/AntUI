package template

// distribute runs the flex-grow and flex-shrink algorithm over one line:
// grow shares the leftover space proportionally, shrink takes space back from
// the fattest bases, both floors clamped by an explicit minimum.
func (b *flexBatch) distribute(line []int, avail, gap int) {
	used, grow := 0, 0.0
	shrinkK := 0.0
	for _, idx := range line {
		it := &b.items[idx]
		used += it.base + it.ml + it.mr
		grow += it.grow
		shrinkK += float64(it.base) * it.shrink
	}
	free := avail - gap*(len(line)-1) - used
	switch {
	case free > 0 && grow > 0:
		for _, idx := range line {
			it := &b.items[idx]
			extra := int(float64(free) * (it.grow / grow))
			size := it.base + extra
			if it.maxMain >= 0 {
				size = min(size, it.maxMain)
			}
			b.setMain(it, size)
		}
	case free < 0 && shrinkK > 0:
		want := -free
		for _, idx := range line {
			it := &b.items[idx]
			cut := int(float64(want) * (float64(it.base) * it.shrink / shrinkK))
			b.setMain(it, max(it.minMain, it.base-cut))
		}
	default:
		for _, idx := range line {
			it := &b.items[idx]
			b.setMain(it, it.base)
		}
	}
}
