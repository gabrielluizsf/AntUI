package svg

// using reports whether an id is already being drawn through a `<use>` on the
// way to this one. A reference that names one of them goes round for ever
// rather than drawing a shape, so it is stopped at the first repeat.
func (img *Image) using(id string) bool {
	for _, seen := range img.uses {
		if seen == id {
			return true
		}
	}
	return false
}
