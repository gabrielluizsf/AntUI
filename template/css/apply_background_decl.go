package css

func applyBackgroundDecl(e *declEnv) []string {
	st, set, raw, ctx, note := e.st, e.set, e.raw, e.ctx, e.note
	switch e.prop {
	case "background":
		if applyBackground(st, raw, *ctx) {
			note()
			for _, kp := range backgroundKeys {
				set[kp] = true
			}
		}
	case "background-color":
		if c, err := ParseColor(raw); err == nil {
			st.Background = c
			note()
		}
	case "background-image":
		if v, ok := parseBackgroundImage(raw, *ctx); ok {
			st.BackgroundImages = v
			note()
		}
	case "background-position":
		if v, ok := parseBackgroundPosition(raw, *ctx); ok {
			st.BackgroundPos = v
			note()
		}
	case "background-size":
		if v, ok := parseBackgroundSize(raw, *ctx); ok {
			st.BackgroundSize = v
			note()
		}
	case "background-repeat":
		if v, ok := parseBackgroundRepeat(raw); ok {
			st.BackgroundRepeat = v
			note()
		}
	case "background-clip":
		if v, ok := parseBackgroundBox(raw); ok {
			st.BackgroundClip = v
			note()
		}
	case "background-origin":
		if v, ok := parseBackgroundBox(raw); ok {
			st.BackgroundOrigin = v
			note()
		}
	case "background-attachment":
		if v, ok := parseBackgroundAttachment(raw); ok {
			st.BackgroundAttach = v
			note()
		}
	}
	return nil
}
