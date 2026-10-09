package css

// applyKeyword resolves one of the four cascade keywords on a property.
// initial and revert reset the property to its CSS initial value; inherit (and
// unset on an inherited property) record the intent for StyleUnits to fold in
// the body's computed value after the cascade. note() keeps the cascade
// honest: a real value that wins later clears the recorded intent, and so does
// a later keyword.
func applyKeyword(st *Style, note func(), kw, prop string) []string {
	switch kw {
	case "initial", "revert":
		applyInitial(st, prop)
		note()
		return nil
	case "inherit":
		note()
		st.inherit[prop] = true
		return nil
	case "unset":
		if inheritedProps[prop] {
			note()
			st.inherit[prop] = true
			return nil
		}
		applyInitial(st, prop)
		note()
		return nil
	}
	return nil
}
