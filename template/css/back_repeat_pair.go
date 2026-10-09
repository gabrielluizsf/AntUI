package css

// backRepeatPair folds a keyword list into one X/Y repeat pair.
func backRepeatPair(s []string) (BackRepeat, bool) {
	switch len(s) {
	case 1:
		switch s[0] {
		case "repeat-x":
			return BackRepeat{BackRepeatRepeat, BackRepeatNoRepeat}, true
		case "repeat-y":
			return BackRepeat{BackRepeatNoRepeat, BackRepeatRepeat}, true
		}
		x, ok := backRepeatWord(s[0])
		if !ok {
			return BackRepeat{}, false
		}
		return BackRepeat{x, x}, true
	case 2:
		x, ok := backRepeatWord(s[0])
		if !ok {
			return BackRepeat{}, false
		}
		y, ok := backRepeatWord(s[1])
		if !ok {
			return BackRepeat{}, false
		}
		return BackRepeat{x, y}, true
	}
	return BackRepeat{}, false
}
