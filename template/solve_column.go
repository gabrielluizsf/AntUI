package template

import "github.com/gabrielluizsf/antui/template/css"

// solveColumn lays the items along a vertical main axis, stacking wrapped
// groups next to each other. Without a declared height the axis is unbounded:
// nothing grows, shrinks or leans on justify-content, the column is a plain
// vertical stack of its natural heights.
func (b *flexBatch) solveColumn() {
	definite := b.st.Has("height")
	avail := b.contentH
	for i := range b.items {
		it := &b.items[i]
		it.base = b.baseMain(it, avail, definite)
	}
	order := b.sorted()
	wrap := definite && b.st.FlexWrap != css.FlexWrapNowrap
	mainGap := b.gapPx(b.st.RowGap)

	total := 0
	prev := -1
	for _, idx := range order {
		it := &b.items[idx]
		if prev >= 0 {
			total += mainGap
		}
		total += it.base + it.mt + it.mb
		prev = idx
	}

	eff := avail
	if !definite || eff < total {
		eff = total
	}
	lines := b.split(order, eff, mainGap, wrap)

	for _, line := range lines {
		b.distribute(line, eff, mainGap)
		b.justify(line, eff, mainGap)
	}

	crosses := make([]int, len(lines))
	for li, line := range lines {
		crosses[li] = b.lineCross(line)
	}
	offsets, _ := b.stackLines(lines, crosses, b.gapPx(b.st.ColumnGap), b.contentW, true)

	for li, line := range lines {
		b.alignCross(line, crosses[li], offsets[li])
	}
	if definite {
		b.contentH = avail
	} else {
		b.contentH = total
	}
}
