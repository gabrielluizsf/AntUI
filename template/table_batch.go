package template

import "github.com/gabrielluizsf/antui/template/css"

// tableBatch is one table: its box, the columns the solver cut out of it and
// the rows the callback drew. The silent pass collects the cells, the solver
// gives every cell the width of its column and every row the height of its
// tallest cell, and the draw pass hands the boxes back.
type tableBatch struct {
	css      *CSS
	st       css.Style
	x, y, w  int // container border box
	h        int
	originX  int // content-box top-left in window coordinates
	originY  int
	top      int // the border and padding above the rows
	bottom   int // the border and padding below them
	caption  int // the height a caption takes above the rows
	contentW int
	contentH int
	cols     []int // the width each column was given
	children []tableChild
	row      *tableRow // the row collecting cells, nil outside one
	rowIdx   int       // the row the draw pass is on
}
