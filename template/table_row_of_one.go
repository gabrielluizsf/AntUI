package template

import "github.com/gabrielluizsf/antui/template/css"

// tableRowOfOne wraps a widget drawn straight in a table in the row it belongs
// to: a row of a single cell, which is what a browser's anonymous boxes make of
// a child that did not ask to be a row. A widget that asked to be a caption is
// not wrapped — it is painted above the table, across its whole width.
func (c *CSS) tableRowOfOne(role, label string) (x, y, w, h int, st css.Style, ok bool) {
	b := c.table
	if c.measure {
		style := c.style.baseStyle(role, State{})
		if style.Caption() {
			cell := c.tableCellMeasure("", "", style)
			b.children = append(b.children, tableChild{caption: cell})
			return 0, 0, 0, 0, style, false
		}
		r := &tableRow{}
		c.rowWithin(b, r, func(*CSS) { c.layoutCell(role, label) })
		b.children = append(b.children, tableChild{row: r})
		return 0, 0, 0, 0, style, false
	}
	if b.rowIdx >= len(b.children) {
		return 0, 0, 0, 0, css.Style{}, false
	}
	child := b.children[b.rowIdx]
	b.rowIdx++
	if child.caption != nil {
		return c.placeCaption(child.caption)
	}
	return c.rowWithinDraw(child.row, role, label)
}
