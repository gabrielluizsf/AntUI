package svg

import "github.com/gabrielluizsf/antui/canvas"

// useNode is the group a `<use>` is replaced by: the style the `<use>` was
// given, moved by the x and y it was given, with the element its href names
// built inside it.
func (img *Image) useNode(e *element, st Style, warn func(string, ...any)) *Node {
	x := readLength(warn, e, "x", "x")
	y := readLength(warn, e, "y", "y")
	if x != 0 || y != 0 {
		st.Transform = st.Transform.Mul(canvas.Translate(x, y))
	}
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
	// second time inside the first. A reference in it is followed here, where
	// the whole drawing is to hand, before what it points at is built.
	filters := st.filters
	st.filters = nil
	filters = img.followFilterRefs(filters, warn)
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
			n.clip = clipUnder(n.clip, st, n, warn)
			n.Kids = append(n.Kids, kid)
			n.mask = maskUnder(def, st, n, warn)
			return n
		}
	}
	n := &Node{Name: "use", Style: st, clip: clip, filters: filters}
	n.mask = maskUnder(def, st, n, warn)
	return n
}
