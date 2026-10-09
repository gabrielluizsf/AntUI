package css

import "strings"

// declEnv is one declaration on its way into the style: the target, the
// canonical property and the resolved value, the measurement context, and
// the note closure that records a winning value. The domain handlers read
// it and return the warnings they produced.
type declEnv struct {
	st   *Style
	set  map[string]bool
	ctx  *Units
	prop string
	raw  string
	note func()
}

// applyDecl folds one declaration into the style, expanding shorthands,
// resolving var() references and evaluating length formulas (calc/min/max/
// clamp) under the measurement context. Properties the engine does not know
// are ignored and warned about, exactly as a browser ignores and cools its
// heels on them. The context's Font is updated by a font-size declaration so
// later em lengths see it. applyDecl returns the warnings it produced.
func applyDecl(st *Style, set map[string]bool, customs, cascaded map[string]string, d Declaration, ctx *Units) []string {
	prop := strings.ToLower(strings.TrimSpace(d.Prop))
	raw := strings.TrimSpace(d.Raw)
	if prop == "" || raw == "" {
		return nil
	}
	if set == nil {
		st.Set = make(map[string]bool, 16)
		set = st.Set
	}
	note := func() { delete(st.inherit, prop); set[prop] = true }

	// Custom property detection runs before the vendor-prefix rule, because
	// --foo starts with a dash too. Their values are captured for var() to
	// resolve later.
	if strings.HasPrefix(prop, "--") {
		st.Custom[prop] = raw
		return nil
	}
	if strings.HasPrefix(prop, "-") {
		return []string{fmtErrf("ignoring vendor-prefixed property %q", prop).Error()}
	}
	raw = resolveVars(raw, cascaded)
	raw = strings.TrimSpace(strings.ToLower(raw))
	if raw == "" {
		return nil
	}
	if raw == "initial" || raw == "inherit" || raw == "unset" || raw == "revert" {
		return applyKeyword(st, note, raw, prop)
	}
	return routeDecl(&declEnv{st: st, set: set, ctx: ctx, prop: prop, raw: raw, note: note})
}
