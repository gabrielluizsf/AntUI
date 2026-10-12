package template

import "github.com/gabrielluizsf/antui/template/css"

// gridItemBox resolves the item's used width and height inside the cell,
// stretching it when the alignment asks and no explicit size was given, then
// clamps it to the min and max the cell resolved percentages give.
func (b *gridBatch) gridItemBox(it *gridItem, cellW, cellH int, justify, align uint8) {
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
}
