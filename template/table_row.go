package template

// tableRow is one row: the cells drawn inside it, and the height the tallest
// one gives it.
type tableRow struct {
	cells  []*tableCell
	height int
	next   int // the cell the draw pass is handing back
}
