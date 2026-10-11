package svg

import "github.com/gabrielluizsf/antui/canvas"

// fade is a colour at an opacity, as a multiply of the two the node was given.
func fade(c canvas.Color, opacity float64) canvas.Color {
	if opacity >= 1 {
		return c
	}
	a := int(float64(c.A())*clampFloat(opacity, 0, 1) + 0.5)
	return canvas.RGBA(c.R(), c.G(), c.B(), uint8(a))
}
