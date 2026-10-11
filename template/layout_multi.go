package template

import "github.com/gabrielluizsf/antui/template/css"

// layoutMulti routes a widget drawn inside a MultiCol callback. The silent pass
// measures it and returns ok=false so nothing paints; a draw pass hands back
// the piece the flow cut for the pass it is on.
func (c *CSS) layoutMulti(role, label string) (x, y, w, h int, st css.Style, ok bool) {
	b := c.multi
	if c.measure {
		it := c.multiItemMeasure(role, label)
		if it.st.Display == css.DisplayNone {
			return 0, 0, 0, 0, it.st, false
		}
		b.items = append(b.items, it)
		return 0, 0, 0, 0, it.st, false
	}

	e := c.style.beginWidget(role, label)
	st = c.style.entryStyle(e, role, e.state)
	if st.Display == css.DisplayNone {
		return 0, 0, 0, 0, st, false
	}
	if b.next >= len(b.items) {
		return 0, 0, 0, 0, st, true
	}
	it := &b.items[b.next]
	b.next++
	piece := b.pieceAt(it)
	// A child of a column is a block in it: it takes the column's width less
	// its own margins, and a piece starts inside its top margin.
	x = b.originX + piece.col*(b.colW+b.gap) + it.ml
	y = b.originY + piece.y
	w = max(b.colW-it.ml-it.mr, 0)
	return x, y, w, piece.h, st, true
}
