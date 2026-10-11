package template

// placeInline packs an inline or inline-block box along the current line,
// wrapping to a new line when it would cross the window edge. Vertical
// margins of inline boxes are ignored, as they are in CSS.
func (c *CSS) placeInline(w, h, winW, ml, mr int) (int, int) {
	if c.lineH > 0 && c.lineX+ml+w+mr > winW {
		c.cursorY += c.lineH
		c.lineH = 0
		c.lineX = 0
	}
	x := c.lineX + ml
	y := c.cursorY
	c.lineX += ml + w + mr
	c.lineH = max(c.lineH, h)
	c.lastMB = 0
	return x, y
}
