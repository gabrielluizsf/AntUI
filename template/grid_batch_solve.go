package template

import "github.com/gabrielluizsf/antui/template/css"

// solve lays every item out: it resolves both axes, places each item on them,
// then hands the axes their space and turns the placements into boxes.
func (b *gridBatch) solve() {
	// A container that declares its height bounds the content box up front, so
	// the row axis has a definite size to divide and a definite space for the
	// rows an auto-repeat counts into.
	rowDefinite := b.st.Has("height") && gridLengthDefinite(b.st.Height)
	rowBase := b.contentW
	if rowDefinite {
		rowBase = b.contentH
	}
	columnAxis := b.resolveGridAxis(b.st.GridTemplateColumns, b.contentW, b.colGap)
	rowAxis := b.resolveGridAxis(b.st.GridTemplateRows, b.gridRowSpace(rowDefinite), b.rowGap)
	columns := columnAxis.Tracks
	rows := rowAxis.Tracks
	areaColumns, areaRows := 0, 0
	for _, row := range b.st.GridTemplateAreas {
		areaColumns = max(areaColumns, len(row))
	}
	areaRows = len(b.st.GridTemplateAreas)
	ensureGridTracks(&columns, areaColumns)
	ensureGridTracks(&rows, areaRows)
	columnShift := b.gridLineShift(columnAxis, false)
	rowShift := b.gridLineShift(rowAxis, true)
	prependGridTracks(&columns, columnShift)
	prependGridTracks(&rows, rowShift)
	areas := b.gridAreas()
	occupied := make(map[[2]int]bool, len(b.items)*2)

	for i := range b.items {
		it := &b.items[i]
		it.w, it.h = it.nw, it.nh
		if area, ok := areas[it.gridAreaName()]; ok {
			it.colStart, it.colEnd = area.col+columnShift, area.col+area.colSpan+columnShift
			it.rowStart, it.rowEnd = area.row+rowShift, area.row+area.rowSpan+rowShift
			it.colDefinite, it.rowDefinite = true, true
			ensureGridTracks(&columns, it.colEnd)
			ensureGridTracks(&rows, it.rowEnd)
			if !it.out {
				b.occupy(occupied, it)
			}
			continue
		}
		it.colStart, it.colEnd, it.colDefinite = gridAxisPlacement(columnAxis, it.column)
		it.colStart += columnShift
		it.colEnd += columnShift
		it.rowStart, it.rowEnd, it.rowDefinite = gridAxisPlacement(rowAxis, it.row)
		it.rowStart += rowShift
		it.rowEnd += rowShift
		if it.out || !(it.colDefinite && it.rowDefinite) {
			continue
		}
		ensureGridTracks(&columns, it.colEnd)
		ensureGridTracks(&rows, it.rowEnd)
		b.occupy(occupied, it)
	}

	order := b.gridOrder()
	rowCursor, colCursor := 0, 0
	place := func(index int) {
		it := &b.items[index]
		if !it.out && !(it.colDefinite && it.rowDefinite) {
			b.placeAutomatic(it, &columns, &rows, occupied, &rowCursor, &colCursor)
		}
	}
	// The items locked to a line come first, in the order they were added, and
	// everything else follows in that same order: an item with a definite column
	// and an automatic one shares the pass, so it takes the cursor where the
	// items before it left it instead of jumping ahead of them.
	if b.st.GridAutoFlow == css.GridAutoFlowRow {
		for _, index := range order {
			it := &b.items[index]
			if !it.out && it.rowDefinite && !it.colDefinite {
				place(index)
			}
		}
		for _, index := range order {
			it := &b.items[index]
			if !it.out && !(it.rowDefinite && !it.colDefinite) {
				place(index)
			}
		}
	} else {
		for _, index := range order {
			it := &b.items[index]
			if !it.out && it.colDefinite && !it.rowDefinite {
				place(index)
			}
		}
		for _, index := range order {
			it := &b.items[index]
			if !it.out && !(it.colDefinite && !it.rowDefinite) {
				place(index)
			}
		}
	}

	rowGapTotal := max(len(rows)-1, 0) * b.rowGap
	colGapTotal := max(len(columns)-1, 0) * b.colGap
	columnSizes := b.resolveGridTracks(columns, true, true, max(b.contentW-colGapTotal, 0))
	rowSizes := b.resolveGridTracks(rows, false, rowDefinite, max(rowBase-rowGapTotal, 0))
	if !rowDefinite {
		b.contentH = sumGridSizes(rowSizes) + rowGapTotal
	}
	b.trackColumns = columns

	columnPositions := gridTrackPositions(columnSizes, b.colGap, sumGridSizes(columnSizes), b.contentW, b.st.GridJustifyContent)
	rowPositions := gridTrackPositions(rowSizes, b.rowGap, sumGridSizes(rowSizes), b.contentH, b.st.AlignContent)
	for i := range b.items {
		if b.items[i].out {
			b.placeOutOfFlow(&b.items[i], columnSizes, rowSizes, columnPositions, rowPositions)
			continue
		}
		b.placeGridItem(&b.items[i], columnSizes, rowSizes, columnPositions, rowPositions)
	}
}
