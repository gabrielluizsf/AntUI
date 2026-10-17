package svg

// cssWinner is the declaration that won one property, with what it takes for
// another to beat it: importance first, then the weight of the selector it
// came from, then the order the rules were written in.
type cssWinner struct {
	value     string
	important bool
	spec      int
	order     int
}

// stronger says whether a declaration beats the one already in hand for the
// same property. An important declaration beats any normal one; between two of
// the same importance the heavier selector wins, and between two of the same
// weight the one written later.
func stronger(d cssDecl, spec, order int, cur cssWinner) bool {
	if d.important != cur.important {
		return d.important
	}
	if spec != cur.spec {
		return spec > cur.spec
	}
	return order > cur.order
}
