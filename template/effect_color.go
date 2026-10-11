package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// effectColor turns a computed effect colour into pixels: a zero colour is the
// theme ink, and every colour is faded by the element's opacity.
func (cs *cssStyle) effectColor(win *antui.Window, st css.Style, c canvas.Color) canvas.Color {
	if c == 0 {
		c = win.Theme().Text
	}
	return canvas.Fade(c, cs.alpha(st))
}
