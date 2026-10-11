package canvas

import "math"

// appendDashes cuts one run of points into the stretches of it that are on,
// beginning at whatever point of the pattern the offset lands on. The walk goes
// by distance along the run rather than by its points, because a dash falls
// where the distances add up, and the points of a curve are not evenly spaced.
func (p *Path) appendDashes(sub []Point, pattern []float64, offset, period float64) {
	phase, gone := dashPhase(pattern, offset, period)
	rest := pattern[phase] - gone
	from, at, gone, ok := dashStart(sub, gone)
	if !ok {
		return
	}
	d := &dasher{p: p, pattern: pattern, phase: phase, rest: rest, on: phase%2 == 0}
	if d.on {
		p.MoveTo(at.X, at.Y)
		d.open = true
	}
	for i := from + 1; i < len(sub); i++ {
		a, b := sub[i-1], sub[i]
		seg := math.Hypot(b.X-a.X, b.Y-a.Y)
		// The segment the offset landed in is cut short by however far into it
		// the walk had already come.
		if i == from+1 {
			a = at
			seg -= gone
		}
		d.cut(a, b, seg)
	}
}
