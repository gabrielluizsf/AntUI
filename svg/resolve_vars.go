package svg

import "strings"

// resolveVars replaces every var(--name[, fallback]) in a raw value with the
// cascaded custom-property value, or the fallback, walking nested definitions
// up to a depth guard so a cycle --a: var(--a) cannot spin forever. A var()
// with no known value and no fallback is left where it is, for [Style.withVars]
// to drop the whole declaration.
func resolveVars(raw string, customs map[string]string) string {
	return resolveVarsDepth(raw, customs, 0)
}

// resolveVarsDepth is resolveVars one definition deep: a var() that resolves to
// another value carrying its own var() is walked again, and after the guard the
// value is handed back as it was rather than resolved for ever.
func resolveVarsDepth(raw string, customs map[string]string, depth int) string {
	if depth > 16 {
		return raw
	}
	var b strings.Builder
	for i := 0; i < len(raw); {
		if raw[i] == 'v' && strings.HasPrefix(raw[i:], "var(") {
			i = writeVar(&b, raw, i, customs, depth)
			continue
		}
		b.WriteByte(raw[i])
		i++
	}
	return b.String()
}

// writeVar writes what the var() standing at i resolves to and answers where
// the rest of the value carries on. A var() nothing answers is written as it
// was, so the declaration it belongs to is dropped further out rather than
// painted with a name it never resolved.
func writeVar(b *strings.Builder, raw string, i int, customs map[string]string, depth int) int {
	end := parenEnd(raw, i+3)
	if end < 0 {
		b.WriteString(raw[i:])
		return len(raw)
	}
	name, fallback, ok := splitVarArg(raw[i+4 : end])
	if ok {
		if v, found := customs[name]; found {
			b.WriteString(resolveVarsDepth(v, customs, depth+1))
			return end + 1
		}
		if fallback != "" {
			b.WriteString(resolveVarsDepth(fallback, customs, depth+1))
			return end + 1
		}
	}
	b.WriteString(raw[i : end+1])
	return end + 1
}
