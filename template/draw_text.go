package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/css"
)

// drawText paints a text run inside a box, aligned by text-align. Justifying
// belongs to wrapped paragraphs, so a lone run is never stretched.
func (cs *cssStyle) drawText(win *antui.Window, st css.Style, x, y, w, h int, text string) {
	cs.drawTextLine(win, st, x, y, w, h, text, false)
}
