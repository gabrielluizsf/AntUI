package css

func applyBoxDecl(e *declEnv) []string {
	st, raw, ctx, note := e.st, e.raw, e.ctx, e.note
	switch e.prop {
	case "width":
		if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.Width = l
			note()
		}
	case "min-width":
		if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.MinWidth = l
			note()
		}
	case "max-width":
		if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.MaxWidth = l
			note()
		}
	case "height":
		if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.Height = l
			note()
		}
	case "min-height":
		if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.MinHeight = l
			note()
		}
	case "max-height":
		if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.MaxHeight = l
			note()
		}
	case "box-sizing":
		switch raw {
		case "border-box":
			st.BoxSizing = true
		case "content-box":
			st.BoxSizing = false
		}
		note()
	}
	return nil
}
