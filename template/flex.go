package template

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
	c.flexLayout(draw)
}
