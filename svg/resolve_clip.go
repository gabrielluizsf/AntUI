package svg

// resolveClip follows the `clip-path` an element named, which only the whole
// drawing can say anything about, and answers the clip that element is cut by:
// the shapes of the clipPath it names, moved into the element's own place
// before anyone paints, since that is where the element is going. The shapes
// in the drawing are shared by everything that points at them, so each takes a
// copy rather than moving the ones everyone else is holding.
//
// A reference the drawing cannot follow keeps nothing and says so, which is
// what the rest of this package does with a reference to nowhere: a `<use>`
// that names nothing draws nothing, a fill that names a gradient that is not
// there is not filled, and a clip that names a clipPath that is not there cuts
// everything away.
//
// A reference from inside a `<clipPath>` cannot be followed here at all: every
// clipPath is being built at this point and none of them has its shapes yet —
// one of them may be named further down the file than the place it is pointed
// at. It waits as the name alone, on the node it was written on, until the
// picture it is written in is walked for its shapes and the name can be
// followed for real — see [Image.resolveClipsIn].
//
// The reference is cleared as it is spent, the same way a fill's is, so that a
// `<g clip-path>` says what it has to say once and not once for every shape
// inside it, and so that the children never see a clip of their own to
// inherit — the picture they are drawn into is what gets cut.
func (img *Image) resolveClip(st *Style, warn func(string, ...any)) *clipRegion {
	if st.clipRef == "" {
		return nil
	}
	id := st.clipRef
	st.clipRef = ""
	cp := img.clips[id]
	switch {
	case cp == nil:
		warn("the clip-path names #%s, which the drawing does not have, so nothing of what it cuts is drawn", id)
		return &clipRegion{}
	case !cp.built:
		return &clipRegion{pending: id}
	}
	return img.clipRegionOf(cp, *st)
}
