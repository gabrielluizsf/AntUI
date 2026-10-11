package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// plantBackdrop seeds the recording layer with what the window shows behind
// the box, pulled forward through the transform: layer point p shows window
// at m(p). A backdrop-filter then has the real page underneath to blur
// instead of the empty layer. A translucent element still blends its own
// backdrop over what it already covered — the backdrop is doped, not
// replaced — a corner opaque backdrops never reach.
func (cs *cssStyle) plantBackdrop(st css.Style, x, y, w, h int) {
	m := cs.transformMatrix(st, x, y, w, h)
	cs.layer.BlitBackdrop(cs.win.Canvas(), m, canvas.Area{X: x, Y: y, Width: w, Height: h})
}
