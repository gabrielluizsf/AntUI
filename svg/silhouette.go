package svg

import "github.com/gabrielluizsf/antui/canvas"

// silhouette is a picture of only where the one it is given has something, in
// the colour of a shadow: the alpha is kept and multiplied by the colour's own,
// and the colour's three channels are what the shape is filled with. That is
// what a drop-shadow is cast from, and blurring it is what softens the shadow
// without touching the element over it.
func silhouette(src *canvas.Canvas, c canvas.Color) *canvas.Canvas {
	out := pictureOf(src)
	if out == nil || out.Pixels == nil {
		return out
	}
	rgb := c & 0x00FFFFFF
	fa := float64(c.A()) / 255
	for y := 0; y < out.Height; y++ {
		row := out.Pixels[y*out.Stride : y*out.Stride+out.Width]
		for x, px := range row {
			row[x] = rgb | canvas.Color(uint32(float64(px.A())*fa+0.5))<<24
		}
	}
	return out
}
