package template

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

	b.gridPlaceFixed(areas, columnAxis, rowAxis, columnShift, rowShift, &columns, &rows, occupied)
	b.gridAutoPlace(&columns, &rows, occupied)

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
