package template

import (
	"math"

	"github.com/gabrielluizsf/antui/template/css"
)

// animationFrame turns the clock into the 0..1 progress read of one
// animation: the delay is skipped first, then the iteration and its
// direction pick the frame. on reports whether the animation is playing or
// pinned there by a fill-mode; pinned animations still return frames.
func animationFrame(a css.Animation, now float64) (float64, bool) {
	dur := a.Duration.Sec()
	if dur <= 0 {
		return 0, false
	}
	el := now - a.Delay.Sec()
	if el < 0 {
		if a.Fill == css.FillBackwards || a.Fill == css.FillBoth {
			return 0, true
		}
		return 0, false
	}
	if !a.Infinite && a.Iterations >= 0 && el >= dur*a.Iterations {
		if a.Fill == css.FillForwards || a.Fill == css.FillBoth {
			return 1, true
		}
		return 0, false
	}
	cycle := int(el / dur)
	local := math.Mod(el, dur) / dur
	switch a.Direction {
	case css.AnimReverse:
		return 1 - local, true
	case css.AnimAlternate:
		if cycle%2 == 1 {
			return 1 - local, true
		}
		return local, true
	case css.AnimAlternateReverse:
		if cycle%2 == 1 {
			return local, true
		}
		return 1 - local, true
	}
	return local, true
}
