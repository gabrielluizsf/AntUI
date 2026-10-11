package template

import "github.com/gabrielluizsf/antui/template/css"

// Flex lays the widgets its callback draws out as flex items inside one flex
// container: a block in the flow that widens to the page, styled by the flex
// role in the stylesheet, and growing to hold its children along the
// flex-direction, wrapped and justified per the container's flex properties.
// Child widgets are placed by the solver, not the top-to-bottom flow, and
// keep all their interaction — a button inside a row still clicks.
//
// The callback runs twice per frame: first silently, to measure every child's
// natural box, then again to draw and interact at the boxes the algorithm
// assigned. It must draw the same widgets in the same order both runs — which
// a stable loop over a data slice does — and an event a child raises is seen
// on the second run, so read it before the callback returns:
//
//	clicked := false
//	app.Flex(func(f *CSS) {
//		for _, label := range items {
//			if f.Button(label).Is(event.Button, event.Click) {
//				clicked = true
//			}
//		}
//	})
//	if clicked { /* the click really happened */ }
//
// A Flex drawn inside another Flex's callback flattens into the outer
// container: the nested block's own face is dropped and its children join
// the outer items, which keeps the two-pass measurement honest. A grid drawn
// inside a flex is the other way round — it stays a box of its own and becomes
// a flex item, while a Flex inside a grid flattens into the grid.
func (c *CSS) Flex(draw func(*CSS)) {
	if c.flexing || c.griding {
		draw(c)
		return
	}

	st := c.style.baseStyle(css.RoleFlex, State{})
	if st.Display == css.DisplayNone {
		return
	}
	u := c.style.u()
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

	bl := st.BorderWidth[3] * u
	br := st.BorderWidth[1] * u
	bt := st.BorderWidth[0] * u
	bb := st.BorderWidth[2] * u
	pl := c.style.length(st, st.Padding[3], winW)
	pr := c.style.length(st, st.Padding[1], winW)
	pt := c.style.length(st, st.Padding[0], winW)
	pb := c.style.length(st, st.Padding[2], winW)

	b := &flexBatch{
		css:      c,
		st:       st,
		x:        x,
		y:        y,
		w:        w,
		originX:  x + bl + pl,
		contentW: max(w-br-bl-pr-pl, 0),
	}
	c.flex = b
	c.flexing = true

	c.measure = true
	draw(c)
	c.measure = false

	// A container that declares its height bounds the content box up front,
	// so the solver has a definite main or cross size to distribute against.
	if st.Has("height") {
		hh := c.clampDim(st, c.style.length(st, st.Height, winH), "height", winH)
		b.contentH = max(hh-bt-bb-pt-pb, 0)
	}
	b.solve()

	h := b.contentH + bt + bb + pt + pb
	if st.Has("height") {
		h = c.clampDim(st, c.style.length(st, st.Height, winH), "height", winH)
	}
	h = max(h, c.clampDim(st, 0, "height", winH))
	b.h = h
	// The content box every child draws against is exactly what the border box
	// leaves inside, so an explicit height or a tall row agrees with the flow.
	b.contentH = max(h-bt-bb-pt-pb, 0)
	b.originY = y + bt + pt

	if c.style.boxOn(st) && visible(st) {
		c.style.paintBox(c.win, st, x, y, w, h)
	}

	b.next = 0
	draw(c)

	c.flexing = false
	c.flex = nil

	c.cursorY = y + h
	c.lastMB = mb
}
