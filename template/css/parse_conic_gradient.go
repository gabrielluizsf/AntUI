package css

import (
	"strings"
)

import (
	"math"
)

// parseConicGradient reads "from <angle>? at <position>?, <stops>".
// The start angle defaults to 0° (pointing up) and the centre to the middle.
func parseConicGradient(args string, ctx Units) (*Gradient, bool) {
	parts := splitFields(args, ',')
	if len(parts) == 0 {
		return nil, false
	}
	g := &Gradient{
		Kind:   GradientConic,
		Center: backMiddle(),
		Shape:  GradientEllipse,
		Size:   GradientFarthestCorner,
	}
	if _, err := ParseColor(strings.TrimSpace(parts[0])); err == nil {
		stops, ok := parseStops(parts, ctx)
		if !ok {
			return nil, false
		}
		g.Stops = stops
		return g, true
	}
	tokens := splitTokens(strings.TrimSpace(parts[0]))
	i := 0
	if i < len(tokens) && tokens[i] == "from" {
		if len(tokens) < 2 {
			return nil, false
		}
		deg, err := parseAngle(tokens[1])
		if err != nil {
			return nil, false
		}
		g.Angle = deg * math.Pi / 180
		i = 2
	}
	if i < len(tokens) && tokens[i] == "at" {
		pos, ok := parseBackPosTokens(tokens[i+1:], ctx)
		if !ok {
			return nil, false
		}
		g.Center = pos
		i = len(tokens)
	}
	stopParts := gradientStopParts(tokens, i, parts[1:])
	stops, ok := parseStops(stopParts, ctx)
	if !ok {
		return nil, false
	}
	g.Stops = stops
	return g, true
}
