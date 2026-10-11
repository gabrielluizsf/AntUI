package svg

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// sliceRegion is the stretch of the picture a sliced box can show: where the
// box falls in the picture's own pixels, which is the only part worth walking.
// The bounds are rounded outwards so that no pixel the box can reach is left
// out — the cut itself is the rectangle [paintPicture] takes to the layer,
// which is the exact one, and this is only about not walking a photograph from
// end to end to show a corner of it. The placement never turns (the turn is on
// the element above it), so the box lands in the picture as a rectangle however
// the drawing has moved.
func sliceRegion(box [4]float64, ox, oy, sx, sy float64, pic *canvas.Canvas) canvas.Area {
	x0 := max(0, int(math.Floor((box[0]-ox)/sx)))
	y0 := max(0, int(math.Floor((box[1]-oy)/sy)))
	x1 := min(pic.Width, int(math.Ceil((box[0]+box[2]-ox)/sx)))
	y1 := min(pic.Height, int(math.Ceil((box[1]+box[3]-oy)/sy)))
	if x1 <= x0 || y1 <= y0 {
		return canvas.Area{}
	}
	return canvas.Area{X: x0, Y: y0, Width: x1 - x0, Height: y1 - y0}
}
