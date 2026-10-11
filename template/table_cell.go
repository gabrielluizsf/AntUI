package template

import "github.com/gabrielluizsf/antui/template/css"

// tableCell is one cell of a row: the box it wants, the box the row gave it,
// and the container it holds when the cell's content is a table or a column
// block of its own.
type tableCell struct {
	role, label string
	st          css.Style
	nw, nh      int // measured border-box size
	w, h        int // solved by the row
	x, y        int // solved, relative to the table's content origin
	multi       *multiNested
	table       *tableNested
}
