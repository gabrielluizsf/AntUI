package template

import "github.com/gabrielluizsf/antui/template/css"

// MultiCol lays the widgets its callback draws out in columns: a block in the
// flow that widens to the page, styled by the multicol role, whose content is
// split into equal columns with a gutter between them and an optional rule
// drawn in the gutter. The number of columns is column-count, the width that
// fits is column-width, the gutter is column-gap and the rule is column-rule —
// the same gutter flex and grid read as their gap, so one declaration spaces
// all three.
//
//	app.MultiCol(func(m *CSS) {
//		for _, item := range items {
//			m.Label(item)
//		}
//	})
//
// A container with no height grows to the tallest of its columns and balances
// the content over them, so column-count: 2 with five widgets puts two in the
// first column and three in the second. A container that declares a height
// spends it: the flow cuts a box that runs out of room and continues it in the
// next column, and a box with break-inside: avoid is moved whole to the next
// column instead of being cut. Content that does not fit the columns the
// container asked for spills into extra columns to the right, the way a browser
// spills a paragraph past a fixed column-count.
//
// The callback runs once silently to measure every child, then once per
// fragment a child was cut into, so it must draw the same widgets in the same
// order every run. Events are read on the last run, as [CSS.Flex] does.
func (c *CSS) MultiCol(draw func(*CSS)) {
	if c.multicoling || c.flexing || c.griding || c.tabling {
		c.nestedMultiCol(draw)
		return
	}

	st := c.style.baseStyle(css.RoleMultiCol, State{})
	if st.Display == css.DisplayNone {
		return
	}
	winW, winH := c.win.Width(), c.win.Height()
	ml := c.style.length(st, st.Margin[3], winW)
	mr := c.style.length(st, st.Margin[1], winW)
	mt := c.style.length(st, st.Margin[0], winW)
	mb := c.style.length(st, st.Margin[2], winW)

	// The container is a block: it closes the inline line, collapses the
	// sibling margins and rests on the flow cursor.
	c.cursorY += c.lineH
	c.lineH = 0
	c.lineX = 0
	c.cursorY += max(c.lastMB, mt)
	y := c.cursorY

	w := max(winW-ml-mr, 0)
	if st.Has("width") {
		w = c.style.length(st, st.Width, winW)
	}
	w = c.clampDim(st, w, "width", winW)
	x := ml + (winW-ml-mr-w)/2

	b := &multiBatch{css: c, st: st}
	c.setMultiBox(b, st, x, y, w, c.multiDeclaredHeight(st), winW)
	c.multi, c.multicoling = b, true

	c.measure = true
	draw(c)
	c.measure = false
	b.flow()

	h := max(b.colH+b.top+b.bottom, c.clampDim(st, 0, "height", winH))
	if b.cuts {
		h = c.multiDeclaredHeight(st)
	}
	b.h = h
	b.contentH = max(h-b.top-b.bottom, 0)

	if c.style.boxOn(st) && visible(st) {
		c.style.paintBox(c.win, st, x, y, w, h)
	}
	c.style.paintColumnRules(c.win, b)

	for pass := range b.passes {
		b.pass = pass
		c.multiWithin(b, draw)
	}

	c.multicoling = false
	c.multi = nil
	c.cursorY = y + h
	c.lastMB = mb
}
