package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// paintScrollbars draws the indicators an overflow: scroll or auto box gets
// when its content is larger than the box. They are visual for now; the
// wheel and drag behaviour arrives with the scroll model.
func (cs *cssStyle) paintScrollbars(win *antui.Window, st css.Style, x, y, w, h, contentW, contentH int) {
	u := cs.u()
	bar := max(3, 3*u)
	inset := max(1, u)
	cv := win.Canvas()
	theme := win.Theme()
	track := canvas.Fade(theme.Surface, trackAlpha)
	thumb := canvas.Shade(theme.Surface, -0.4)

	vScroll := st.Overflow[1] == css.OverflowScroll || st.Overflow[1] == css.OverflowAuto
	if vScroll && contentH > h && h > 2*inset {
		tx := x + w - bar - inset
		cv.FillRoundRect(canvas.Round(canvas.Area{X: tx, Y: y + inset, Width: bar, Height: h - 2*inset}, bar/2), track)
		th := max(h*(h-2*inset)/max(contentH, 1), 3*bar)
		th = min(th, h-2*inset)
		cv.FillRoundRect(canvas.Round(canvas.Area{X: tx, Y: y + inset, Width: bar, Height: th}, bar/2), thumb)
	}
	hScroll := st.Overflow[0] == css.OverflowScroll || st.Overflow[0] == css.OverflowAuto
	if hScroll && contentW > w && w > 2*inset {
		ty := y + h - bar - inset
		cv.FillRoundRect(canvas.Round(canvas.Area{X: x + inset, Y: ty, Width: w - 2*inset, Height: bar}, bar/2), track)
		tw := max(w*(w-2*inset)/max(contentW, 1), 3*bar)
		tw = min(tw, w-2*inset)
		cv.FillRoundRect(canvas.Round(canvas.Area{X: x + inset, Y: ty, Width: tw, Height: bar}, bar/2), thumb)
	}
}
