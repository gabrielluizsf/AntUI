package svg

// maskNamed is the mask an element named, followed before anything is built —
// where a clip is followed — so that what cannot be followed says so once,
// while the element is being read rather than while it is being painted.
//
// Three things come out of it: nothing at all, where no mask was named or the
// element asks for `mask="none"`; the mask itself; and [maskNowhere], where
// the drawing does not have the id that was named, which keeps nothing.
func (img *Image) maskNamed(ref string, warn func(string, ...any)) *maskDef {
	if ref == "" {
		return nil
	}
	def := img.masks[ref]
	if def == nil {
		warn("the mask names #%s, which the drawing does not have, so nothing of what it masks is drawn", ref)
		return maskNowhere
	}
	return def
}
