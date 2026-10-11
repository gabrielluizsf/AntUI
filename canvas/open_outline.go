package canvas

// openOutline is the band along an open run: the same rectangles and joins,
// with a cap at either end so the line stops where it was told to. It is a list
// of pieces rather than one walk, because the fill is a union and a piece does
// not have to be connected to the next one to be filled.
func openOutline(sub []Point, half float64, style StrokeStyle, limit float64) [][]Point {
	var out [][]Point
	out = append(out, segmentOutline(sub[0], sub[1], half))
	for i := 1; i+1 < len(sub); i++ {
		out = append(out, segmentOutline(sub[i], sub[i+1], half))
		out = append(out, joinOutline(sub[i-1], sub[i], sub[i+1], half, style, limit))
	}
	out = append(out, segmentOutline(sub[len(sub)-2], sub[len(sub)-1], half))
	// Each cap reaches away from the line, so the first one looks back the way
	// the line came and the last one carries on the way it was going.
	out = append(out, capOutline(sub[0], unit(sub[0].X-sub[1].X, sub[0].Y-sub[1].Y), half, style.Cap))
	out = append(out, capOutline(sub[len(sub)-1], unit(sub[len(sub)-1].X-sub[len(sub)-2].X, sub[len(sub)-1].Y-sub[len(sub)-2].Y), half, style.Cap))
	return out
}
