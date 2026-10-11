package template

import "github.com/gabrielluizsf/antui/template/css"

// placeCaption gives a caption the table's whole width and the run above the
// rows, with the height it asks for: a caption is not in a column, so it is
// never given one.
func (c *CSS) placeCaption(cell *tableCell) (x, y, w, h int, st css.Style, ok bool) {
	b := c.table
	e := c.style.beginWidget(cell.role, cell.label)
	st = c.style.entryStyle(e, cell.role, e.state)
	if st.Display == css.DisplayNone {
		return 0, 0, 0, 0, st, false
	}
	return b.originX, b.originY - cell.nh, max(cell.nw, b.contentW), cell.nh, st, true
}
