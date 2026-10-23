package css

import "strings"

// resolveVars replaces every var(--name[, fallback]) in a raw value with the
// cascaded custom-property value, or the fallback, walking nested definitions
// up to a depth guard so a cycle --a: var(--a) cannot spin forever. A var()
// with no known value and no fallback is left in place, which makes the
// property unrecognised and dropped by applyDecl.
//
// A value that would resolve to more than [maxVarBytes] is handed back as it
// was written, which drops it the same way: ten properties each naming the
// one before ten times is a few hundred bytes that read into gigabytes, and
// the CSS Variables spec has such a value be invalid rather than worked out.
func resolveVars(raw string, customs map[string]string) string {
	budget := maxVarBytes
	if v, ok := resolveVarsDepth(raw, customs, 0, &budget); ok {
		return v
	}
	return raw
}

// resolveVarsDepth is resolveVars one definition deep, answering false once
// the budget every level shares is spent.
func resolveVarsDepth(raw string, customs map[string]string, depth int, budget *int) (string, bool) {
	if depth > 16 {
		return raw, true
	}
	out := raw
	for i := 0; i < len(out); i++ {
		if out[i] != 'v' || !strings.HasPrefix(out[i:], "var(") {
			continue
		}
		end := parenEnd(out, i+3)
		if end < 0 {
			return out, true
		}
		name, fallback, ok := splitVarArg(out[i+4 : end])
		if !ok {
			i = end
			continue
		}
		v, found := customs[name]
		if !found {
			v = fallback
		}
		if !found && fallback == "" {
			i = end
			continue
		}
		resolved, ok := resolveVarsDepth(v, customs, depth+1, budget)
		if !ok {
			return "", false
		}
		*budget -= len(resolved)
		if *budget < 0 {
			return "", false
		}
		out = out[:i] + resolved + out[end+1:]
	}
	return out, true
}
