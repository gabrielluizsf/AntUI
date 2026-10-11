package template

import "github.com/gabrielluizsf/antui/template/css"

// layoutTable routes a widget drawn inside a Table callback. A widget drawn
// inside a row is a cell of it: the silent pass measures the cell and the draw
// pass hands back the box the solver gave it. A widget drawn straight in the
// table is a row of its own, the way an anonymous row works in a browser.
func (c *CSS) layoutTable(role, label string) (x, y, w, h int, st css.Style, ok bool) {
	if c.row == nil {
		return c.tableRowOfOne(role, label)
	}
	return c.layoutCell(role, label)
}
