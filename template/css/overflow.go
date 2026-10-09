package css

// Overflow is the per-axis overflow behaviour of a box.
const (
	OverflowVisible uint8 = iota
	OverflowHidden
	OverflowScroll
	OverflowAuto
	OverflowClip
)

func parseOverflow(raw string) (uint8, bool) {
	switch raw {
	case "visible":
		return OverflowVisible, true
	case "hidden":
		return OverflowHidden, true
	case "scroll":
		return OverflowScroll, true
	case "auto":
		return OverflowAuto, true
	case "clip":
		return OverflowClip, true
	}
	return 0, false
}
