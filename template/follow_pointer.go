package template

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// followPointer makes an interactive transformed box answer the pointer where
// it really is. For the duration of the layer pass the mouse is projected
// through the inverse transform, so the widget's own hit tests run in its
// local space; the map is exact — m⁻¹(p) is inside the box precisely when p
// is inside m(box) — so no bounding-box approximation is involved. The
// returned closure restores the pointer the backend reported.
func (cs *cssStyle) followPointer(st css.Style, m canvas.Matrix) func() {
	win := cs.win
	mx, my := win.MouseX(), win.MouseY()
	if !interactive(st) {
		return func() {}
	}
	inv, ok := m.Inverse()
	if !ok {
		return func() {}
	}
	lx, ly := inv.Map(float64(mx)+0.5, float64(my)+0.5)
	px, py := int(math.Floor(lx)), int(math.Floor(ly))
	if px == mx && py == my {
		return func() {}
	}
	win.SetMouse(px, py)
	return func() { win.SetMouse(mx, my) }
}
