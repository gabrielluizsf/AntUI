package template

import "github.com/gabrielluizsf/antui/template/css"

// solveRow lays the items along a horizontal main axis: wrap, grow, justify,
// then cross-align each line and stack the lines with align-content.
func (b *flexBatch) solveRow() {
	for i := range b.items {
		it := &b.items[i]
		it.base = b.baseMain(it, b.contentW, true)
	}
	order := b.sorted()
	wrap := b.st.FlexWrap != css.FlexWrapNowrap
	mainGap := b.gapPx(b.st.ColumnGap)
	lines := b.split(order, b.contentW, mainGap, wrap)

	for _, line := range lines {
		b.distribute(line, b.contentW, mainGap)
		b.justify(line, b.contentW, mainGap)
	}

	crosses := make([]int, len(lines))
	for li, line := range lines {
		crosses[li] = b.lineCross(line)
	}

	definite := b.st.Has("height")
	avail := 0
	if definite {
		avail = b.contentH
	}
	offsets, content := b.stackLines(lines, crosses, b.gapPx(b.st.RowGap), avail, definite)
	b.contentH = content

	for li, line := range lines {
		b.alignCross(line, crosses[li], offsets[li])
	}
}
