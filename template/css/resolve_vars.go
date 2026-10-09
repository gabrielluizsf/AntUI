package css

import "strings"

// resolveVars replaces every var(--name[, fallback]) in a raw value with the
// cascaded custom-property value, or the fallback, walking nested definitions
// up to a depth guard so a cycle --a: var(--a) cannot spin forever. A var()
// with no known value and no fallback is left in place, which makes the
// property unrecognised and dropped by applyDecl.
func resolveVars(raw string, customs map[string]string) string {
	return resolveVarsDepth(raw, customs, 0)
}

func resolveVarsDepth(raw string, customs map[string]string, depth int) string {
	if depth > 16 {
		return raw
	}
	out := raw
	for i := 0; i < len(out); i++ {
		if out[i] != 'v' || !strings.HasPrefix(out[i:], "var(") {
			continue
		}
		end := parenEnd(out, i+3)
		if end < 0 {
			return out
		}
		name, fallback, ok := splitVarArg(out[i+4 : end])
		if !ok {
			i = end
			continue
		}
		if v, found := customs[name]; found {
			out = out[:i] + resolveVarsDepth(v, customs, depth+1) + out[end+1:]
			continue
		}
		if fallback != "" {
			out = out[:i] + resolveVarsDepth(fallback, customs, depth+1) + out[end+1:]
			continue
		}
		i = end
	}
	return out
}
