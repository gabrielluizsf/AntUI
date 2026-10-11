package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// transformed paints the widget behind a transform. With no transform list the
// widget draws straight into the window, as always. With one, the painting is
// recorded on a window-sized layer in its own space and then composited back
// through the transform matrix — a rotated scale costs one full-surface blit
// and leaves the untransformed neighbours underneath untouched.
//
// Three things keep the layer honest with the page. A backdrop-filter sees the
// window content that will land behind the box, planted into the layer through
// the forward transform before the widget paints. The pointer is projected
// through the inverse transform while an interactive widget paints, so hit
// tests follow the box wherever it went. And the composite walks the canvas's
// own dirty bounds, so whatever the widget drew — a select's menu below its
// box, say — comes over with it instead of being cut at the padded box.
func (cs *cssStyle) transformed(st css.Style, x, y, w, h int, paint func()) {
	if len(st.Transform) == 0 {
		paint()
		return
	}
	win := cs.win
	if cs.needsLayer() {
		l, err := canvas.NewCanvas(win.Width(), win.Height())
		if err != nil {
			paint()
			return
		}
		cs.layer = l
		cs.layerW, cs.layerH = win.Width(), win.Height()
	}
	clear(cs.layer.Pixels)
	cs.layer.ResetDirty()

	if len(st.BackdropFilters) > 0 {
		cs.plantBackdrop(st, x, y, w, h)
	}
	m := cs.transformMatrix(st, x, y, w, h)

	prev := win.SetCanvas(cs.layer)
	restore := cs.followPointer(st, m)
	paint()
	restore()
	win.SetCanvas(prev)

	region := cs.transformRegion(st, x, y, w, h)
	if d := cs.layer.Dirty; d.Width > 0 && d.Height > 0 {
		region = unionArea(region, d)
	}
	if region.Width <= 0 || region.Height <= 0 {
		return
	}
	win.Canvas().BlitMatrix(cs.layer, m, region)
}
