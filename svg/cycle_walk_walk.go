package svg

// walk says which mask or pattern every reference under n closes the circle
// with, and descends into the pictures that point at pictures.
func (w *cycleWalk) walk(n *Node) {
	if n == nil {
		return
	}
	if n.mask != nil {
		def := n.mask.def
		switch {
		case w.gray[def]:
			w.img.warnNamed(n.Name,
				"the mask here points back at the mask this picture is inside, which would draw for ever, so it is left off when it is painted")
		case !w.done[def]:
			w.visitMask(def)
		}
	}
	// A fill and a stroke may each name a pattern, and both walk the same way:
	// what closes the circle is the pattern, not the side of the shape it was
	// asked for.
	for _, p := range []*pattern{n.Style.fillPattern, n.Style.strokePattern} {
		if p == nil {
			continue
		}
		switch {
		case w.gray[p]:
			w.img.warnNamed(n.Name,
				"the pattern here points back at the pattern this picture is inside, which would draw for ever, so it is left off when it is painted")
		case !w.done[p]:
			w.visitPattern(p)
		}
	}
	for _, k := range n.Kids {
		w.walk(k)
	}
}
