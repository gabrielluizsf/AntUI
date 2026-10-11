package template

import "github.com/gabrielluizsf/antui/template/css"

// gridLineShift is how many implicit lines have to stand before the explicit
// ones, for the items that count lines from the end of the axis.
func (b *gridBatch) gridLineShift(axis css.GridAxis, rows bool) int {
	shift := 0
	for i := range b.items {
		placement := b.items[i].column
		if rows {
			placement = b.items[i].row
		}
		start, end, ok := gridAxisPlacement(axis, placement)
		if ok {
			shift = min(shift, start, end)
		}
	}
	return max(-shift, 0)
}
