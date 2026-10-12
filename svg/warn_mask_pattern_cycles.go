package svg

// warnMaskPatternCycles says which reference closes a circle a drawing drew
// between its masks or between its patterns: a mask whose picture masks
// itself — or two that name each other — would draw the same picture while it
// was already being drawn, for ever, and a pattern that fills itself with
// itself is the same round trip. The painting already leaves the reference
// that closes the circle off, so the drawing comes out as a picture rather
// than as a program that does not come back — see [paintMasked] and
// [patternShade]. What was missing was the drawing saying which reference
// closed it, the way [Image.cutMarkerCycles] says for a marker.
//
// It only walks and only says: nothing here is cut, because what the painting
// does when it meets the circle again is what the warning says it does. The
// walk starts at the masks and the patterns themselves — the circles only run
// through pictures that point at pictures, and nothing points back at the
// drawing that names them.
func (img *Image) warnMaskPatternCycles() {
	w := &cycleWalk{img: img, gray: map[any]bool{}, done: map[any]bool{}}
	for _, def := range img.masks {
		w.visitMask(def)
	}
	for _, p := range img.patterns {
		w.visitPattern(p)
	}
}
