package template

// Reset starts a new "page": the vertical cursor returns to the top of the
// window, the widget timeline begins a fresh run, and the widgets drawn
// afterwards flow down from it again. Call it once per frame before drawing.
//
// It is also where the loading screen takes the frame, when the program is
// still loading: Reset runs once at the start of every frame of the
// coord-free loop, which is the one place such a loop passes through, so
// the screen can draw above the frame's widgets from there without the
// developer wiring anything.
func (c *CSS) Reset() {
	c.style.resetTimeline()
	c.cursorY = 0
	c.lineX = 0
	c.lineH = 0
	c.lastMB = 0
	c.flex = nil
	c.flexing = false
	c.grid = nil
	c.griding = false
	c.multi = nil
	c.multicoling = false
	c.table = nil
	c.tabling = false
	c.row = nil
	c.measure = false
	c.ui.ld.arm(c.win, c.ui)
}
