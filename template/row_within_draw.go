package template

import "github.com/gabrielluizsf/antui/template/css"

// rowWithinDraw hands one row its cells back in the draw pass.
func (c *CSS) rowWithinDraw(r *tableRow, role, label string) (x, y, w, h int, st css.Style, ok bool) {
	if r == nil {
		return 0, 0, 0, 0, st, false
	}
	c.rowWithin(c.table, r, func(*CSS) {
		x, y, w, h, st, ok = c.layoutCell(role, label)
	})
	return x, y, w, h, st, ok
}
