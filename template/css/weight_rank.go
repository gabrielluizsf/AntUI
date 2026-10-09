package css

// weightRank is how soon CSS Fonts asks for w when want is the weight in
// hand: 0 is the exact weight, and the rest follow the order the search
// turns — down first below 500, up first at 500 or above, with 500 ahead of
// every weight below 400 when 400 is asked for.
func weightRank(want, w uint16) int {
	w = normWeight(w)
	if w == want {
		return 0
	}
	switch {
	case want == 400:
		if w == 500 {
			return 1
		}
		if w < want {
			return 1 + int((want-w)/100)
		}
		// Past the weights below 400, the climb starts at 600: 500 has
		// its own place already.
		return 4 + int((w-500)/100)
	case want < 400:
		if w < want {
			return int((want - w) / 100)
		}
		return int((want-100)/100) + int((w-want)/100)
	case want == 500:
		if w < want {
			return int((want - w) / 100)
		}
		return 5 + int((w-600)/100)
	default:
		if w > want {
			return int((w - want) / 100)
		}
		return int((maxWeight-want)/100) + int((want-w)/100)
	}
}

const maxWeight = 1000
