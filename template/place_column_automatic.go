package template

import "github.com/gabrielluizsf/antui/template/css"

func (b *gridBatch) placeColumnAutomatic(it *gridItem, columns *[]css.GridTrack, rows *[]css.GridTrack, occupied map[[2]int]bool, rowCursor, colCursor *int, colSpan, rowSpan int) {
	if it.rowDefinite {
		startCol := 0
		if !b.st.GridAutoFlowDense && it.rowStart == *rowCursor {
			startCol = *colCursor
		}
		for col := startCol; ; col++ {
			ensureGridTracks(columns, col+colSpan)
			if b.fits(occupied, col, it.rowStart, colSpan, rowSpan) {
				b.reserve(occupied, it, col, it.rowStart, colSpan, rowSpan)
				*rowCursor, *colCursor = it.rowStart, col+colSpan
				return
			}
		}
	}
	if it.colDefinite {
		startRow := 0
		if !b.st.GridAutoFlowDense {
			startRow = *rowCursor
			if it.colStart < *colCursor {
				startRow++
			}
		}
		for row := startRow; ; row++ {
			ensureGridTracks(rows, row+rowSpan)
			if b.fits(occupied, it.colStart, row, colSpan, rowSpan) {
				b.reserve(occupied, it, it.colStart, row, colSpan, rowSpan)
				*rowCursor, *colCursor = row+rowSpan, it.colStart
				return
			}
		}
	}
	if b.st.GridAutoFlowDense {
		for col := 0; col < len(*columns); col++ {
			for row := 0; row < len(*rows); row++ {
				if b.fits(occupied, col, row, colSpan, rowSpan) {
					b.reserve(occupied, it, col, row, colSpan, rowSpan)
					return
				}
			}
		}
	}
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
