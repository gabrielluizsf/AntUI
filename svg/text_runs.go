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
		switch {
		case p.Elem == nil:
			if p.Text != "" {
				runs = append(runs, TextRun{Text: p.Text, Style: st})
			}
		case p.Elem.Name == "textPath":
			if sub := img.textPathRuns(p.Elem, st, warn); sub != nil {
				runs = append(runs, sub...)
			}
		case p.Elem.Name == "tspan":
			runs = append(runs, img.childRuns(p.Elem, st, warn)...)
		default:
			warn("<%s> inside the writing is not writing this package draws, so it is left out", p.Elem.Name)
		}
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
