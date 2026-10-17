package svg

// parseDeclarations reads the inside of a `style="…"` attribute into the value
// each name was last given, the last one of a name counting as in CSS. It is
// the map view of [parseDecls], for the callers that ask about one name.
func parseDeclarations(s string) map[string]string {
	out := map[string]string{}
	for _, d := range parseDecls(s) {
		out[d.name] = d.value
	}
	return out
}
