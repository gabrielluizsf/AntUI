package canvas

import (
	"cmp"
	"math"
	"slices"
)

// How many lines a pixel row is cut into before its coverage is worked out, and
// with them the finest step along an edge a fill can describe. Sixteen lines
// put the error of a shallow edge at a fraction of a pixel, which is as fine as
// a screen can show; the coverage of a pixel is not spread out, it is gone, so
// going finer buys nothing and costs work in fours.
const fillSubLines = 16

// FillPath fills the inside of a path with a colour, painted through the
// canvas's clip and blend the way every other draw call is. The rule decides
// which insides count: nonzero for shapes that wind around a point, where two
// nested rings drawn the same way are a ring and a hole, and evenodd for shapes
// that must leave a hole wherever they cross themselves.
func (cv *Canvas) FillPath(path *Path, c Color, rule FillRule) {
	if path == nil || c.A() == 0 {
		return
	}
	minX, minY, maxX, maxY, ok := path.Bounds()
	if !ok {
		return
	}
	// The band is the path's box clipped to the clip region: a path reaching out
	// of the canvas is still filled where the canvas is, and one lying wholly
	// outside it is not filled at all.
	x0 := clampInt(int(math.Floor(minX)), cv.Clip.X, cv.Clip.X+cv.Clip.Width)
	y0 := clampInt(int(math.Floor(minY)), cv.Clip.Y, cv.Clip.Y+cv.Clip.Height)
	x1 := clampInt(int(math.Ceil(maxX)), cv.Clip.X, cv.Clip.X+cv.Clip.Width)
	y1 := clampInt(int(math.Ceil(maxY)), cv.Clip.Y, cv.Clip.Y+cv.Clip.Height)
	if x0 >= x1 || y0 >= y1 {
		return
	}
	for _, s := range fillCoverage(path, x0, y0, x1, y1, rule) {
		cv.fillCoverageSpan(s, c)
	}
}

// span is a run of pixels in one row written at one coverage, which is all the
// fill has to say about a pixel: how much of it the shape covers.
type span struct {
	x0, x1 int
	y      int
	a      float64
}

// fillCoverageSpan paints one run. A run at the colour's own alpha goes through
// the fastest write the canvas has, a partial one scales the alpha and goes a
// pixel at a time, because a blend is a blend of one pixel and not of a run.
func (cv *Canvas) fillCoverageSpan(s span, c Color) {
	a := uint8(float64(c.A())*s.a + 0.5)
	if a == 0 {
		return
	}
	switch {
	case a != c.A():
		sc := RGBA(c.R(), c.G(), c.B(), a)
		for x := s.x0; x < s.x1; x++ {
			cv.Pixel(x, s.y, sc)
		}
	case cv.Pixels != nil:
		cv.blendRun(s.y, s.x0, s.x1, c)
	default:
		cv.narrowRun(s.x0, s.x1, s.y, c)
		cv.markRun(s.y, s.x0, s.x1)
	}
}

// fillCoverage walks the path a row at a time. Each row is cut into lines, each
// line is asked which stretches of it fall inside the shape, and the exact
// length of every stretch is added into the row of pixels above it. Sixteen
// lines down and an exact width across is what makes an edge smooth in both
// directions rather than stair-stepped in one and smeared in the other.
func fillCoverage(path *Path, x0, y0, x1, y1 int, rule FillRule) []span {
	edges := pathEdges(path, y0, y1)
	if len(edges) == 0 {
		return nil
	}
	weight := 1.0 / fillSubLines
	row := make([]float64, x1-x0)
	var (
		spans  []span
		cross  []crossing
		inside []crossing
	)
	for y := y0; y < y1; y++ {
		clear(row)
		for line := range fillSubLines {
			y := float64(y) + (float64(line)+0.5)*weight
			cross, inside = insideRuns(edges, y, rule, cross[:0], inside[:0])
			for i := 0; i+1 < len(inside); i += 2 {
				addCoverage(row, x0, inside[i].x, inside[i+1].x)
			}
		}
		spans = appendRow(spans, row, x0, y, weight)
	}
	return spans
}

// crossing is where a horizontal line meets an edge, and which way the edge was
// heading when it was met. Sorted by x, every second pair of crossings is a
// stretch of line that falls inside the shape.
type crossing struct {
	x   float64
	dir int
}

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

// addCoverage adds the length of a stretch of line that falls inside, spread
// over the pixels it falls on. The stretch is clipped to the band, and a pixel
// it starts or stops inside is worth only the part of it the line really
// covers, which is the whole of the smoothing along a near-vertical edge.
func addCoverage(row []float64, x0 int, from, to float64) {
	from, to = max(from, float64(x0)), min(to, float64(x0+len(row)))
	if to <= from {
		return
	}
	first := max(int(math.Floor(from)), x0)
	last := min(int(math.Ceil(to)), x0+len(row))
	for x := first; x < last; x++ {
		row[x-x0] += min(to, float64(x+1)) - max(from, float64(x))
	}
}

// appendRow turns the coverage of one row into runs of neighbouring pixels under
// the same coverage, leaving out the pixels the shape did not reach at all.
// Coverages are compared as the alpha byte they turn into, so two pixels count
// as equal exactly when they would be painted the same.
func appendRow(spans []span, row []float64, x0, y int, weight float64) []span {
	start, level := -1, 0
	for i, v := range row {
		a := int(v*weight*255 + 0.5)
		if a == level {
			continue
		}
		if level > 0 {
			spans = append(spans, span{x0 + start, x0 + i, y, float64(level) / 255})
		}
		start, level = i, a
	}
	if level > 0 {
		spans = append(spans, span{x0 + start, x0 + len(row), y, float64(level) / 255})
	}
	return spans
}

// edge is one straight piece of a path, kept only where it can cross a line
// inside the band being filled.
type edge struct {
	y0, y1 float64
	dxdy   float64
	dir    int
	xAtY0  float64
}

// pathEdges turns every subpath into the edges a line can cross: each is a run
// of points joined back up to its first when closed, and left out when it is a
// single point, which has no area to cross.
func pathEdges(path *Path, y0, y1 int) []edge {
	pts, closed := path.Points()
	edges := make([]edge, 0, len(pts))
	for i, sub := range pts {
		if len(sub) < 2 {
			continue
		}
		if closed[i] {
			sub = slices.Concat(sub, sub[:1])
		}
		for j := 1; j < len(sub); j++ {
			if e, ok := makeEdge(sub[j-1], sub[j], y0, y1); ok {
				edges = append(edges, e)
			}
		}
	}
	return edges
}

// makeEdge is a segment ready to be sampled: clipped to the rows the band is
// being filled over, its slope inverted into "how far right for one down", and
// the way it was heading recorded, which is what tells the nonzero rule which
// side is inside. A level edge is left out — a line lies along it for a whole
// row, and that is not a crossing but a tie.
//
// An edge lying wholly past one side of the band in x is still kept: a stretch
// of line it bounds may reach into the band from outside it, and without its
// crossing the pairs come out lopsided — the stretch that reaches in never
// gets a start, and the fill drops it rather than painting it, which empties
// the whole shape wherever a shape reaches past its band with a sloped edge.
// The crossing costs nothing even where the stretch it pairs up lies wholly
// outside: the coverage of a stretch is clipped to the band where it is added,
// so a pair that never comes near the band paints no pixel at all.
func makeEdge(a, b Point, y0, y1 int) (edge, bool) {
	if a.Y == b.Y {
		return edge{}, false
	}
	dir := 1
	if b.Y < a.Y {
		a, b, dir = b, a, -1
	}
	top, bottom := max(a.Y, float64(y0)), min(b.Y, float64(y1))
	if bottom <= top {
		return edge{}, false
	}
	e := edge{y0: top, y1: bottom, dxdy: (b.X - a.X) / (b.Y - a.Y), dir: dir}
	e.xAtY0 = a.X + (top-a.Y)*e.dxdy
	return e, true
}
