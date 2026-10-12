package template

import "github.com/gabrielluizsf/antui/template/css"

// placeDenseRow scans the whole grid row by row for the first cell the item
// fits in, which is what grid-auto-flow: dense asks for.
func (b *gridBatch) placeDenseRow(occupied map[[2]int]bool, columns, rows *[]css.GridTrack, it *gridItem, colSpan, rowSpan int) bool {
	for row := 0; row < len(*rows); row++ {
		for col := 0; col < len(*columns); col++ {
			if b.fits(occupied, col, row, colSpan, rowSpan) {
				b.reserve(occupied, it, col, row, colSpan, rowSpan)
				return true
			}
		}
	}
	return false
}
