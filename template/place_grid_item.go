package template

import "github.com/gabrielluizsf/antui/template/css"

func (b *gridBatch) placeGridItem(it *gridItem, columns, rows, columnPositions, rowPositions []int) {
	if it.colStart < 0 || it.colStart >= len(columns) || it.rowStart < 0 || it.rowStart >= len(rows) {
		return
	}
	cellX := columnPositions[it.colStart]
	cellY := rowPositions[it.rowStart]
	cellW, cellH := b.gridCellSize(it, columns, rows)
	justify, align := b.gridItemAlignment(it)
	b.gridItemBox(it, cellW, cellH, justify, align)

	switch justify {
	case css.AlignCenter:
		it.x = cellX + (cellW-it.ml-it.mr-it.w)/2 + it.ml
	case css.AlignFlexEnd:
		it.x = cellX + cellW - it.mr - it.w
	default:
		it.x = cellX + it.ml
	}
	switch align {
	case css.AlignCenter:
		it.y = cellY + (cellH-it.mt-it.mb-it.h)/2 + it.mt
	case css.AlignFlexEnd:
		it.y = cellY + cellH - it.mb - it.h
	default:
		it.y = cellY + it.mt
	}
}
