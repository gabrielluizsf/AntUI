package css

// backRepeatWord maps one repeat keyword onto the BackRepeat constants.
func backRepeatWord(s string) (uint8, bool) {
	switch s {
	case "repeat":
		return BackRepeatRepeat, true
	case "no-repeat":
		return BackRepeatNoRepeat, true
	case "space":
		return BackRepeatSpace, true
	case "round":
		return BackRepeatRound, true
	}
	return 0, false
}
