package css

// parseFlexDirection turns "row", "column" and their reverses into a
// FlexDirection* constant.
func parseFlexDirection(raw string) (uint8, bool) {
	switch raw {
	case "row":
		return FlexDirectionRow, true
	case "row-reverse":
		return FlexDirectionRowReverse, true
	case "column":
		return FlexDirectionColumn, true
	case "column-reverse":
		return FlexDirectionColumnReverse, true
	}
	return 0, false
}
