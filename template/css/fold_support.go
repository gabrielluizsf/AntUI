package css

// foldSupport runs one declaration through the cascade's own fold on a
// scratch style. There is no sheet, no inheritance and a fixed measuring
// context, so the answer depends on the engine and on nothing else — least of
// all on the window the next frame will draw at.
func foldSupport(prop, raw string) ([]string, map[string]bool) {
	var st Style
	st.Custom = make(map[string]string, 4)
	st.inherit = make(map[string]bool, 4)
	set := map[string]bool{}
	ctx := Units{Width: 1000, Height: 1000, Font: DefaultFontSize, Root: DefaultFontSize}
	warns := applyDecl(&st, set, st.Custom, nil, Declaration{Prop: prop, Raw: raw}, &ctx)
	return warns, set
}
