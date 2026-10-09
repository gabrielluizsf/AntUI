package css

// styleUnits resolves one element's style. inherit folds the body's computed
// value into the element for the inherited properties; the body itself is
// computed with inherit off, so the recursion bottoms out in the theme.
func (sh *Sheet) styleUnits(tag string, classes []string, state State, ctx Units, inherit bool) Style {
	got := sh.matchCandidates(tag, classes, state, ctx)
	sortCandidates(got)
	customs := resolveCustoms(got)
	normal, imp := partitionDecls(got)

	var st Style
	st.Custom = make(map[string]string, 8)
	st.inherit = make(map[string]bool, 4)
	for p, v := range customs {
		st.Custom[p] = v
	}
	set := map[string]bool{}
	var warns []string
	for _, d := range normal {
		warns = append(warns, applyDecl(&st, set, st.Custom, customs, d, &ctx)...)
	}
	for _, d := range imp {
		warns = append(warns, applyDecl(&st, set, st.Custom, customs, d, &ctx)...)
	}
	if inherit {
		applyInheritance(&st, set, sh, state, ctx)
	}
	st.Set = set
	resolveCurrentColors(&st)
	finishTransitions(&st)
	finishAnimations(&st)
	sh.Warn = append(sh.Warn, warns...)
	return st
}
