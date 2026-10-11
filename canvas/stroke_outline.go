package canvas

import "slices"

// StrokeOutline is the shape a stroke covers, ready to be filled: what a
// shadow, a glow or a second stroke laid over the first one needs. Every piece
// of it is wound the same way round, because the fill is nonzero and two pieces
// that overlapped the wrong way over would cancel each other out and leave a
// hole where the join should be. It is nil when there is nothing to draw.
func StrokeOutline(path *Path, style StrokeStyle) *Path {
	if path == nil || style.Width <= 0 {
		return nil
	}
	limit := style.MiterLimit
	if limit <= 0 {
		limit = DefaultMiterLimit
	}
	pts, closed := path.Points()
	outline := NewPath()
	half := style.Width / 2
	for i, sub := range pts {
		// Compacted into a copy, because the path belongs to the caller and
		// squeezing runs of the same point out of it in place would show.
		sub = slices.Compact(slices.Clone(sub))
		switch {
		case len(sub) < 2:
			// A lone point is a line with no length, and a cap is all that
			// decides what a stroke of one covers: a round or square cap still
			// reaches past the end, so the point is a dot, and a butt cap stops
			// at it and draws nothing.
			if dot := zeroLengthCap(sub, half, style.Cap); len(dot) > 0 {
				outline.AddPolyline(wound(dot), true)
			}
		case closed[i]:
			for _, piece := range ringOutline(sub, half, style, limit) {
				outline.AddPolyline(wound(piece), true)
			}
		default:
			for _, piece := range openOutline(sub, half, style, limit) {
				outline.AddPolyline(wound(piece), true)
			}
		}
	}
	if outline.Empty() {
		return nil
	}
	return outline
}
