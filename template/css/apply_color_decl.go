package css

func applyColorDecl(e *declEnv) []string {
	st, raw, note := e.st, e.raw, e.note
	switch e.prop {
	case "color":
		if c, err := ParseColor(raw); err == nil {
			st.Color = c
			note()
		}
	case "opacity":
		if v, ok := parseOpacity(raw); ok {
			st.Opacity = v
			note()
		}
	}
	return nil
}
