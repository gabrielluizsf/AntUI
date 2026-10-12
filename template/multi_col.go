package template

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
	c.multiColLayout(draw)
}
