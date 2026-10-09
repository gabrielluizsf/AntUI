package css

// frameStyle computes the style one keyframe paints: its declarations folded
// into the base's custom-property map, so var(--x) reads the element's own
// values. The frame style carries only what its own Set map saw.
func frameStyle(frame Keyframe, base Style, ctx Units) Style {
	var st Style
	st.Custom = base.Custom
	st.inherit = map[string]bool{}
	set := map[string]bool{}
	for _, d := range frame.Decls {
		applyDecl(&st, set, base.Custom, base.Custom, d, &ctx)
	}
	st.Set = set
	resolveCurrentColors(&st)
	return st
}
