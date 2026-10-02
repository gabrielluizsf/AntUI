package svg

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// A `<use>` is one shape written once and put in the picture as many times as
// it is asked for. The element its href names is built as though it had been
// written where the `<use>` is, with the style the `<use>` and the elements
// above it gave it and never the style of the place it originally stood in —
// which is what makes one definition come out in the colour of every place it
// is used rather than in the colour it was written in.
//
// What it names is usually inside a `<defs>`, which is kept out of the picture
// for exactly this reason, so the same shape is drawn once and not twice.

// useNode is the group a `<use>` is replaced by: the style the `<use>` was
// given, moved by the x and y it was given, with the element its href names
// built inside it.
func (img *Image) useNode(e *element, st Style, warn func(string, ...any)) *Node {
	x := readLength(warn, e, "x", "x")
	y := readLength(warn, e, "y", "y")
	if x != 0 || y != 0 {
		// The translate goes on the right of the transform the `<use>` already
		// carries, which is the inside of it: what is referenced is moved in
		// its own space first and then everything the `<use>` does is applied
		// to the result. SVG spells it out as a translate appended to the end
		// of the transform, and a rectangle's own x is inside the transform on
		// it for the same reason.
		st.Transform = st.Transform.Mul(canvas.Translate(x, y))
	}
	// The clip is followed here rather than by the caller, and for two reasons
	// that both come from what a `<use>` is. It has to be spent before what the
	// use points at is built, or every element inside it would find the same
	// reference to follow and say the same thing about it once each. And it is
	// followed after the x and y have moved the style, so that the clip goes
	// where everything else the `<use>` draws goes — before the viewBox of a
	// `<symbol>` is fitted over it, which moves what is drawn and not the place
	// the `<use>` itself is.
	clip := img.resolveClip(&st, warn)
	// The mask is named and spent here for the same two reasons, and what it
	// turns into is worked out against the box of everything the `<use>` ends
	// up drawing, so it waits until the element it names has been built.
	maskRef := st.maskRef
	st.maskRef = ""
	def := img.maskNamed(maskRef, warn)
	// The filter list is spent here too, on the `<use>` itself rather than on
	// what it points at: the filter is put through the whole picture the
	// `<use>` makes, and what it names would otherwise run the same list a
	// second time inside the first.
	filters := st.filters
	st.filters = nil
	target, id := img.useTarget(e, warn)
	switch {
	case st.Hidden:
		// The `<use>` is not drawn, so nothing behind it is built either.
	case target == nil:
		// useTarget has already said what is wrong with the reference.
	case img.using(id):
		warn("the href names #%s, which is already being drawn through a <use>, so it is left out rather than going round for ever", id)
	default:
		img.fitSymbol(e, target, &st, warn)
		if !st.Hidden {
			img.uses = append(img.uses, id)
			kid := img.build(target, st, target.Name == "symbol")
			img.uses = img.uses[:len(img.uses)-1]
			n := &Node{Name: "use", Style: st, clip: clip, filters: filters}
			n.Kids = append(n.Kids, kid)
			n.mask = maskUnder(def, st, n)
			return n
		}
	}
	n := &Node{Name: "use", Style: st, clip: clip, filters: filters}
	n.mask = maskUnder(def, st, n)
	return n
}

// useTarget is the element a `<use>` names and the id it named it by. Only the
// whole drawing can say what an id is, so the reference is followed here
// rather than where it is written. A reference that cannot be followed — none
// at all, one into another drawing, or one to an id this drawing does not have
// — says so and stands for nothing, since a `<use>` with nothing behind it is
// a hole in the picture either way.
func (img *Image) useTarget(e *element, warn func(string, ...any)) (*element, string) {
	href := e.attr("href")
	if href == "" {
		// The older spelling of the same attribute, which is the one a drawing
		// written for SVG 1.1 carries.
		href = e.attr("xlink:href")
	}
	if href == "" {
		warn("it has no href, so there is nothing for it to draw")
		return nil, ""
	}
	id, ok := strings.CutPrefix(href, "#")
	if !ok {
		warn("the href %q names something in another drawing, which is not read", href)
		return nil, ""
	}
	if id == "" {
		warn("the href %q names nothing at all", href)
		return nil, ""
	}
	target := img.ids[id]
	if target == nil {
		warn("the href names #%s, which the drawing does not have", id)
		return nil, ""
	}
	return target, id
}

// using reports whether an id is already being drawn through a `<use>` on the
// way to this one. A reference that names one of them goes round for ever
// rather than drawing a shape, so it is stopped at the first repeat.
func (img *Image) using(id string) bool {
	for _, seen := range img.uses {
		if seen == id {
			return true
		}
	}
	return false
}

// fitSymbol is the size a referenced `<symbol>` is drawn at: the viewBox it was
// drawn with, mapped onto the width and height the `<use>` asked for. A symbol
// is a template for a picture of its own, and this is what puts that picture
// where it was asked for and at the size it was asked for.
//
// A `<use>` that asks for no size gets the whole drawing, which is what the
// 100% an absent width and height mean; a symbol with no viewBox has no size of
// its own to map and is drawn where it stands; and a size of nothing holds
// nothing, so a width or height of zero draws nothing at all, as SVG has it.
func (img *Image) fitSymbol(e, target *element, st *Style, warn func(string, ...any)) {
	if target.Name != "symbol" {
		return
	}
	vb, ok := parseViewBox(target.attr("viewBox"))
	if !ok {
		return
	}
	w, h := img.ViewBox[2], img.ViewBox[3]
	if e.hasAttr("width") {
		w = readLength(warn, e, "width", "width")
	}
	if e.hasAttr("height") {
		h = readLength(warn, e, "height", "height")
	}
	if w <= 0 || h <= 0 {
		st.Hidden = true
		return
	}
	st.Transform = st.Transform.Mul(fitTransform(vb, w, h))
}
