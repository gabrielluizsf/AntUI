package css

func applyOutlineDecl(e *declEnv) []string {
	st, raw, ctx, note := e.st, e.raw, e.ctx, e.note
	switch e.prop {
	case "outline":
		if applyOutline(st, raw, *ctx) {
			note()
		}
	case "outline-width":
		if l, err := parseLengthAt(raw, *ctx); err == nil && !l.IsPct() && !l.Auto() && !l.None() {
			st.OutlineWidth = l.Resolve(*ctx)
			note()
		}
	case "outline-style":
		if s, ok := outlineStyle(raw); ok {
			st.OutlineStyle = s
			note()
		}
	case "outline-color":
		if c, err := ParseColor(raw); err == nil {
			st.OutlineColor = c
			note()
		}
	case "outline-offset":
		if l, err := parseLengthAt(raw, *ctx); err == nil && !l.Auto() && !l.None() {
			st.OutlineOffset = l.Resolve(*ctx)
			note()
		}
	}
	return nil
}
