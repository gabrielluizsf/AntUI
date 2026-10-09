package css

// parseFontStyle reads the font-style keyword.
func parseFontStyle(raw string) (uint8, bool) {
	switch raw {
	case "normal":
		return FontStyleNormal, true
	case "italic":
		return FontStyleItalic, true
	case "oblique":
		return FontStyleOblique, true
	}
	return 0, false
}
