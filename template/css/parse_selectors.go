package css

// parseSelectors splits a comma-separated selector list into one Selector per
// part. Every part is read tolerantly — an unsupported part marks its Selector
// as never-matching and warns instead of failing the sheet.
func parseSelectors(text string) ([]Selector, []string) {
	var out []Selector
	var warns []string
	for _, part := range splitTopLevel(text) {
		s, w := parseSelector(part)
		warns = append(warns, w...)
		out = append(out, s)
	}
	return out, warns
}
