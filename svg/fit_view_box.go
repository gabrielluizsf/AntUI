package svg

import "github.com/gabrielluizsf/antui/canvas"

// fitViewBox maps a viewBox onto the box a Picture was given, the way
// preserveAspectRatio asks: the meet, slice or none of the fit and the corner
// or edge the alignment pins the picture to, rather than always a centred
// meeting. A viewBox with no area of its own has no proportions to keep, so it
// is stretched to the box as it is everywhere else it is fitted.
func fitViewBox(view [4]float64, p pictureFit) canvas.Matrix {
	if view[2] <= 0 || view[3] <= 0 {
		return canvas.Scale(p.box[2], p.box[3])
	}
	sx, sy, ox, oy := picturePlacement(p, view[2], view[3])
	return canvas.Translate(ox-view[0]*sx, oy-view[1]*sy).Mul(canvas.Scale(sx, sy))
}
