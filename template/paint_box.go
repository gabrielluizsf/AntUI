package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/css"
)

// paintBox fills a rounded box with the background and strokes it with the
// border, honouring border-box widths, an elliptical radius and the per-side
// border styles. Around that it layers the visual effects: a backdrop-filter
// sees the pixels behind the box, the outer shadows fall behind it, the inset
// shadows fall on the surface, an outline rings it, and a filter sees the
// box once it is painted. A missing background leaves the window's own fill
// showing through.
func (cs *cssStyle) paintBox(win *antui.Window, st css.Style, x, y, w, h int) {
	cs.paintBoxOf(win, st, x, y, w, h, true)
}
