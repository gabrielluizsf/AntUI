package canvas

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
