package css

// applyVisibilityDecl resets visibility and the two hit-test properties
// pointer-events and cursor: the trio says whether the box shows, answers,
// and what shape the pointer takes over it.
func applyVisibilityDecl(e *declEnv) []string {
	st, raw, note := e.st, e.raw, e.note
	switch e.prop {
	case "visibility":
		if v, ok := parseVisibility(raw); ok {
			st.Visibility = v
			note()
		}
	case "pointer-events":
		if v, ok := parsePointerEvents(raw); ok {
			st.PointerEvents = v
			note()
		}
	case "cursor":
		if v, ok := parseCursor(raw); ok {
			st.Cursor = v
			note()
		}
	}
	return nil
}
