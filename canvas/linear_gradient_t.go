package canvas

import "math"

// linearGradientT is where a point lands along the gradient line: 0 at the
// start of the angle's direction, 1 at the far edge. The line runs through
// the rectangle's centre and its extent is set by projecting the four corners
// onto it, which keeps any angle inside the box and puts the corner
// directions exactly corner to corner.
func linearGradientT(x, y, w, h int, angle float64) float64 {
	dx := math.Sin(angle)
	dy := -math.Cos(angle)
	cw, ch := float64(w)/2, float64(h)/2
	minP, maxP := math.Inf(1), math.Inf(-1)
	for _, c := range [4][2]float64{{0, 0}, {float64(w), 0}, {0, float64(h)}, {float64(w), float64(h)}} {
		p := (c[0]-cw)*dx + (c[1]-ch)*dy
		if p < minP {
			minP = p
		}
		if p > maxP {
			maxP = p
		}
	}
	if span := maxP - minP; span > 0 {
		return ((float64(x)-cw)*dx + (float64(y)-ch)*dy - minP) / span
	}
	return 0
}
