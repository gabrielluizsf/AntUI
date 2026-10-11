package template

// placeBlock starts a new line for a block box: it closes whatever inline line
// was open, collapses the vertical margins and centres the box.
func (c *CSS) placeBlock(w, h, winW, ml, mr, mt, mb int) (int, int) {
	c.cursorY += c.lineH
	c.lineH = 0
	c.lineX = 0
	c.cursorY += max(c.lastMB, mt)
	x := ml + (winW-ml-mr-w)/2
	y := c.cursorY
	c.cursorY += h
	c.lastMB = mb
	return x, y
}
