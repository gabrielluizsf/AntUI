package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// paintBackground paints the surface of a box: the background colour first,
// then each layer above it — the last declared image sitting lowest, the
// first on top, the CSS stacking. A layer is rasterized into its own canvas
// and composited over the box so the alphas of several stacked layers aculate
// instead of overwriting one another, and each is cut at its background-clip
// box so a border-box layer does not bleed over the border.
func (cs *cssStyle) paintBackground(win *antui.Window, st css.Style, x, y, w, h, rx, ry int) {
	cv := win.Canvas()
	if st.Has("background-color") || st.Has("background") {
		cv.FillRoundRect(canvas.Rounded{Area: canvas.Area{X: x, Y: y, Width: w, Height: h}, RX: rx, RY: ry}, cs.fill(st))
	}
	n := len(st.BackgroundImages)
	if n == 0 {
		return
	}
	// Rounded corners are fractions of the box side, so a layer clipped to a
	// smaller box keeps the same curvature.
	cr := 0.0
	if min(w, h) > 0 {
		cr = max(float64(rx), float64(ry)) / float64(min(w, h))
	}
	for i := n - 1; i >= 0; i-- {
		layer := cs.paintBackgroundLayer(win, st, i, x, y, w, h, cr)
		if layer != nil {
			cv.BlitOver(x, y, layer)
		}
	}
}
