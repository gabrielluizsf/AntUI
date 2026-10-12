package svg

// beingDrawn reports whether pt is already being drawn further up: following it
// again would draw it while it is being drawn, for ever. A pattern that names
// itself, or two that name each other, is left off so the drawing comes out as
// a picture rather than not at all.
func (pt *pattern) beingDrawn(patterning []*pattern) bool {
	for _, being := range patterning {
		if being == pt {
			return true
		}
	}
	return false
}
