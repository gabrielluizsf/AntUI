package css

// supportsSelector answers selector() with the engine's own selector reader:
// a selector it would keep is one it can evaluate, and one it marks
// unsupported — a combinator, an attribute, an id, a pseudo-element, a
// structural pseudo-class — is one it would drop from the sheet. A selector
// with nothing in it at all selects no widget the engine knows how to name,
// which is how parseSelector reads the empty one.
func supportsSelector(text string) bool {
	sels, _ := parseSelectors(text)
	if len(sels) == 0 {
		return false
	}
	for _, s := range sels {
		if s.Unsupported != "" {
			return false
		}
		if s.Tag == "" && !s.All && s.State == 0 && len(s.Classes) == 0 {
			return false
		}
	}
	return true
}
