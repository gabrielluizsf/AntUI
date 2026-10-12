package svg

// resolveRefs follows the references that are spent on this element rather than
// inherited, and hands back the ones that live on the node instead: the paints
// and the markers stay on the style, while the clip, the mask and the filter
// list come off it — the children are drawn into the picture this element makes
// and the picture is cut, masked and filtered once, so a child must not each
// find the same reference to follow again. Only the id of a reference is in the
// style; what it turns into is looked up here, once, the same way a fill is.
func (img *Image) resolveRefs(st *Style, warn func(string, ...any)) (clip *clipRegion, def *maskDef, filters []filterOp) {
	// The gradients are read before any node is built, so a `url(#id)` named
	// here can be followed now and the children inherit whatever it turned out
	// to be, the same way they inherit a colour. The marker is the one
	// reference that stays on the style after it has been spent: it is
	// inherited, so what the children draw on their own vertices is this
	// marker, looked up once here rather than once each down there.
	img.resolvePaints(st, warn)
	img.resolveMarkers(st, warn)
	clip = img.resolveClip(st, warn)
	maskRef := st.maskRef
	st.maskRef = ""
	def = img.maskNamed(maskRef, warn)
	filters = st.filters
	st.filters = nil
	filters = img.followFilterRefs(filters, warn)
	return clip, def, filters
}
