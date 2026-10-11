package template

import "github.com/gabrielluizsf/antui/template/css"

// nestedMultiCol is [CSS.MultiCol] called from inside another container: the
// container becomes an item of the parent, placed like any other and holding
// its own children.
func (c *CSS) nestedMultiCol(draw func(*CSS)) {
	st := c.style.baseStyle(css.RoleMultiCol, State{})
	if st.Display == css.DisplayNone {
		return
	}
	switch {
	case c.multicoling:
		c.multiInMulti(st, draw)
	case c.flexing:
		c.flexNestedMultiCol(st, draw)
	case c.griding:
		c.gridNestedMultiCol(st, draw)
	default:
		// Straight in a table: the content of the cell in place, or of a row
		// of its own when nothing put one there.
		c.multiInTableCell(st, draw)
	}
}
