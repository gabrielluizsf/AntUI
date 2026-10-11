package canvas

import "math"

// Length is how far it is along the path: every segment of every subpath, with
// a closed one counted back to the point it started from. Curves were flattened
// into points as they were added, so this measures the polygon the path has
// become — the same walk a stroke takes along it and the same walk a dash is
// cut into. A path with nothing in it, or only points that lead nowhere, is
// nothing long.
func (p *Path) Length() float64 {
	if p == nil {
		return 0
	}
	pts, closed := p.Points()
	total := 0.0
	for i, sub := range pts {
		for j := 1; j < len(sub); j++ {
			total += math.Hypot(sub[j].X-sub[j-1].X, sub[j].Y-sub[j-1].Y)
		}
		if closed[i] && len(sub) > 1 {
			from, to := sub[len(sub)-1], sub[0]
			total += math.Hypot(to.X-from.X, to.Y-from.Y)
		}
	}
	return total
}
