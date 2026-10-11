package canvas

import (
	"cmp"
	"slices"
)

// insideRuns lists the stretches of the line at y that the shape covers. It
// sorts the edges the line meets and walks them while the rule keeps saying
// inside: the nonzero rule adds up which way each edge went and calls a point
// inside when the total is not zero, and the evenodd rule only counts how many
// edges were crossed. An odd count means the shape was not closed where the
// line ran out, and that last stretch is dropped rather than painted out to the
// edge of the band.
func insideRuns(edges []edge, y float64, rule FillRule, cross, inside []crossing) ([]crossing, []crossing) {
	for _, e := range edges {
		if y >= e.y0 && y < e.y1 {
			cross = append(cross, crossing{e.xAtY0 + (y-e.y0)*e.dxdy, e.dir})
		}
	}
	slices.SortFunc(cross, func(a, b crossing) int { return cmp.Compare(a.x, b.x) })

	winding, start := 0, -1
	for i, c := range cross {
		if rule == FillNonZero {
			winding += c.dir
		} else {
			winding ^= 1
		}
		switch now := winding != 0; {
		case now && start < 0:
			start = i
		case !now && start >= 0:
			inside = append(inside, cross[start], cross[i])
			start = -1
		}
	}
	return cross, inside
}
