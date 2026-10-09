package css

// isRepeatKeyword reports whether a token belongs to the repeat vocabulary,
// including the repeat-x/repeat-y shorthands.
func isRepeatKeyword(s string) bool {
	if _, ok := backRepeatWord(s); ok {
		return true
	}
	return s == "repeat-x" || s == "repeat-y"
}
