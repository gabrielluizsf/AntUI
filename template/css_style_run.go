package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// run draws one styled run with no alignment, for widgets that place their own
// text (switches, calendar cells).
func (cs *cssStyle) run(cv *canvas.Canvas, st css.Style, x, y int, text string) {
	o := cs.textOpts(st)
	cv.DrawStyled(x, y, transformText(text, st.TextTransform), cs.ink(st), o)
}
