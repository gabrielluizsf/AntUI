package css

// parseTextTransform reads the text-transform keyword.
func parseTextTransform(raw string) (uint8, bool) {
	switch raw {
	case "none":
		return TextTransformNone, true
	case "uppercase":
		return TextTransformUppercase, true
	case "lowercase":
		return TextTransformLowercase, true
	case "capitalize":
		return TextTransformCapitalize, true
	}
	return 0, false
}
