package css

import (
	"math"
)

// animationParts splits a parsed shorthand into the parallel lists the
// longhands also fill, so the cascade merges them in one finisher. An
// infinite count lands in the list as positive infinity.
func animationParts(as []Animation) (names []string, durs []Time, tims []Timing, dels []Time, iters []float64, dirs, fills []uint8) {
	for _, a := range as {
		names = append(names, a.Name)
		durs = append(durs, a.Duration)
		tims = append(tims, a.Timing)
		dels = append(dels, a.Delay)
		if a.Infinite {
			iters = append(iters, math.Inf(1))
		} else {
			iters = append(iters, a.Iterations)
		}
		dirs = append(dirs, a.Direction)
		fills = append(fills, a.Fill)
	}
	return names, durs, tims, dels, iters, dirs, fills
}
