package svg

// readDocFit is the drawing's own preserveAspectRatio, kept for the day the
// drawing is read as a picture of its own: a `defer` on the `<image>`
// pointing at it asks to be fitted by exactly this value, which travels with
// the frames. Whether the value also fits the drawing's own viewBox is the
// viewport's business, and not something this reading decides.
func (img *Image) readDocFit(e *element) {
	if !e.hasAttr("preserveAspectRatio") {
		return
	}
	img.docFit, img.docAlign, _ = readPreserveAspect(e.attr("preserveAspectRatio"), discardWarn)
	img.hasDocFit = true
}
