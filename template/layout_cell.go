package template

import "github.com/gabrielluizsf/antui/template/css"

// layoutCell places one cell of the row in place. The silent pass measures it
// and reports ok=false so nothing paints; the draw pass hands back the box the
// solver gave it.
func (c *CSS) layoutCell(role, label string) (x, y, w, h int, st css.Style, ok bool) {
	r := c.row
	b := c.table
	if c.measure {
		cell := c.tableCellMeasure(role, label, c.style.baseStyle(role, State{}))
		r.cells = append(r.cells, cell)
		return 0, 0, 0, 0, cell.st, false
	}
	e := c.style.beginWidget(role, label)
	st = c.style.entryStyle(e, role, e.state)
	if st.Display == css.DisplayNone {
		return 0, 0, 0, 0, st, false
	}
	if r.next >= len(r.cells) {
		return 0, 0, 0, 0, st, false
	}
	cell := r.cells[r.next]
	r.next++
	x, y = b.originX+cell.x, b.originY+cell.y
	return x, y, cell.w, cell.h, st, true
}
