package css

// parseMedia reads an @media prelude into a Media. The full grammar is
// walked — not/only, and/or, comma-separated query lists — while the
// features the viewport answers constrain the rule (its two edges, the
// window's shape, the display's density and the system's color scheme);
// everything else warns and is treated as satisfied, so a canvas never drops
// a rule it merely cannot measure.
func parseMedia(prelude string) (Media, []string) {
	var m Media
	var warns []string
	for _, branch := range splitTopLevel(prelude) {
		qs, w := parseMediaQuery(branch)
		warns = append(warns, w...)
		m.Queries = append(m.Queries, qs...)
	}
	return m, warns
}
