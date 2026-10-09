package css

// CopyStyle returns a style that owns its Set map, so a caller can blend or
// freeze values into it without touching the cached cascade's Set. The
// Custom and inherit maps are read-only after the cascade and stay shared.
func CopyStyle(s Style) Style {
	out := make(map[string]bool, len(s.Set))
	for k := range s.Set {
		out[k] = true
	}
	s.Set = out
	return s
}
