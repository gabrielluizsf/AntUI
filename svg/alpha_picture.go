package svg

import "github.com/gabrielluizsf/antui/canvas"

// alphaPicture is a picture of only where there is something — the colour
// taken off, keeping the shape a filter can put its own colour into, which is
// what SourceAlpha names.
func alphaPicture(src *canvas.Canvas) *canvas.Canvas {
	out := pictureOf(src)
	if out == nil || out.Pixels == nil {
		return out
	}
	for y := 0; y < out.Height; y++ {
		row := out.Pixels[y*out.Stride : y*out.Stride+out.Width]
		for x := range row {
			row[x] &= 0xFF000000
		}
	}
	return out
}
