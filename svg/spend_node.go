package svg

// spendNode counts one more node against [maxBuiltNodes] and says whether it
// may be built. The first node past the limit says why the rest of the drawing
// is missing; every one after it is left out in silence, since a million
// copies of the same warning would be the memory the limit is there to keep.
func (img *Image) spendNode(warn func(string, ...any)) bool {
	if img.built >= maxBuiltNodes {
		if !img.overBudget {
			img.overBudget = true
			warn("the drawing builds more than %d shapes, so the rest of it is left out", maxBuiltNodes)
		}
		return false
	}
	img.built++
	return true
}
