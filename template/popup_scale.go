package template

// popupScale is the scale an open calendar draws at, read back out of the
// width it was handed: the popup is measured in [datePickerSize] steps, so its
// width says which step it was measured at. The styles are given the popup's
// box and nothing else, and this is how they learn the size it wants.
func popupScale(w int) int {
	base, _ := datePickerSize(1)
	if w < base {
		return 1
	}
	return w / base
}
