package svg

// place puts where an element says its writing starts onto the pieces of that
// writing. `x` and `y` may name a place for every letter rather than one place
// for the whole of it, and both are read here: the list runs through the
// letters of the pieces in the order they were written, taking one value for
// each letter even where the value does not land — a piece a nearer element
// has already placed keeps that place, but its letters still count against the
// list, so the value after them goes to the letter after them. `dx` and `dy`
// are a move of the pen rather than a place for it, and are never one per
// letter.
func place(runs []TextRun, e *element, warn func(string, ...any)) {
	placeAt(runs, e, "x", warn)
	placeAt(runs, e, "y", warn)
	if v, ok := coord(e, "dx", warn); ok {
		runs[0].DX += v
	}
	if v, ok := coord(e, "dy", warn); ok {
		runs[0].DY += v
	}
}
