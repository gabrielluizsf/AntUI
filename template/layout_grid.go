package template

import "github.com/gabrielluizsf/antui/template/css"

// layoutGrid routes a widget drawn inside a Grid callback. The silent pass
// measures it and returns ok=false so nothing paints; the draw pass hands back
// the box the solver placed it in. An item the flow left out is parked by its
// placement — the area grid-area names, or the padding box — and its insets.
func (c *CSS) layoutGrid(role, label string) (x, y, w, h int, st css.Style, ok bool) {
	b := c.grid
	if c.measure {
		it := c.gridItemMeasure(role, label)
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
		return 0, 0, 0, 0, st, false
	}
	it := b.items[b.next]
	b.next++
	if it.out {
		return c.gridPlaceOut(st, it)
	}
	x = b.originX + it.x
	y = b.originY + it.y
	x, y = c.shiftInset(st, x, y, b.contentW, b.contentH)
	return x, y, it.w, it.h, st, true
}
