package template

import "github.com/gabrielluizsf/antui/template/css"

// tableCellMeasure reads one cell's style and measures the box it asks for:
// the natural size of the widget, capped by a width or height the cell's own
// style declares and clamped by its min and max. A percentage answers to the
// widest the table can be, which is the room the flow has for it — the table's
// own width is decided afterwards, by the columns.
func (c *CSS) tableCellMeasure(role, label string, st css.Style) *tableCell {
	cell := &tableCell{role: role, label: label, st: st}
	cell.nw, cell.nh = c.style.natural(role, label, st, c.style.u())
	winW := c.win.Width()
	if st.Has("width") {
		cell.nw = c.style.length(st, st.Width, winW)
	}
	if st.Has("max-width") {
		cell.nw = min(cell.nw, c.style.length(st, st.MaxWidth, winW))
	}
	if st.Has("min-width") {
		cell.nw = max(cell.nw, c.style.length(st, st.MinWidth, winW))
	}
	winH := c.win.Height()
	if st.Has("height") {
		cell.nh = c.style.length(st, st.Height, winH)
	}
	if st.Has("max-height") {
		cell.nh = min(cell.nh, c.style.length(st, st.MaxHeight, winH))
	}
	if st.Has("min-height") {
		cell.nh = max(cell.nh, c.style.length(st, st.MinHeight, winH))
	}
	return cell
}
