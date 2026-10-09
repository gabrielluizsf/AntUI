package css

// Position is the CSS positioning scheme. Absolute and fixed boxes are taken
// out of flow; relative and sticky boxes keep their flow slot and only shift
// where they are painted.
const (
	PositionStatic uint8 = iota
	PositionAbsolute
	PositionFixed
	PositionRelative
	PositionSticky
)

func parsePosition(raw string) (uint8, bool) {
	switch raw {
	case "static":
		return PositionStatic, true
	case "absolute":
		return PositionAbsolute, true
	case "fixed":
		return PositionFixed, true
	case "relative":
		return PositionRelative, true
	case "sticky":
		return PositionSticky, true
	}
	return 0, false
}
