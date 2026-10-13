package template

import "github.com/gabrielluizsf/antui/canvas"

// halo is a shape that light spills out of: the rounded rectangle it spills
// from, the colour of the light, how far it reaches, and the unit the shells are
// proportioned to — so the halo keeps its shape and its weight on every
// resolution rather than being a blur drawn at one fixed size.
type halo struct {
	x, y, w, h int
	radius     int
	colour     canvas.Color
	strength   int
	u          int
}

// glow draws translucent shells around a rounded rectangle, the strongest
// closest to it, spreading colour outward like light.
func glow(cv *canvas.Canvas, h halo) {
	for i := 2; i >= 0; i-- {
		alpha := h.strength / (3 - i)
		if alpha <= 0 {
			continue
		}
		g := (i + 1) * h.u
		cv.FillRoundRect(canvas.Round(canvas.Area{X: h.x - g, Y: h.y - g, Width: h.w + g*2, Height: h.h + g*2}, h.radius+g), canvas.Fade(h.colour, alpha))
	}
}
