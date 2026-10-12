package template

import "github.com/gabrielluizsf/antui/template/css"

// placeDenseColumn scans the whole grid column by column for the first cell
// the item fits in, which is what grid-auto-flow: column dense asks for.
func (b *gridBatch) placeDenseColumn(occupied map[[2]int]bool, columns, rows *[]css.GridTrack, it *gridItem, colSpan, rowSpan int) bool {
	for col := 0; col < len(*columns); col++ {
		for row := 0; row < len(*rows); row++ {
			if b.fits(occupied, col, row, colSpan, rowSpan) {
				b.reserve(occupied, it, col, row, colSpan, rowSpan)
				return true
			}
		}
	}
	return false
}
