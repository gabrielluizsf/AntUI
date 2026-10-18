package svg

import "strings"

// withVars settles the custom properties the element declares over the ones it
// inherited, then resolves every var() an attribute wrote out of the cascade
// those make: this is the computed value CSS speaks of, and what the element is
// read with. A value whose var() nothing answers is left out, as CSS drops a
// declaration that is invalid at computed-value time, so the element keeps
// whatever it inherited instead of being painted with a name it never resolved.
func (s *Style) withVars(attrs map[string]string) map[string]string {
	if declared := customProps(attrs); len(declared) > 0 {
		s.Custom = mergedCustom(s.Custom, declared)
	}
	out := make(map[string]string, len(attrs))
	for name, raw := range attrs {
		if strings.HasPrefix(name, "--") {
			continue
		}
		if v := resolveVars(raw, s.Custom); !strings.Contains(v, "var(") {
			out[name] = v
		}
	}
	return out
}
