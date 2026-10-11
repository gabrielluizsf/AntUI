package canvas

import "math"

// arcSweepCoverage is how much of the angular edge a point at dx, dy is on,
// from 0 outside the sweep to 1 well inside it, feathered over the last feather
// radians at either end. A full turn has no edge and comes back 1.
func arcSweepCoverage(dx, dy int, start, span, dir, feather float64, full bool) float64 {
	if full {
		return 1
	}
	delta := math.Mod(dir*(math.Atan2(float64(dy), float64(dx))-start), 2*math.Pi)
	if delta < 0 {
		delta += 2 * math.Pi
	}
	if delta > span {
		return 0
	}
	acov := 1.0
	if delta < feather {
		acov = delta / feather
	}
	if span-delta < feather {
		acov = min(acov, (span-delta)/feather)
	}
	return max(acov, 0)
}
