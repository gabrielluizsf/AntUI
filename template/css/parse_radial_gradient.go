package css

import (
	"strings"
)

// parseRadialGradient reads the prelude of a radial gradient — the ending
// shape, the size and an optional "at <position>", each optional — then the
// colour stops. The centre defaults to the middle of the painting area.
func parseRadialGradient(args string, ctx Units) (*Gradient, bool) {
	parts := splitFields(args, ',')
	if len(parts) == 0 {
		return nil, false
	}
	g := &Gradient{
		Kind:   GradientRadial,
		Center: backMiddle(),
		Shape:  GradientEllipse,
		Size:   GradientFarthestCorner,
	}
	// A first part that is already a colour stop begins the list with no
	// prelude at all.
	if _, err := ParseColor(strings.TrimSpace(parts[0])); err == nil {
		stops, ok := parseStops(parts, ctx)
		if !ok {
			return nil, false
		}
		g.Stops = stops
		return g, true
	}
	tokens := splitTokens(strings.TrimSpace(parts[0]))
	i := readRadialShape(tokens, g)
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
