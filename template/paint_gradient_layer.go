package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// paintGradientLayer rasterizes one gradient into a fresh box-sized canvas
// and clips it. The gradient family, shape and size constants of the engine
// coincide with the canvas's own, so the wire is a straight cast.
func (cs *cssStyle) paintGradientLayer(st css.Style, i int, g css.Gradient, w, h int, cr float64) *canvas.Canvas {
	u := cs.u()
	origin := propAt(st.BackgroundOrigin, i, css.BackPadding)
	ox, oy, ow, oh := cs.originRect(st, 0, 0, w, h, origin)
	if ow <= 0 || oh <= 0 {
		return nil
	}
	cg := canvas.Gradient{
		Kind:    canvas.GradientKind(uint8(g.Kind)),
		Angle:   g.Angle,
		CenterX: cs.centerFrac(st, g.Center[0], u, ow),
		CenterY: cs.centerFrac(st, g.Center[1], u, oh),
		Shape:   canvas.GradientShape(uint8(g.Shape)),
		Size:    canvas.GradientSize(uint8(g.Size)),
		Stops:   make([]canvas.GradientStop, 0, len(g.Stops)),
	}
	alpha := cs.alpha(st)
	for _, s := range g.Stops {
		a := int(s.Color.A()) * alpha / 255
		cg.Stops = append(cg.Stops, canvas.GradientStop{Offset: s.Offset, Color: canvas.Fade(s.Color, a)})
	}
	layer := boxLayer(w, h)
	if layer == nil {
		return nil
	}
	layer.FillGradient(canvas.Area{X: ox, Y: oy, Width: ow, Height: oh}, cg)
	cs.maskLayer(layer, st, propAt(st.BackgroundClip, i, css.BackBorder), 0, 0, w, h, cr)
	return layer
}
