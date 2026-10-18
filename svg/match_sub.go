package svg

// matchSub says whether any selector of a `:not`, `:is` or `:where` list
// stands on the element.
func matchSub(sub []complexSelector, n *selNode) bool {
	for _, cs := range sub {
		if matchComplex(cs, n) {
			return true
		}
	}
	return false
}
