package template

// gridAreaBox is the stretch of the tracks an item spans: where it starts on
// screen and how wide it is, gaps included.
func gridAreaBox(sizes, positions []int, start, end, gap int) (int, int) {
	if len(sizes) == 0 {
		return 0, 0
	}
	start = min(max(start, 0), len(sizes)-1)
	end = min(max(end, start+1), len(sizes))
	width := max(end-start-1, 0) * gap
	for track := start; track < end; track++ {
		width += sizes[track]
	}
	return positions[start], width
}
