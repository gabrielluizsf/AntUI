package svg

// specificity weighs a selector against another: the number of ids, then of
// classes, attribute conditions and pseudo-classes, then of type names, each
// read as a digit of one number. The universal selector and the combinators
// count for nothing.
func specificity(cs complexSelector) int {
	a, b, c := specificityTriple(cs)
	return a*10000 + b*100 + c
}

// specificityTriple is the three rows of a selector's weight, kept apart so
// the heaviest argument of `:not`/`:is`/`:has` can be picked before they are
// folded into one number.
func specificityTriple(cs complexSelector) (a, b, c int) {
	for _, comp := range cs.compounds {
		for _, s := range comp.simples {
			sa, sb, sc := simpleSpecificity(s)
			a, b, c = a+sa, b+sb, c+sc
		}
	}
	return a, b, c
}
