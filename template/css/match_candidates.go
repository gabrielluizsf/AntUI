package css

// matchCandidates collects every rule whose media and selector match the
// element, together with the specificity and order the cascade sorts by.
func (sh *Sheet) matchCandidates(tag string, classes []string, state State, ctx Units) []candidate {
	var got []candidate
	for _, r := range sh.rules {
		if !r.Media.matches(Viewport{
			Width: ctx.Width, Height: ctx.Height, Scale: ctx.Scale, Scheme: ctx.Scheme,
		}) {
			continue
		}
		for _, sel := range r.Selectors {
			if !sel.match(tag, classes, state) {
				continue
			}
			_, b, c := sel.specificity()
			got = append(got, candidate{[3]int{0, b, c}, r.Order, r.Decls})
			break
		}
	}
	return got
}
