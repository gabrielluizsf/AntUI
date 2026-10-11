package svg

// readSize takes the size of the drawing from its width and height, falling
// back to the viewBox and then to a hundred square, and makes sure the viewBox
// is one that describes the drawing: an absent or unreadable one is taken from
// the size. There is always one of the three, so this never fails — a drawing
// that says nothing about its size is still a drawing, and it is painted at
// the size a caller asks for either way.
func (img *Image) readSize(e *element) {
	img.Width, img.Height = 100, 100
	hasW, hasH := e.hasAttr("width"), e.hasAttr("height")
	if w, ok := parseLength(e.attr("width")); ok {
		img.Width = w
	} else {
		hasW = false
	}
	if h, ok := parseLength(e.attr("height")); ok {
		img.Height = h
	} else {
		hasH = false
	}
	if v, ok := parseViewBox(e.attr("viewBox")); ok {
		img.ViewBox = v
		if !hasW && !hasH {
			// A drawing with only a viewBox is sized by it, which is how a file
			// meant to scale says how big it is without saying.
			img.Width, img.Height = v[2], v[3]
		}
		return
	}
	// No viewBox: the drawing's coordinates are the size, starting at the
	// origin, so that is what it is.
	img.ViewBox = [4]float64{0, 0, img.Width, img.Height}
}
