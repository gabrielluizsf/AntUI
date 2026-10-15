package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// run draws one styled run with no alignment, for widgets that place their own
// text (switches, calendar cells).
func (cs *cssStyle) run(cv *canvas.Canvas, st css.Style, x, y int, text string) {
	o := cs.textOpts(st)
	o.Color = cs.ink(st)
	cv.DrawStyled(canvas.Point{X: float64(x), Y: float64(y)}, transformText(text, st.TextTransform), o)
}
