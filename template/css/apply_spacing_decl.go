package css

func applySpacingDecl(e *declEnv) []string {
	st, set, raw, ctx, note := e.st, e.set, e.raw, e.ctx, e.note
	switch e.prop {
	case "margin":
		if sides, err := parseFourAt(raw, *ctx); err == nil {
			st.Margin = sides
			note()
		}
	case "margin-top":
		setLenSide(st.Margin[:], "margin", set, 0, raw, *ctx)
	case "margin-right":
		setLenSide(st.Margin[:], "margin", set, 1, raw, *ctx)
	case "margin-bottom":
		setLenSide(st.Margin[:], "margin", set, 2, raw, *ctx)
	case "margin-left":
		setLenSide(st.Margin[:], "margin", set, 3, raw, *ctx)
	case "padding":
		if sides, err := parseFourAt(raw, *ctx); err == nil {
			st.Padding = sides
			note()
		}
	case "padding-top":
		setLenSide(st.Padding[:], "padding", set, 0, raw, *ctx)
	case "padding-right":
		setLenSide(st.Padding[:], "padding", set, 1, raw, *ctx)
	case "padding-bottom":
		setLenSide(st.Padding[:], "padding", set, 2, raw, *ctx)
	case "padding-left":
		setLenSide(st.Padding[:], "padding", set, 3, raw, *ctx)
	}
	return nil
}
