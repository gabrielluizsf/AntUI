package svg

// simpleSpecificity is how one simple selector counts on each of the three
// rows. Every pseudo-class counts as a class, a pseudo-element as a type, and
// `:where` counts for nothing; `:not`, `:is` and `:has` count the heaviest of
// their arguments.
func simpleSpecificity(s simpleSelector) (a, b, c int) {
	switch s.kind {
	case kindID:
		a = 1
	case kindClass, kindAttribute:
		b = 1
	case kindType:
		c = 1
	case kindUniversal, kindWhere:
		// nothing
	case kindNot, kindIs:
		a, b, c = maxSpecificity(s.sub)
	case kindHas:
		a, b, c = maxSpecificityRelative(s.has)
	case kindUnmatchable:
		if s.pe {
			c = 1
		} else {
			b = 1
		}
	default:
		b = 1
	}
	return a, b, c
}
