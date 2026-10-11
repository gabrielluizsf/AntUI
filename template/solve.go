package template

// solve turns the measured items into boxes: decide the axes, split lines,
// distribute the free space and answer where every item sits and how tall the
// content grew. It runs once per pass; every run recomputes from the measured
// nw/nh, so the draw pass agrees with the silent one.
func (b *flexBatch) solve() {
	if len(b.items) == 0 {
		return
	}
	if b.col() {
		b.solveColumn()
		return
	}
	b.solveRow()
}
