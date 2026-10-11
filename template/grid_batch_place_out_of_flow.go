package template

// placeOutOfFlow positions an item the flow left out, in coordinates from the
// content origin like every other item. The box it is measured against is the
// area its placement names, and the padding box when the placement names none;
// the insets then move it inside that box, and a pair of opposite insets with
// no width of its own fills what is left between them.
func (b *gridBatch) placeOutOfFlow(it *gridItem, columns, rows, columnPositions, rowPositions []int) {
	baseX, baseW := 0, b.contentW
	if it.colDefinite {
		baseX, baseW = gridAreaBox(columns, columnPositions, it.colStart, it.colEnd, b.colGap)
	}
	baseY, baseH := 0, b.contentH
	if it.rowDefinite {
		baseY, baseH = gridAreaBox(rows, rowPositions, it.rowStart, it.rowEnd, b.rowGap)
	}

	it.w, it.h = it.nw, it.nh
	// A pair of opposite insets on an axis whose size was not declared fills
	// what is left between them, so the width is settled before the position.
	if !it.explicitW && it.inLeft >= 0 && it.inRight >= 0 {
		it.w = max(baseW-it.inLeft-it.inRight, 0)
	}
	if !it.explicitH && it.inTop >= 0 && it.inBottom >= 0 {
		it.h = max(baseH-it.inTop-it.inBottom, 0)
	}
	it.x, it.y = baseX, baseY
	if it.inLeft >= 0 {
		it.x = baseX + it.inLeft
	} else if it.inRight >= 0 {
		it.x = baseX + baseW - it.inRight - it.w
	}
	if it.inTop >= 0 {
		it.y = baseY + it.inTop
	} else if it.inBottom >= 0 {
		it.y = baseY + baseH - it.inBottom - it.h
	}
	if it.minW >= 0 {
		it.w = max(it.w, it.minW)
	}
	if it.maxW >= 0 {
		it.w = min(it.w, it.maxW)
	}
	if it.minH >= 0 {
		it.h = max(it.h, it.minH)
	}
	if it.maxH >= 0 {
		it.h = min(it.h, it.maxH)
	}
}
