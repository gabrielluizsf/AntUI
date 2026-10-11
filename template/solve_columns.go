package template

import "github.com/gabrielluizsf/antui/template/css"

// solveColumns reads the count and the width apart. An explicit column-count
// says how many; a column-width says how many of them fit, and the smaller of
// the two is what the container makes. With neither there is one column and
// the block flows as it would have without any of this. The width is measured
// the way every other length in the sheet is — through the template's length,
// so it answers to the window's unit like the gutter beside it.
func (b *multiBatch) solveColumns(st css.Style) {
	n := 0
	if st.ColumnCount > 0 {
		n = st.ColumnCount
	}
	if st.Has("column-width") && gridLengthDefinite(st.ColumnWidth) {
		w := b.css.style.length(st, st.ColumnWidth, b.contentW)
		if w > 0 {
			if fit := (b.contentW + b.gap) / (w + b.gap); n == 0 || fit < n {
				n = fit
			}
		}
	}
	if n < 1 {
		n = 1
	}
	b.count = n
	if n == 1 {
		b.colW = b.contentW
		return
	}
	b.colW = max((b.contentW-b.gap*(n-1))/n, 0)
}
