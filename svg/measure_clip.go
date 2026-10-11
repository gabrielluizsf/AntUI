package svg

// measureClip works out the shapes one clipPath cuts with, following the
// `clip-path`s written inside it first, and hands back the clipPath with
// them in it. The gray map is the walk: an id in it is a clipPath whose
// shapes are being worked out right now, so a reference back at one of them
// is a circle the walk would go round for ever without. A clipPath whose
// shapes are already worked out is handed back as it stands, which is what
// stops the walk from reading the same picture twice — and saying everything
// inside it twice with it.
func (img *Image) measureClip(id string, gray map[string]bool) *clipPath {
	cp := img.clips[id]
	if cp == nil || cp.built {
		return cp
	}
	gray[id] = true
	img.resolveClipsIn(cp.content, gray)
	cp.shapes = img.clipShapesOf(cp.content)
	cp.built = true
	delete(gray, id)
	return cp
}
