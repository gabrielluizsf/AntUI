package css

import (
	"math"
	"strconv"
	"strings"
)

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

// bezierX and bezierY evaluate the cubic applied to a parameter u in [0,1].
func bezierX(t Timing, u float64) float64 {
	return 3*(1-u)*(1-u)*u*t.x1 + 3*(1-u)*u*u*t.x2 + u*u*u
}

func bezierY(t Timing, u float64) float64 {
	return 3*(1-u)*(1-u)*u*t.y1 + 3*(1-u)*u*u*t.y2 + u*u*u
}

// parseTiming reads one easing keyword or a cubic-bezier(x1 y1 x2 y2)
// function. The horizontal controls must stay inside [0,1] for the curve to
// be monotonic, which the parser checks the way a browser does.
func parseTiming(raw string) (Timing, bool) {
	s := strings.TrimSpace(strings.ToLower(raw))
	switch s {
	case "linear":
		return Linear, true
	case "ease":
		return Timing{x1: 0.25, y1: 0.1, x2: 0.25, y2: 1}, true
	case "ease-in":
		return Timing{x1: 0.42, y1: 0, x2: 1, y2: 1}, true
	case "ease-out":
		return Timing{x1: 0, y1: 0, x2: 0.58, y2: 1}, true
	case "ease-in-out":
		return Timing{x1: 0.42, y1: 0, x2: 0.58, y2: 1}, true
	}
	if strings.HasPrefix(s, "cubic-bezier(") && strings.HasSuffix(s, ")") {
		parts := splitFields(s[len("cubic-bezier("):len(s)-1], ',')
		if len(parts) != 4 {
			return Timing{}, false
		}
		vals := [4]float64{}
		for i, p := range parts {
			v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
			if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
				return Timing{}, false
			}
			vals[i] = v
		}
		if vals[0] < 0 || vals[0] > 1 || vals[2] < 0 || vals[2] > 1 {
			return Timing{}, false
		}
		return Timing{x1: vals[0], y1: vals[1], x2: vals[2], y2: vals[3]}, true
	}
	return Timing{}, false
}
