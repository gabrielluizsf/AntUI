package template

import "github.com/gabrielluizsf/antui/template/css"

func (b *gridBatch) placeGridItem(it *gridItem, columns, rows, columnPositions, rowPositions []int) {
	if it.colStart < 0 || it.colStart >= len(columns) || it.rowStart < 0 || it.rowStart >= len(rows) {
		return
	}
	cellX := columnPositions[it.colStart]
	cellY := rowPositions[it.rowStart]
	cellW, cellH := 0, 0
	for col := it.colStart; col < it.colEnd && col < len(columns); col++ {
		cellW += columns[col]
	}
	for row := it.rowStart; row < it.rowEnd && row < len(rows); row++ {
		cellH += rows[row]
	}
	cellW += max(it.colEnd-it.colStart-1, 0) * b.colGap
	cellH += max(it.rowEnd-it.rowStart-1, 0) * b.rowGap

	justify := b.st.JustifyItems
	if !b.st.Has("justify-items") {
		justify = css.AlignStretch
	}
	if it.justifySelf != css.AlignAuto {
		justify = it.justifySelf
	}
	align := b.st.AlignItems
	if !b.st.Has("align-items") {
		align = css.AlignStretch
	}
	if it.alignSelf != css.AlignAuto {
		align = it.alignSelf
	}

	it.w, it.h = it.nw, it.nh
	if it.st.Width.IsPct() {
		it.w = b.css.style.length(it.st, it.st.Width, cellW)
	}
	if it.st.Height.IsPct() {
		it.h = b.css.style.length(it.st, it.st.Height, cellH)
	}
	if justify == css.AlignStretch && !it.explicitW {
		it.w = max(cellW-it.ml-it.mr, 0)
	}
	if align == css.AlignStretch && !it.explicitH {
		it.h = max(cellH-it.mt-it.mb, 0)
	}
	minW, maxW := it.minW, it.maxW
	if it.st.MinWidth.IsPct() {
		minW = b.css.style.length(it.st, it.st.MinWidth, cellW)
	}
	if it.st.MaxWidth.IsPct() {
		maxW = b.css.style.length(it.st, it.st.MaxWidth, cellW)
	}
	minH, maxH := it.minH, it.maxH
	if it.st.MinHeight.IsPct() {
		minH = b.css.style.length(it.st, it.st.MinHeight, cellH)
	}
	if it.st.MaxHeight.IsPct() {
		maxH = b.css.style.length(it.st, it.st.MaxHeight, cellH)
	}
	it.w = clampGridItem(it.w, minW, maxW)
	it.h = clampGridItem(it.h, minH, maxH)

	switch justify {
	case css.AlignCenter:
		it.x = cellX + (cellW-it.ml-it.mr-it.w)/2 + it.ml
	case css.AlignFlexEnd:
		it.x = cellX + cellW - it.mr - it.w
	default:
		it.x = cellX + it.ml
	}
	switch align {
	case css.AlignCenter:
		it.y = cellY + (cellH-it.mt-it.mb-it.h)/2 + it.mt
	case css.AlignFlexEnd:
		it.y = cellY + cellH - it.mb - it.h
	default:
		it.y = cellY + it.mt
	}
}
