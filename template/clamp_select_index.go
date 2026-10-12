package template

// clampSelectIndex keeps the selection inside the option list.
func clampSelectIndex(index *int, options []string) {
	if *index < 0 {
		*index = 0
	}
	if *index >= len(options) {
		*index = len(options) - 1
	}
}
