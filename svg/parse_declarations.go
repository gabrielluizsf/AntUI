package svg

import "strings"

// parseDeclarations reads the inside of a `style="…"` attribute, which is a list
// of `name: value` pairs separated by semicolons. The last one of a name is the
// one that counts, as in CSS.
func parseDeclarations(s string) map[string]string {
	out := map[string]string{}
	for _, decl := range strings.Split(s, ";") {
		name, value, ok := strings.Cut(decl, ":")
		if !ok {
			continue
		}
		name = strings.ToLower(strings.TrimSpace(name))
		if name == "" {
			continue
		}
		out[name] = strings.TrimSpace(value)
	}
	return out
}
