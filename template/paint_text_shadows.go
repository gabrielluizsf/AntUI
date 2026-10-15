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
			o.Color = col
			cv.DrawStyled(canvas.Point{X: float64(tx + dx), Y: float64(ty + dy)}, text, o)
			continue
		}
		pad := blur + 2
		layer, err := canvas.NewCanvas(cv.TextWidthStyled(text, o)+2*pad, o.Height()+2*pad)
		if err != nil {
			continue
		}
		o.Color = col | 0xFF000000
		layer.DrawStyled(canvas.Point{X: float64(pad), Y: float64(pad)}, text, o)
		layer.Blur(canvas.Area{X: 0, Y: 0, Width: layer.Width, Height: layer.Height}, max(blur/2, 1))
		canvas.ScaleAlpha(layer, float64(col.A())/255)
		cv.BlitOver(tx+dx-pad, ty+dy-pad, layer)
	}
}
