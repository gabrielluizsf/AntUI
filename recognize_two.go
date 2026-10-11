package antui

import (
	"math"
)

// recognizeTwo watches exactly two fingers for a pinch and a twist. Three or
// more is nobody's gesture and is left alone.
func (win *Window) recognizeTwo() {
	g := &win.gesture
	live := win.touches[:0:0] // a view, not a copy of the backing array
	for _, t := range win.touches {
		if !t.Ended {
			live = append(live, t)
		}
	}
	if len(live) != 2 {
		g.twoFingers = false
		return
	}

	a, b := live[0], live[1]
	dx := float64(b.X - a.X)
	dy := float64(b.Y - a.Y)
	dist := math.Hypot(dx, dy)
	angle := math.Atan2(dy, dx)
	g.out.X, g.out.Y = (a.X+b.X)/2, (a.Y+b.Y)/2

	if !g.twoFingers {
		g.twoFingers = true
		g.prevDist, g.prevAngle = dist, angle
		return
	}
	// A pinch is meaningless when the fingers are on top of each other, and
	// dividing by that distance is worse than meaningless.
	if g.prevDist > 1 && dist > 1 && dist != g.prevDist {
		g.out.Pinch = true
		g.out.Scale = dist / g.prevDist
	}
	// Wrap the angle into (-pi, pi] so that crossing the top of the circle
	// is a small turn and not a full one backwards.
	turn := math.Mod(angle-g.prevAngle+3*math.Pi, 2*math.Pi) - math.Pi
	if turn != 0 {
		g.out.Rotate = true
		g.out.Angle = turn
	}
	g.prevDist, g.prevAngle = dist, angle
}
