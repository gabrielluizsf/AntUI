package canvas

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
