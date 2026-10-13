package template

import "github.com/gabrielluizsf/antui/canvas"

// fadeColor applies the element opacity without discarding the colour's own
// alpha. Effects such as rgba() outlines and shadows carry meaningful alpha
// before the stylesheet opacity is applied.
func fadeColor(c canvas.Color, opacity int) canvas.Color {
	return canvas.Fade(c, int(c.A())*opacity/255)
}
