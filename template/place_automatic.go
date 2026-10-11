package template

import "github.com/gabrielluizsf/antui/template/css"

func (b *gridBatch) placeAutomatic(it *gridItem, columns *[]css.GridTrack, rows *[]css.GridTrack, occupied map[[2]int]bool, rowCursor, colCursor *int) {
	colSpan, rowSpan := gridPlacementSpan(it.column), gridPlacementSpan(it.row)
	if it.colDefinite {
		colSpan = max(it.colEnd-it.colStart, colSpan)
	}
	if it.rowDefinite {
		rowSpan = max(it.rowEnd-it.rowStart, rowSpan)
	}

	if b.st.GridAutoFlow == css.GridAutoFlowColumn {
		b.placeColumnAutomatic(it, columns, rows, occupied, rowCursor, colCursor, colSpan, rowSpan)
		return
	}
	b.placeRowAutomatic(it, columns, rows, occupied, rowCursor, colCursor, colSpan, rowSpan)
}
