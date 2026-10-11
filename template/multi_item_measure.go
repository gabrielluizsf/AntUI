package template

// multiItemMeasure reads one child's style and measures its box: the natural
// size a widget wants, with the margins the style declares. The width and
// height it ends up with are resolved by the flow, which knows the column it
// will be drawn in.
func (c *CSS) multiItemMeasure(role, label string) multiItem {
	st := c.style.baseStyle(role, State{})
	nw, nh := c.style.natural(role, label, st, c.style.u())
	return c.multiItemSized(st, nw, nh)
}
