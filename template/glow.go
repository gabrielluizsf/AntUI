package template

import "github.com/gabrielluizsf/antui/canvas"

// glow draws translucent shells around a rounded rectangle, the strongest
// closest to it, spreading colour outward like light. u is the scale, so the
// shells stay proportioned to the control on every resolution.
func glow(cv *canvas.Canvas, x, y, w, h, radius int, c canvas.Color, strength, u int) {
	for i := 2; i >= 0; i-- {
		alpha := strength / (3 - i)
		if alpha <= 0 {
			continue
		}
		g := (i + 1) * u
		cv.FillRoundRect(x-g, y-g, w+g*2, h+g*2, radius+g, canvas.Fade(c, alpha))
	}
}
