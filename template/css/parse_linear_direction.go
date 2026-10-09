package css

import (
	"strings"
)

import (
	"math"
)

// parseLinearDirection reads the direction a linear gradient heads toward: an
// angle or a "to" corner/side phrase, in radians measured clockwise from the
// top. Only an angle or a "to" form counts as a direction.
func parseLinearDirection(s string) (float64, bool) {
	switch s {
	case "to top":
		return 0, true
	case "to bottom":
		return math.Pi, true
	case "to left":
		return -math.Pi / 2, true
	case "to right":
		return math.Pi / 2, true
	case "to top left", "to left top":
		return -math.Pi / 4, true
	case "to top right", "to right top":
		return math.Pi / 4, true
	case "to bottom left", "to left bottom":
		return 3 * math.Pi / 4, true
	case "to bottom right", "to right bottom":
		return -3 * math.Pi / 4, true
	default:
		if strings.HasPrefix(s, "to ") {
			return 0, false
		}
		deg, err := parseAngle(s)
		if err != nil {
			return 0, false
		}
		return deg * math.Pi / 180, true
	}
}
