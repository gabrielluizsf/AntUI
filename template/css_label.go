package template

import "github.com/gabrielluizsf/antui/template/css"

// Label paints one line of text in the window, placed and styled by CSS.
// The widget never reports an event.
func (c *CSS) Label(text string) {
	x, y, w, h, st, ok := c.layout(css.RoleLabel, text)
	if !ok || !visible(st) {
		return
	}
	c.style.transformed(st, x, y, w, h, func() {
		c.style.Label(c.win, x, y, text)
	})
}
