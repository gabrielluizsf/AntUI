package css

import (
	"strings"
)

// parseDeclaration reads "property: value". A missing ':' is not a
// declaration. A trailing !important is recorded, not parsed as a value.
func parseDeclaration(text string) (Declaration, bool) {
	text = stripComments(text)
	idx := strings.IndexByte(text, ':')
	if idx < 0 {
		return Declaration{}, false
	}
	prop := strings.TrimSpace(text[:idx])
	val := strings.TrimSpace(text[idx+1:])
	if prop == "" {
		return Declaration{}, false
	}
	val, important := splitImportant(val)
	return Declaration{Prop: prop, Raw: val, Important: important}, true
}

// splitImportant peels a trailing "!important" (in any spacing) off a value.
func splitImportant(raw string) (val string, important bool) {
	v := strings.TrimSpace(raw)
	if i := strings.LastIndexByte(v, '!'); i >= 0 {
		rest := strings.TrimSpace(strings.ToLower(v[i+1:]))
		if rest == "important" {
			return strings.TrimSpace(v[:i]), true
		}
	}
	return v, false
}
