package template

import "github.com/gabrielluizsf/antui/template/css"

// layoutFlex routes a widget drawn inside a Flex callback. The silent pass
// measures it and returns ok=false so nothing paints; the draw pass hands
// back the box the solver placed it in. Out-of-flow children are parked by
// their insets against the container's content box and skip distribution.
func (c *CSS) layoutFlex(role, label string) (x, y, w, h int, st css.Style, ok bool) {
	b := c.flex
	if c.measure {
		it := c.flexItemMeasure(role, label)
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
		return c.flexPlaceOut(st, it)
	}
	return b.originX + it.x, b.originY + it.y, it.w, it.h, st, true
}
