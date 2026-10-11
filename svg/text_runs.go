package svg

// textRuns is everything written inside a `<text>` or a `<tspan>`, flattened
// into one list in the order it was written. The writing and the elements
// between it are taken in turn, so a `<tspan>`'s writing goes where it was
// written — after the writing before it, not at the origin.
//
// A child that is not writing at all is left out with a warning of its own: a
// drawing may put a `<rect>` inside a `<text>` for reasons of its own, and
// there is nothing here that knows what to do with it.
func (img *Image) textRuns(e *element, st Style, warn func(string, ...any)) []TextRun {
	runs := make([]TextRun, 0, len(e.Parts))
	for _, p := range e.Parts {
		if p.Elem == nil {
			if p.Text != "" {
				runs = append(runs, TextRun{Text: p.Text, Style: st})
			}
			continue
		}
		if p.Elem.Name == "textPath" {
			// The writing is not hung on a pen at all: it is laid along the
			// outline the href names, with `startOffset` saying how far along
			// it the first letter goes. What it names could not be followed,
			// and then there is nowhere to write, which is worth one sentence
			// rather than the writing turning up in the wrong place.
			along := img.textPathOutline(p.Elem, warn)
			if along == nil {
				continue
			}
			kid := st.with(p.Elem, warn)
			img.resolvePaints(&kid, warn)
			sub := img.textRuns(p.Elem, kid, warn)
			off := startOffset(p.Elem, along, warn)
			for i := range sub {
				sub[i].Path, sub[i].Offset = along, off
			}
			runs = append(runs, sub...)
			continue
		}
		if p.Elem.Name != "tspan" {
			warn("<%s> inside the writing is not writing this package draws, so it is left out", p.Elem.Name)
			continue
		}
		// The `<tspan>` is read the way any other element is read, over the
		// style the writing around it has: a size or a colour on it is the
		// size or colour of what it holds and nothing else.
		kid := st.with(p.Elem, warn)
		img.resolvePaints(&kid, warn)
		runs = append(runs, img.textRuns(p.Elem, kid, warn)...)
	}
	if len(runs) == 0 {
		return nil
	}
	// Where the writing starts is on the element it was written on, and it
	// lands on the pieces of that writing. A piece that a `<tspan>` closer
	// to the writing has already placed keeps the place it was given: the
	// innermost of them is the one that counts, the same way an attribute on
	// a child beats the one written on its parent.
	place(runs, e, warn)
	return runs
}
