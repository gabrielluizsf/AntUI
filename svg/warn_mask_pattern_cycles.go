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
	gray := map[any]bool{}
	done := map[any]bool{}
	var walk func(n *Node)
	var visitMask func(def *maskDef)
	var visitPattern func(p *pattern)
	walk = func(n *Node) {
		if n == nil {
			return
		}
		if n.mask != nil {
			def := n.mask.def
			switch {
			case gray[def]:
				img.warnings.warn(n.Name,
					"the mask here points back at the mask this picture is inside, which would draw for ever, so it is left off when it is painted")
			case !done[def]:
				visitMask(def)
			}
		}
		// A fill and a stroke may each name a pattern, and both walk the same
		// way: what closes the circle is the pattern, not the side of the
		// shape it was asked for.
		for _, p := range []*pattern{n.Style.fillPattern, n.Style.strokePattern} {
			if p == nil {
				continue
			}
			switch {
			case gray[p]:
				img.warnings.warn(n.Name,
					"the pattern here points back at the pattern this picture is inside, which would draw for ever, so it is left off when it is painted")
			case !done[p]:
				visitPattern(p)
			}
		}
		for _, k := range n.Kids {
			walk(k)
		}
	}
	visitMask = func(def *maskDef) {
		if done[def] {
			return
		}
		gray[def] = true
		walk(def.content)
		gray[def] = false
		done[def] = true
	}
	visitPattern = func(p *pattern) {
		if done[p] {
			return
		}
		gray[p] = true
		walk(p.content)
		gray[p] = false
		done[p] = true
	}
	for _, def := range img.masks {
		visitMask(def)
	}
	for _, p := range img.patterns {
		visitPattern(p)
	}
}
