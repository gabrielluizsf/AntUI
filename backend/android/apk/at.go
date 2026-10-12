package apk

import (
	"image/color"
	_ "image/jpeg"
)

func (c canvasImage) At(x, y int) color.Color {
	p := c.cv.At(x, y)
	return color.NRGBA{R: p.R(), G: p.G(), B: p.B(), A: p.A()}
}
