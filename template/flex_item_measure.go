package template

// flexItemMeasure reads one flex item's style and measures its box: the
// natural size a widget wants, capped by its own width/height and clamped by
// min/max, with the margins, order and flex factors the style declares. The
// flex properties a container reads are here; the box the widget paints is
// the same one it would have taken in the flow. Properties the style never
// set keep their CSS initial values, which do not always match the struct's
// zero value: flex-shrink is 1, flex-basis and align-self are auto.
func (c *CSS) flexItemMeasure(role, label string) flexItem {
	st := c.style.baseStyle(role, State{})
	nw, nh := c.style.natural(role, label, st, c.style.u())
	return c.flexItemSized(st, nw, nh)
}
