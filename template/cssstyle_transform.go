package template

import (
	"math"

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

// unionArea is the bounding box of two areas; an empty one is absorbed.
func unionArea(a, b canvas.Area) canvas.Area {
	if b.Width <= 0 || b.Height <= 0 {
		return a
	}
	if a.Width <= 0 || a.Height <= 0 {
		return b
	}
	x1 := max(a.X+a.Width, b.X+b.Width)
	y1 := max(a.Y+a.Height, b.Y+b.Height)
	a.X = min(a.X, b.X)
	a.Y = min(a.Y, b.Y)
	a.Width, a.Height = x1-a.X, y1-a.Y
	return a
}

// needsLayer reports whether the scratch layer has to be rebuilt: a window
// that grew or shrank since the last frame would otherwise draw into a stale
// surface.
func (cs *cssStyle) needsLayer() bool {
	w, h := cs.win.Width(), cs.win.Height()
	return cs.layer == nil || cs.layerW != w || cs.layerH != h
}

// transformMatrix is the whole transform list pinned to the box and its
// origin: the user's functions run around the origin point, which itself sits
// at the box's top-left plus transform-origin. An undeclared origin is the
// CSS initial "50% 50%", the centre.
func (cs *cssStyle) transformMatrix(st css.Style, x, y, w, h int) canvas.Matrix {
	ox, oy := cs.originPoint(st, x, y, w, h)

	m := canvas.Translate(ox, oy)
	for i := range st.Transform {
		f := st.Transform[i]
		switch f.Kind {
		case css.TransformTranslate:
			m = m.Mul(canvas.Translate(float64(cs.length(st, f.Dx, w)), float64(cs.length(st, f.Dy, h))))
		case css.TransformScale:
			m = m.Mul(canvas.Scale(f.Sx, f.Sy))
		case css.TransformRotate:
			m = m.Mul(canvas.Rotate(f.Ax.Rad()))
		case css.TransformSkew:
			m = m.Mul(canvas.Skew(f.Ax.Rad(), f.Ay.Rad()))
		case css.TransformMatrix:
			m = m.Mul(canvas.Matrix{A: f.M[0], B: f.M[1], C: f.M[2], D: f.M[3], E: f.M[4], F: f.M[5]})
		}
	}
	return m.Mul(canvas.Translate(-ox, -oy))
}

// originPoint is the pivot of the transform in device pixels: the box's
// top-left plus transform-origin, the initial value landing dead centre.
func (cs *cssStyle) originPoint(st css.Style, x, y, w, h int) (float64, float64) {
	if !st.Has("transform-origin") {
		return float64(x + w/2), float64(y + h/2)
	}
	return float64(x + cs.length(st, st.TransformOrigin[0], w)),
		float64(y + cs.length(st, st.TransformOrigin[1], h))
}

// transformRegion is the part of the recorded layer worth compositing: the
// box widened by everything a box can shed beyond its edge — drops, spread,
// blur, outlines. The caller widens it further to what the widget actually
// drew, so content past the pad, like a select's menu, still comes over.
func (cs *cssStyle) transformRegion(st css.Style, x, y, w, h int) canvas.Area {
	pad := cs.transformPad(st)
	return canvas.Area{X: x - pad, Y: y - pad, Width: w + 2*pad, Height: h + 2*pad}
}

// transformPad is how far, in device pixels, a widget's pixels can reach past
// its border box: shadows and outlines. Whatever is closer to the box is
// inside the region already.
func (cs *cssStyle) transformPad(st css.Style) int {
	u := cs.u()
	pad := 0
	span := func(s css.Shadow) int {
		ex := max(s.X, -s.X)
		if s.Y > ex {
			ex = s.Y
		}
		if -s.Y > ex {
			ex = -s.Y
		}
		return ex + s.Spread + s.Blur + 2
	}
	for i := range st.Filters {
		if f := st.Filters[i]; f.Drop != nil {
			pad = max(pad, span(*f.Drop))
		} else if f.Kind == canvas.FilterBlur {
			pad = max(pad, int(f.Amount)+1)
		}
	}
	for i := range st.BoxShadow {
		if s := st.BoxShadow[i]; !s.Inset {
			pad = max(pad, span(s))
		}
	}
	for i := range st.TextShadow {
		pad = max(pad, span(st.TextShadow[i]))
	}
	if cs.outlineOn(st) {
		pad = max(pad, st.OutlineOffset+st.OutlineWidth)
	}
	if pad < 0 {
		return 0
	}
	return pad * u
}
