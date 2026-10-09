package css

func applyLayoutDecl(e *declEnv) []string {
	st, raw, ctx, note := e.st, e.raw, e.ctx, e.note
	switch e.prop {
	case "display":
		if v, ok := parseDisplay(raw); ok {
			st.Display = v
			note()
		}
	case "position":
		if v, ok := parsePosition(raw); ok {
			st.Position = v
			note()
		}
	case "top":
		if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.Top = l
			note()
		}
	case "right":
		if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.Right = l
			note()
		}
	case "bottom":
		if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.Bottom = l
			note()
		}
	case "left":
		if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.Left = l
			note()
		}
	case "z-index":
		if v, ok := parseZIndex(raw); ok {
			st.ZIndex = v
			note()
		}
	}
	return nil
}
