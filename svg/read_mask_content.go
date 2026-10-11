package svg

// readMaskContent builds the picture one `<mask>` holds, which is what forceKids
// says: what is inside a definition is not part of the picture where it stands,
// but here it is not being painted either, it is being read as the measure of
// something else.
func (img *Image) readMaskContent(id string, e *element) {
	img.masks[id].content = img.build(e, baseStyle(), true)
}
