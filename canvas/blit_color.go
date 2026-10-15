package canvas

import (
	"math"

	"github.com/gabrielluizsf/antui/font"
)

// blitColor draws one colour picture into a canvas with the pen at x and the
// baseline at y: scaled from the strike it was cut at to the size this face
// draws at, with its left edge bearingX to the right of the pen and its top
// bearingY above the baseline. The weight and slant of a style lean and
// double the picture the same way they do an outline.
func (f *Face) blitColor(cv *Canvas, cg *font.ColorGlyph, x, y float64, o TextStyle) {
	k := f.colorScale(cg)
	bounds := cg.Image.Bounds()
	srcW, srcH := bounds.Dx(), bounds.Dy()
	dstW, dstH := max(1, int(float64(srcW)*k+0.5)), max(1, int(float64(srcH)*k+0.5))
	left := int(math.Round(x + float64(cg.BearingX)*k))
	top := int(math.Round(y - float64(cg.BearingY)*k))
	step := max(o.Scale, 1)
	for row := range dstH {
		sy := bounds.Min.Y + row*srcH/dstH
		shift := 0
		if o.Italic {
			shift = (dstH - 1 - row) / 4
		}
		for col := range dstW {
			r, g, blue, a := cg.Image.At(bounds.Min.X+col*srcW/dstW, sy).RGBA()
			if a == 0 {
				continue
			}
			// A decoder hands the channels back already multiplied by the
			// alpha; the canvas takes them plain.
			c := RGBA(RGBAParams{R: uint8(r * 255 / a), G: uint8(g * 255 / a), B: uint8(blue * 255 / a), A: uint8(a >> 8)})
			px := left + col + shift
			cv.Pixel(px, top+row, c)
			if o.Bold {
				cv.Pixel(px+step, top+row, c)
			}
		}
	}
}
