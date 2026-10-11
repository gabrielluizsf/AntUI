package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/css"
)

// paintFrame is paintBox without the surface: the border, the shadows and the
// outline of a box whose background is not the box's to paint, because the
// widget inside it draws a surface of its own over a smaller part of it.
func (cs *cssStyle) paintFrame(win *antui.Window, st css.Style, x, y, w, h int) {
	cs.paintBoxOf(win, st, x, y, w, h, false)
}
