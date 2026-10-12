package template

import "github.com/gabrielluizsf/antui/template/css"

// placeCursorRow walks the flow cursor forward until the item fits, growing
// the tracks as it goes.
func (b *gridBatch) placeCursorRow(occupied map[[2]int]bool, columns, rows *[]css.GridTrack, it *gridItem, rowCursor, colCursor *int, colSpan, rowSpan int) {
	row, col := *rowCursor, *colCursor
	for {
		if col > 0 && col+colSpan > len(*columns) {
			col = 0
			row++
			*rowCursor, *colCursor = row, col
			continue
		}
		ensureGridTracks(columns, col+colSpan)
		ensureGridTracks(rows, row+rowSpan)
		if b.fits(occupied, col, row, colSpan, rowSpan) {
			b.reserve(occupied, it, col, row, colSpan, rowSpan)
			*rowCursor, *colCursor = row, col+colSpan
			return
		}
		col++
	}
}
