package template

func clampGridItem(value, minimum, maximum int) int {
	if maximum >= 0 {
		value = min(value, maximum)
	}
	if minimum >= 0 {
		value = max(value, minimum)
	}
	return value
}
