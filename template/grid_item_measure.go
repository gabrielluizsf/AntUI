package template

func (c *CSS) gridItemMeasure(role, label string) gridItem {
	st := c.style.baseStyle(role, State{})
	u := c.style.u()
	naturalW, naturalH := c.style.natural(role, label, st, u)
	// A track asking for min-content wants the narrowest the item can be, and
	// the size the item declares or constrains is the narrowest it will ever be.
	minNaturalW, minNaturalH := c.style.minContent(role, label, st, u)
	return c.gridItemSized(role, label, st, naturalW, naturalH, minNaturalW, minNaturalH)
}
