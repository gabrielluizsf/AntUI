package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// Background is the colour the window is cleared to. It comes from the body
// of the stylesheet when one is loaded, and from the theme otherwise.
func (c *CSS) Background(win *antui.Window) canvas.Color {
	return c.style.Background(win)
}
