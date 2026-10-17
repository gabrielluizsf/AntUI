package svg

// matchCompound says whether every simple selector of a compound holds for the
// element at node.
func matchCompound(c compoundSelector, n *selNode) bool {
	for _, s := range c.simples {
		if !matchSimple(s, n) {
			return false
		}
	}
	return true
}
