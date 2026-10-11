package svg

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// fitTransform is the transform that maps a viewBox onto a rectangle of the
// given size, keeping the two in the proportion the drawing was made in. A
// drawing whose viewBox is not the same shape as the canvas is fitted inside it
// and centred, the way a picture is, rather than stretched to fill it. The size
// is any two numbers in the coordinates the drawing is being put into: pixels
// when it is painted, the units of the drawing when a `<symbol>` is fitted into
// the box a `<use>` gave it.
func fitTransform(view [4]float64, width, height float64) canvas.Matrix {
	vw, vh := view[2], view[3]
	if vw <= 0 || vh <= 0 {
		return canvas.Scale(width, height)
	}
	scale := math.Min(width/vw, height/vh)
	sx, sy := scale, scale
	// Keep the drawing's own proportions inside the canvas.
	if vw*scale > width {
		sx = width / vw
	}
	if vh*scale > height {
		sy = height / vh
	}
	// The drawing is centred in what is left over, which is nothing at all when
	// the canvas and the viewBox have the same shape, as an icon's usually do.
	ox := (width - vw*sx) / 2
	oy := (height - vh*sy) / 2
	return canvas.Translate(ox-view[0]*sx, oy-view[1]*sy).Mul(canvas.Scale(sx, sy))
}
