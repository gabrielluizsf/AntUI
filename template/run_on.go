package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// runOn draws one styled run in a given colour, for a cell that paints its
// text against a fill of its own choosing (a calendar's selected day).
func (cs *cssStyle) runOn(cv *canvas.Canvas, st css.Style, colour canvas.Color, x, y int, text string) {
	o := cs.textOpts(st)
	cv.DrawStyled(x, y, transformText(text, st.TextTransform), colour, o)
}
