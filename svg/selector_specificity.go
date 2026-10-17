package svg

// specificity weighs a selector against another: the number of ids, then of
// classes and attribute conditions, then of type names, each read as a digit
// of one number. The universal selector and the combinators count for nothing.
func specificity(cs complexSelector) int {
	a, b, c := 0, 0, 0
	for _, comp := range cs.compounds {
		for _, s := range comp.simples {
			switch s.kind {
			case kindID:
				a++
			case kindClass, kindAttribute:
				b++
			case kindType:
				c++
			}
		}
	}
	return a*10000 + b*100 + c
}
