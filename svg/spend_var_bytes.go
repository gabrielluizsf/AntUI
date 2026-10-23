package svg

// spendVarBytes takes n bytes from what one value may still grow by, and says
// whether there were that many left.
func spendVarBytes(budget *int, n int) bool {
	*budget -= n
	return *budget >= 0
}
