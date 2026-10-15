package svg

// imageNode is the `<image>` read whole: where the picture comes from, the box
// it goes into, how it is fitted there, and the opacity it is drawn at. It
// answers the node to paint, with nothing on it where there was nothing to
// paint — a drawing that asked for a picture this cannot fetch or read leaves
// the hole and says which it was, the same way a `url(#id)` that is not there
// leaves a shape unfilled.
//
// The picture is fetched while the drawing is read rather than when it is
// painted, because reading is where every other warning this package gives
// comes from and where a drawing is shared: what went wrong is said once, to
// whoever called Parse, instead of once a frame for as long as the drawing is
// on screen.
func (img *Image) imageNode(e *element, st Style, warn func(string, ...any)) *Node {
	n := &Node{Name: e.Name, Style: st}
	href := e.attr("href")
	if href == "" {
		// The older spelling of the same attribute, which is the one a drawing
		// written for SVG 1.1 carries — the same fallback a `<use>` takes.
		href = e.attr("xlink:href")
	}
	if href == "" {
		warn("it has no href, so there is no picture for it to draw")
		return n
	}
	box := img.pictureBox(e, warn)
	if box[2] <= 0 || box[3] <= 0 {
		// A box of no area is where the picture would go and there is no there,
		// the same as a rectangle with no width: nothing is drawn, and this
		// says so rather than the drawing quietly keeping a picture nobody can
		// see. Fetching first would read a file only to throw it away.
		warn("it has no width or height to draw the picture into, so the picture is left out")
		return n
	}
	fit, align := readPreserveAspect(e.attr("preserveAspectRatio"), warn)
	if st.Opacity <= 0 {
		// Opacity of nothing: the picture would come out transparent however
		// good it is, so it is not worth fetching. Unlike display:none this is
		// not about hiding the element from the drawing — it is already not
		// going to be seen, so there is nothing to say about it either.
		return n
	}
	pixels := img.pictureFor(st.base, href, warn)
	if pixels == nil {
		return n
	}
	if st.Opacity < 1 {
		// The opacity is put into pixels of their own rather than onto the
		// ones in the cache, which every other element pointing at the same
		// address shares: what is in the cache is the picture as it was
		// written, and one element drawing it faint does not fade the rest.
		pixels = fadePicture(pixels, st.Opacity)
	}
	n.Pic = &Picture{Pixels: pixels, Box: box, Fit: fit, Align: align}
	return n
}
