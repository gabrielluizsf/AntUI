package css

// parseWhiteSpace reads the white-space keyword.
func parseWhiteSpace(raw string) (uint8, bool) {
	switch raw {
	case "normal":
		return WhiteSpaceNormal, true
	case "nowrap":
		return WhiteSpaceNowrap, true
	case "pre":
		return WhiteSpacePre, true
	case "pre-wrap":
		return WhiteSpacePreWrap, true
	case "pre-line":
		return WhiteSpacePreLine, true
	}
	return 0, false
}
