package svg

// paintedBox is the box a node paints over, in the drawing's own coordinates:
// every shape it draws and every shape inside it draws, each where it was put
// — the shapes carry their transforms already, so this is the box as it lands
// on the page rather than as it was written. A picture has no transform put
// into it — the transform meets it while it is painted — so its box is taken
// through the element's own here for the same result. Writing has no box of
// its own until a font is asked, which is what [walkRuns] and [writingBox] do
// where this walks, and the letters land through the element's transform the
// same way a picture's box does. A node with nothing in it that paints has no
// box, and says so rather than answering an empty one.
func paintedBox(n *Node) ([4]float64, bool) {
	var b boxBounds
	b.walk(n)
	if !b.found {
		return [4]float64{}, false
	}
	return [4]float64{
		b.minMax[0], b.minMax[1],
		b.minMax[2] - b.minMax[0], b.minMax[3] - b.minMax[1],
	}, true
}
