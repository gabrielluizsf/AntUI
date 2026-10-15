package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// Progress paints a read-only progress bar with progress from 0 to 1. When
// the stylesheet rounds the box into a full circle, it paints a ring instead,
// and the optional label floats in its middle. It never reports an event: a
// bar has nothing to say.
func (c *CSS) Progress(progress float32, center ...string) {
	x, y, w, h, st, ok := c.layout(css.RoleProgress, "")
	if !ok {
		return
	}
	label := ""
	if len(center) > 0 {
		label = center[0]
	}
	c.paint(st, css.RoleProgress, "", x, y, w, h, func() {
		c.ui.Progress(c.win, canvas.Area{X: x, Y: y, Width: w, Height: h}, progress)
		if label != "" {
			tw := c.style.measure(st, label)
			th := c.style.textHeight(st)
			c.style.run(c.win.Canvas(), st, x+(w-tw)/2, y+(h-th)/2, label)
		}
	})
}
