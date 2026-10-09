package css

// Timing is one easing function a transition or animation runs on: the
// linear ramp, one of the named CSS easings, or the four numbers of a
// cubic-bezier(). Ease turns a progress value — 0 at the start, 1 at the
// end — into the value to draw, so a transition that reads room to grow can
// accelerate or overshoot on the way.
type Timing struct {
	x1, y1, x2, y2 float64
}

// Linear is the default easing: no curve, the value equals the progress. It
// is also the reading of a Timing that was never written.
var Linear = Timing{x2: 1, y2: 1}

// Ease applies the easing to a progress in [0,1]. Cubic-bezier easing solves
// the curve's horizontal position for its parameter, then returns the height
// there; a valid cubic-bezier() keeps its x monotonic, so bisection settles.
func (t Timing) Ease(p float64) float64 {
	if p <= 0 {
		return 0
	}
	if p >= 1 {
		return 1
	}
	if t.x1 == 0 && t.y1 == 0 && t.x2 == 1 && t.y2 == 1 {
		return p
	}
	lo, hi := 0.0, 1.0
	for i := 0; i < 24; i++ {
		mid := (lo + hi) / 2
		if bezierX(t, mid) < p {
			lo = mid
		} else {
			hi = mid
		}
	}
	u := (lo + hi) / 2
	return bezierY(t, u)
}
