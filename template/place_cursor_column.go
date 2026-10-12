package template

import "github.com/gabrielluizsf/antui/template/css"

// placeCursorColumn walks the flow cursor forward until the item fits, growing
// the tracks as it goes.
func (b *gridBatch) placeCursorColumn(occupied map[[2]int]bool, columns, rows *[]css.GridTrack, it *gridItem, rowCursor, colCursor *int, colSpan, rowSpan int) {
	col, row := *colCursor, *rowCursor
	for {
		if row > 0 && row+rowSpan > len(*rows) {
			row = 0
			col++
			*rowCursor, *colCursor = row, col
			continue
		}
		ensureGridTracks(columns, col+colSpan)
		ensureGridTracks(rows, row+rowSpan)
		if b.fits(occupied, col, row, colSpan, rowSpan) {
			b.reserve(occupied, it, col, row, colSpan, rowSpan)
			*rowCursor, *colCursor = row+rowSpan, col
			return
		}
		row++
	}
}
