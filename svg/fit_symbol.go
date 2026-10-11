package svg

// fitSymbol is the size a referenced `<symbol>` or `<svg>` is drawn at: the
// viewBox it was drawn with, mapped onto the width and height the `<use>` asked
// for. A symbol is a template for a picture of its own, and a nested drawing is
// a viewport of its own, and this is what puts that picture where it was asked
// for and at the size it was asked for.
//
// A `<use>` that asks for no size gets the whole drawing, which is what the
// 100% an absent width and height mean — unless the target is a `<svg>`, which
// has a size of its own to fall back on. A target with no viewBox has no size of
// its own to map and is drawn where it stands; and a size of nothing holds
// nothing, so a width or height of zero draws nothing at all, as SVG has it.
func (img *Image) fitSymbol(e, target *element, st *Style, warn func(string, ...any)) {
	if target.Name != "symbol" && target.Name != "svg" {
		return
	}
	vb, ok := parseViewBox(target.attr("viewBox"))
	if !ok {
		return
	}
	w, h := img.ViewBox[2], img.ViewBox[3]
	// A `<svg>` has a viewport of its own, so the size it was written with is
	// the size it keeps when the `<use>` does not ask for one; a `<symbol>` has
	// no size of its own — everything it is drawn at comes from the use.
	isSVG := target.Name == "svg"
	if e.hasAttr("width") {
		w = readLength(warn, e, "width", "width")
	} else if isSVG && target.hasAttr("width") {
		w = readLength(warn, target, "width", "width")
	}
	if e.hasAttr("height") {
		h = readLength(warn, e, "height", "height")
	} else if isSVG && target.hasAttr("height") {
		h = readLength(warn, target, "height", "height")
	}
	if w <= 0 || h <= 0 {
		st.Hidden = true
		return
	}
	st.Transform = st.Transform.Mul(fitTransform(vb, w, h))
}
