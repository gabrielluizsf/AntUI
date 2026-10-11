package template

import "github.com/gabrielluizsf/antui/template/css"

// Grid lays the widgets its callback draws out as grid items inside one grid
// container: a block in the flow that widens to the page, styled by the grid
// role in the stylesheet, and split into tracks per the container's
// grid-template properties. Child widgets are placed by the solver, not the
// top-to-bottom flow, and keep all their interaction — a button in a cell still
// clicks.
//
// The callback runs twice per frame: first silently, to measure every child's
// natural box, then again to draw and interact at the boxes the algorithm
// assigned. It must draw the same widgets in the same order both runs — which
// a stable loop over a data slice does — and an event a child raises is seen on
// the second run, so read it before the callback returns:
//
//	selected := 0
//	app.Grid(func(f *CSS) {
//		for i, cell := range cells {
//			if f.Button(cell.label).Is(event.Button, event.Click) {
//				selected = i
//			}
//		}
//	})
//
// Tracks are the fixed lengths, auto, fr and minmax() the grammar allows, plus
// repeat() — with a count, or with auto-fill and auto-fit to work out the count
// from the room the container has — and the content keywords min-content,
// max-content and fit-content(). Lines may be named in a template and placed
// by name, in either direction, and grid-template-areas gives a name to a
// rectangle of cells an item can ask for by grid-area.
//
// A grid drawn inside another grid's callback, or inside a flex callback, is an
// item of that container: the outer solver picks the cell it lands in, the
// container keeps its own face, and its children are laid out inside the box it
// was given, so a percentage or an fr track in it answers to the cell it ended
// up in. Containers carry no classes, so a nested grid is styled by the same
// grid rule as the one around it — a width declared there is a width on both.
//
//	subgrid is not supported: a nested grid lays its own tracks out from
//	scratch, and asking it to adopt the parent's lines is an error.
func (c *CSS) Grid(draw func(*CSS)) {
	if c.griding {
		c.nestedGrid(draw)
		return
	}
	if c.flexing {
		c.flexNestedGrid(draw)
		return
	}

	st := c.style.baseStyle(css.RoleGrid, State{})
	if st.Display == css.DisplayNone {
		return
	}
	winW, winH := c.win.Width(), c.win.Height()
	ml := c.style.length(st, st.Margin[3], winW)
	mr := c.style.length(st, st.Margin[1], winW)
	mt := c.style.length(st, st.Margin[0], winW)
	mb := c.style.length(st, st.Margin[2], winW)

	c.cursorY += c.lineH
	c.lineH = 0
	c.lineX = 0
	c.cursorY += max(c.lastMB, mt)
	y := c.cursorY

	w := max(winW-ml-mr, 0)
	if st.Has("width") && gridLengthDefinite(st.Width) {
		w = c.style.length(st, st.Width, winW)
	}
	w = c.clampDim(st, w, "width", winW)
	x := ml + (winW-ml-mr-w)/2

	b := &gridBatch{css: c, st: st}
	c.setGridBox(b, st, x, y, w, c.gridDeclaredHeight(st), winW)
	c.grid, c.griding = b, true
	c.measure = true
	draw(c)
	c.measure = false
	b.solve()

	_, _, top, bottom := c.gridChrome(st, winW)
	if st.Has("min-height") && gridLengthDefinite(st.MinHeight) {
		limit := c.clampDim(st, c.style.length(st, st.MinHeight, winH), "height", winH)
		b.contentH = max(b.contentH, max(limit-top-bottom, 0))
	}
	if st.Has("max-height") && gridLengthDefinite(st.MaxHeight) {
		limit := c.clampDim(st, c.style.length(st, st.MaxHeight, winH), "height", winH)
		b.contentH = min(b.contentH, max(limit-top-bottom, 0))
	}
	h := b.contentH + top + bottom
	if st.Has("height") && gridLengthDefinite(st.Height) {
		h = c.style.length(st, st.Height, winH)
	}
	h = c.clampDim(st, h, "height", winH)
	c.setGridBox(b, st, x, y, w, h, winW)

	if c.style.boxOn(st) && visible(st) {
		c.style.paintBox(c.win, st, x, y, w, h)
	}

	c.gridWithin(b, draw)

	c.grid = nil
	c.griding = false
	c.cursorY = y + h
	c.lastMB = mb
}
