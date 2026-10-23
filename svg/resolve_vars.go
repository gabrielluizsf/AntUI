package svg

import "strings"

// resolveVars replaces every var(--name[, fallback]) in a raw value with the
// cascaded custom-property value, or the fallback, walking nested definitions
// up to a depth guard so a cycle --a: var(--a) cannot spin forever. A var()
// with no known value and no fallback is left where it is, for [Style.withVars]
// to drop the whole declaration.
//
// A value that would resolve to more than [maxVarBytes] is handed back as it
// was written, which drops it the same way: ten properties each naming the one
// before ten times is a few hundred bytes that read into gigabytes, and the
// CSS Variables spec has such a value be invalid rather than worked out.
func resolveVars(raw string, customs map[string]string) string {
	budget := maxVarBytes
	if v, ok := resolveVarsDepth(raw, customs, 0, &budget); ok {
		return v
	}
	return raw
}

// resolveVarsDepth is resolveVars one definition deep: a var() that resolves to
// another value carrying its own var() is walked again, and after the guard the
// value is handed back as it was rather than resolved for ever. It answers
// false once the budget every level shares is spent.
func resolveVarsDepth(raw string, customs map[string]string, depth int, budget *int) (string, bool) {
	if depth > 16 {
		return raw, spendVarBytes(budget, len(raw))
	}
	var b strings.Builder
	for i := 0; i < len(raw); {
		if raw[i] == 'v' && strings.HasPrefix(raw[i:], "var(") {
			next, ok := writeVar(&b, raw, i, customs, depth, budget)
			if !ok {
				return "", false
			}
			i = next
			continue
		}
		if !spendVarBytes(budget, 1) {
			return "", false
		}
		b.WriteByte(raw[i])
		i++
	}
	return b.String(), true
}

// writeVar writes what the var() standing at i resolves to and answers where
// the rest of the value carries on. A var() nothing answers is written as it
// was, so the declaration it belongs to is dropped further out rather than
// painted with a name it never resolved.
func writeVar(b *strings.Builder, raw string, i int, customs map[string]string, depth int, budget *int) (int, bool) {
	end := parenEnd(raw, i+3)
	if end < 0 {
		b.WriteString(raw[i:])
		return len(raw), spendVarBytes(budget, len(raw)-i)
	}
	name, fallback, ok := splitVarArg(raw[i+4 : end])
	if ok {
		v, found := customs[name]
		if !found {
			v = fallback
		}
		if found || fallback != "" {
			// What the nested value spends is spent by the call itself; the
			// write here only copies text that has already been paid for.
			resolved, ok := resolveVarsDepth(v, customs, depth+1, budget)
			if !ok {
				return 0, false
			}
			b.WriteString(resolved)
			return end + 1, true
		}
	}
	b.WriteString(raw[i : end+1])
	return end + 1, spendVarBytes(budget, end+1-i)
}
