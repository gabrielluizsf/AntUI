package css

// parseFlexWrap turns "nowrap", "wrap" and "wrap-reverse" into a FlexWrap*
// constant.
func parseFlexWrap(raw string) (uint8, bool) {
	switch raw {
	case "nowrap":
		return FlexWrapNowrap, true
	case "wrap":
		return FlexWrapWrap, true
	case "wrap-reverse":
		return FlexWrapWrapReverse, true
	}
	return 0, false
}
