package apk

import (
	_ "image/jpeg"

	"github.com/gabrielluizsf/antui/canvas"
)

// rounded takes the corners off a square icon.
func rounded(src *canvas.Canvas) *canvas.Canvas {
	out, err := canvas.NewCanvas(src.Width, src.Height)
	if err != nil {
		return src
	}
	clear(out)
	// The radius is half the width, and the comparison is in squared
	// distance so that there is no square root per pixel.
	radius := src.Width / 2
	limit := radius * radius
	for y := range src.Height {
		dy := y - radius
		for x := range src.Width {
			dx := x - radius
			if dx*dx+dy*dy <= limit {
				out.Put(x, y, src.At(x, y))
			}
		}
	}
	return out
}
