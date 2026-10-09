package css

// applyInheritance folds the body's computed style into the element in two
// passes: first the properties whose cascade explicitly chose the inherit
// keyword (or unset on an inherited property), then every inherited property
// the element never mentioned. The body copy is lazy and computed once per
// StyleUnits call with inheritance off, so there is no recursion loop.
func applyInheritance(st *Style, set map[string]bool, sh *Sheet, state State, ctx Units) {
	body := sh.styleUnits("body", nil, state, ctx, false)
	for prop := range st.inherit {
		inheritOne(st, set, body, prop)
	}
	for prop := range inheritedProps {
		if !set[prop] && body.Set[prop] {
			inheritOne(st, set, body, prop)
		}
	}
}
