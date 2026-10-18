package svg

// maxSpecificity is the heaviest weight among the selectors of a `:not`,
// `:is` or `:where` list, which is what the pseudo-class itself counts as.
func maxSpecificity(subs []complexSelector) (int, int, int) {
	var best [3]int
	for _, cs := range subs {
		a, b, c := specificityTriple(cs)
		if heavier(a, b, c, best) {
			best = [3]int{a, b, c}
		}
	}
	return best[0], best[1], best[2]
}

// maxSpecificityRelative is [maxSpecificity] for the relative selectors of a
// `:has` list.
func maxSpecificityRelative(rels []relativeSelector) (int, int, int) {
	var best [3]int
	for _, rel := range rels {
		a, b, c := specificityTriple(rel.cs)
		if heavier(a, b, c, best) {
			best = [3]int{a, b, c}
		}
	}
	return best[0], best[1], best[2]
}

// heavier says whether one weight beats another, comparing the rows in order.
func heavier(a, b, c int, best [3]int) bool {
	if a != best[0] {
		return a > best[0]
	}
	if b != best[1] {
		return b > best[1]
	}
	return c > best[2]
}
