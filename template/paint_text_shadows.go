package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// paintTextShadows draws the text-shadow list under a text run. A shadow with
// no blur is the glyphs painted again, offset; a blurred one is rendered to a
// small layer so the glyph edges can soften.
func (cs *cssStyle) paintTextShadows(win *antui.Window, st css.Style, tx, ty int, text string, o canvas.TextStyle) {
	if text == "" {
		return
	}
	cv := win.Canvas()
	u := cs.u()
	for i := range st.TextShadow {
		sh := st.TextShadow[i]
		col := cs.effectColor(win, st, sh.Color)
		dx, dy := sh.X*u, sh.Y*u
		blur := sh.Blur * u
		if blur <= 0 {
			cv.DrawStyled(tx+dx, ty+dy, text, col, o)
			continue
		}
		pad := blur + 2
		layer, err := canvas.NewCanvas(cv.TextWidthStyled(text, o)+2*pad, o.Height()+2*pad)
		if err != nil {
			continue
		}
		layer.DrawStyled(pad, pad, text, col|0xFF000000, o)
		layer.Blur(0, 0, layer.Width, layer.Height, max(blur/2, 1))
		canvas.ScaleAlpha(layer, float64(col.A())/255)
		cv.BlitOver(tx+dx-pad, ty+dy-pad, layer)
	}
}
