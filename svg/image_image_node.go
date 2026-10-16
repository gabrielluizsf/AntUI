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
// comes from: what went wrong is said once, to whoever called Parse, instead of
// once a frame for as long as the drawing is on screen. An address that is a
// server is the one exception — it is read off the frame, so the drawing is
// built with the hole where the picture will be, the photo arriving in a later
// paint and a failure said from the goroutine that fetched it, which
// [Image.Warnings] reads back after the fact.
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
	box, autoW, autoH := img.pictureBox(e, warn)
	if (!autoW && box[2] <= 0) || (!autoH && box[3] <= 0) {
		// A box of no area is where the picture would go and there is no there,
		// the same as a rectangle with no width: nothing is drawn, and this
		// says so rather than the drawing quietly keeping a picture nobody can
		// see. Fetching first would read a file only to throw it away — except
		// that a side left open is the picture's to say, which is only known
		// once it is read, and that case has an area coming either way.
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
	slot := &Picture{Box: box, Fit: fit, Align: align, declared: box, autoW: autoW, autoH: autoH}
	// The opacity is put into pixels of their own rather than onto the ones in
	// the cache, which every other element pointing at the same address shares:
	// what is in the cache is the picture as it was written, and one element
	// drawing it faint does not fade the rest. The request carries it because a
	// picture still on its way has to be faded when it lands.
	ask := pictureRequest{slot: slot, opacity: st.Opacity, href: href, warn: warn, chain: img.chain}
	if !img.fillPicture(resolveHref(st.base, href), ask) {
		return n
	}
	n.Pic = slot
	return n
}
