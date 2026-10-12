package template

import "github.com/gabrielluizsf/antui/template/css"

// gridPlaceFixed puts every item whose place is already known — named by
// grid-area, or pinned to both a column and a row — and marks the cells they
// occupy.
func (b *gridBatch) gridPlaceFixed(areas map[string]gridArea, columnAxis, rowAxis css.GridAxis, columnShift, rowShift int, columns, rows *[]css.GridTrack, occupied map[[2]int]bool) {
	for i := range b.items {
		it := &b.items[i]
		it.w, it.h = it.nw, it.nh
		if area, ok := areas[it.gridAreaName()]; ok {
			it.colStart, it.colEnd = area.col+columnShift, area.col+area.colSpan+columnShift
			it.rowStart, it.rowEnd = area.row+rowShift, area.row+area.rowSpan+rowShift
			it.colDefinite, it.rowDefinite = true, true
			ensureGridTracks(columns, it.colEnd)
			ensureGridTracks(rows, it.rowEnd)
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
		ensureGridTracks(columns, it.colEnd)
		ensureGridTracks(rows, it.rowEnd)
		b.occupy(occupied, it)
	}
}
