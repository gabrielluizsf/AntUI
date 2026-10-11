package template

import "github.com/gabrielluizsf/antui/template/css"

// nestedTable is [CSS.Table] called from inside another container: the table
// becomes an item of the parent, placed like any other and holding its own
// rows.
func (c *CSS) nestedTable(draw func(*CSS)) {
	st := c.style.baseStyle(css.RoleTable, State{})
	if st.Display == css.DisplayNone {
		return
	}
	switch {
	case c.multicoling:
		c.tableInColumns(st, draw)
	case c.flexing:
		c.flexNestedTable(st, draw)
	case c.griding:
		c.gridNestedTable(st, draw)
	case c.row != nil:
		c.tableInCell(st, draw)
	default:
		// Straight in a table's callback: a row of its own, holding the table.
		c.tableRowOfOneBlock(func(*CSS) { c.tableInCell(st, draw) })
	}
}
