package svg

import "strings"

// matchDecls answers the declarations the sheets give one element, or nil
// where no rule matches it. For every property the whole cascade of matching
// rules is weighed, the important and the heavier and the later winning.
func matchDecls(n *selNode, rules []cssRule) *sheetDecls {
	var best map[string]cssWinner
	for order, r := range rules {
		spec, matched := ruleWeight(r, n)
		if !matched {
			continue
		}
		for _, d := range r.decls {
			if strings.HasPrefix(d.name, "--") {
				continue // a custom property is read where it is used
			}
			cur, seen := best[d.name]
			if !seen || stronger(d, spec, order, cur) {
				if best == nil {
					best = map[string]cssWinner{}
				}
				best[d.name] = cssWinner{value: d.value, important: d.important, spec: spec, order: order}
			}
		}
	}
	return sheetMap(best)
}

// ruleWeight is the heaviest selector of a rule that matches the element, and
// whether any of them did.
func ruleWeight(r cssRule, n *selNode) (int, bool) {
	best, ok := 0, false
	for _, sel := range r.selectors {
		if !matchComplex(sel, n) {
			continue
		}
		if s := specificity(sel); !ok || s > best {
			best, ok = s, true
		}
	}
	return best, ok
}
