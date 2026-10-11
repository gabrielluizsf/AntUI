package canvas

// ringOutline is the band around a closed loop: a rectangle along every segment
// and the wedge the two rectangles meeting at a corner leave between them. The
// pieces are kept apart rather than stitched into one loop, because a stroke is
// filled as a union: walking from one rectangle to the next would draw seams
// across the band, and those seams cancel each other out under the nonzero rule.
func ringOutline(sub []Point, half float64, style StrokeStyle, limit float64) [][]Point {
	var out [][]Point
	for i := range sub {
		a, b := sub[i], sub[(i+1)%len(sub)]
		if a == b {
			continue
		}
		out = append(out, segmentOutline(a, b, half))
		// The corner is at the far end of this segment, which is where the
		// next one starts, so the join needs the two points around it.
		if wedge := joinOutline(a, b, sub[(i+2)%len(sub)], half, style, limit); len(wedge) > 0 {
			out = append(out, wedge)
		}
	}
	return out
}
