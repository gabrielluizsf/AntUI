package css

import (
	"strings"
)

import (
	"math"
)

// parseLinearGradient reads "45deg, red, blue" or "to top right, red, blue".
// A missing direction defaults to to bottom, the CSS initial.
func parseLinearGradient(args string, ctx Units) (*Gradient, bool) {
	parts := splitFields(args, ',')
	if len(parts) == 0 {
		return nil, false
	}
	angle := math.Pi
	start := 0
	if a, ok := parseLinearDirection(strings.TrimSpace(parts[0])); ok {
		angle = a
		start = 1
	}
	stops, ok := parseStops(parts[start:], ctx)
	if !ok {
		return nil, false
	}
	return &Gradient{
		Kind:   GradientLinear,
		Angle:  angle,
		Center: backMiddle(),
		Shape:  GradientEllipse,
		Size:   GradientFarthestCorner,
		Stops:  stops,
	}, true
}
