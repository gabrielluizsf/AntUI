package canvas

import "math"

// radialGradientT scales the distance from the centre up to the ending shape
// into a point along the gradient line. The ending shape's size keywords are
// resolved against the rectangle, and a circle uses the largest circle that
// the keywords name while an ellipse keeps its two radii live.
func radialGradientT(x, y float64, w, h int, g Gradient) float64 {
	cx := g.CenterX * float64(w)
	cy := g.CenterY * float64(h)
	if g.Shape == GradCircle {
		r := circleRadius(cx, cy, w, h, g.Size)
		if r <= 0 {
			return 0
		}
		px, py := x-cx, y-cy
		return math.Sqrt(px*px+py*py) / r
	}
	var rx, ry float64
	switch g.Size {
	case GSClosestSide:
		rx = min(cx, float64(w)-cx)
		ry = min(cy, float64(h)-cy)
	case GSFarthestSide:
		rx = max(cx, float64(w)-cx)
		ry = max(cy, float64(h)-cy)
	case GSClosestCorner:
		ex, ey := nearestCorner(cx, cy, w, h)
		rx, ry = math.Abs(ex-cx), math.Abs(ey-cy)
	default:
		ex, ey := farthestCorner(cx, cy, w, h)
		rx, ry = math.Abs(ex-cx), math.Abs(ey-cy)
	}
	if rx <= 0 || ry <= 0 {
		return 0
	}
	px, py := x-cx, y-cy
	return math.Sqrt(px*px/(rx*rx) + py*py/(ry*ry))
}
