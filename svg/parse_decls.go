package svg

import "strings"

// parseDecls reads the inside of a rule or of a `style="…"` attribute: a list
// of `name: value` pairs separated by semicolons, with any `!important` a value
// was written with kept. A name is read the way CSS reads one, without case,
// except a custom property, whose name is its own.
func parseDecls(s string) []cssDecl {
	var out []cssDecl
	for _, part := range splitTopSemicolon(s) {
		name, value, ok := splitDecl(part)
		if !ok {
			continue
		}
		value, important := cutImportant(value)
		out = append(out, cssDecl{name: name, value: value, important: important})
	}
	return out
}

// declName is the name a declaration is read under.
func declName(name string) string {
	if strings.HasPrefix(name, "--") {
		return name
	}
	return keyword(name)
}
