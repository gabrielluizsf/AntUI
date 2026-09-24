package css

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// parseBackgroundImage reads a comma-separated background-image list: url()
// names the template looks up in its image registry, the gradient functions
// paint themselves, and "none" clears the layers. A layer the engine does not
// recognise drops the whole property, exactly as an unknown value does.
func parseBackgroundImage(raw string, ctx Units) ([]BackImage, bool) {
	if raw == "none" {
		return nil, true
	}
	var out []BackImage
	for _, part := range splitFields(raw, ',') {
		img, ok := parseBackImage(strings.TrimSpace(part), ctx)
		if !ok {
			return nil, false
		}
		out = append(out, img)
	}
	if len(out) == 0 {
		return nil, false
	}
	return out, true
}

// parseBackImage reads one layer: a url() or one of the gradient functions.
func parseBackImage(part string, ctx Units) (BackImage, bool) {
	open := strings.IndexByte(part, '(')
	if open < 0 || !strings.HasSuffix(part, ")") {
		return BackImage{}, false
	}
	name := strings.TrimSpace(part[:open])
	args := strings.TrimSpace(part[open+1 : len(part)-1])
	switch name {
	case "url":
		u := strings.TrimSpace(strings.Trim(args, `"'`))
		if u == "" {
			return BackImage{}, false
		}
		return BackImage{URL: u}, true
	case "linear-gradient", "radial-gradient", "conic-gradient":
		g, ok := parseGradient(name, args, ctx)
		return BackImage{Grad: g}, ok
	}
	return BackImage{}, false
}

// parseGradient reads the arguments of one gradient function into a Gradient.
func parseGradient(name, args string, ctx Units) (*Gradient, bool) {
	switch name {
	case "linear-gradient":
		return parseLinearGradient(args, ctx)
	case "radial-gradient":
		return parseRadialGradient(args, ctx)
	default:
		return parseConicGradient(args, ctx)
	}
}

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
	i := 0
	for i < len(tokens) {
		switch tokens[i] {
		case "circle":
			g.Shape = GradientCircle
		case "ellipse":
			g.Shape = GradientEllipse
		case "closest-side":
			g.Size = GradientClosestSide
		case "farthest-side":
			g.Size = GradientFarthestSide
		case "closest-corner":
			g.Size = GradientClosestCorner
		case "farthest-corner":
			g.Size = GradientFarthestCorner
		default:
			goto shapesDone
		}
		i++
	}
shapesDone:
	if i < len(tokens) && tokens[i] == "at" {
		pos, ok := parseBackPosTokens(tokens[i+1:], ctx)
		if !ok {
			return nil, false
		}
		g.Center = pos
		i = len(tokens)
	}
	stopParts := parts[1:]
	if i < len(tokens) {
		// Words that never became prelude belong to the first colour stop,
		// which was missing its comma.
		stopParts = append([]string{strings.Join(tokens[i:], " ")}, parts[1:]...)
	}
	stops, ok := parseStops(stopParts, ctx)
	if !ok {
		return nil, false
	}
	g.Stops = stops
	return g, true
}

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
	stopParts := parts[1:]
	if i < len(tokens) {
		stopParts = append([]string{strings.Join(tokens[i:], " ")}, parts[1:]...)
	}
	stops, ok := parseStops(stopParts, ctx)
	if !ok {
		return nil, false
	}
	g.Stops = stops
	return g, true
}

// parseStops reads the colour stops of a gradient. A bare percentage or
// length between stops is a colour hint and is skipped, the way a browser
// reads the hint while keeping the interpolation the two neighbours ask for.
func parseStops(parts []string, ctx Units) ([]Stop, bool) {
	var out []Stop
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if isStopHint(p, ctx) {
			continue
		}
		stop, ok := parseStop(p, ctx)
		if !ok {
			return nil, false
		}
		out = append(out, stop)
	}
	return out, len(out) > 0
}

// isStopHint reports whether a stop is a bare offset with no colour, the
// colour-hint form the engine skips.
func isStopHint(p string, ctx Units) bool {
	if _, err := parseStopOffset(p, ctx); err != nil {
		return false
	}
	if _, err := ParseColor(p); err == nil {
		return false
	}
	return true
}

// parseStop reads "red", "red 50%" or "red 10% 20%": a colour, then up to two
// offsets. Two offsets pin the colour over a span, so the painter spreads it
// from the earlier one.
func parseStop(p string, ctx Units) (Stop, bool) {
	tokens := splitTokens(p)
	if len(tokens) == 0 {
		return Stop{}, false
	}
	var first, second float64 = -1, -1
	for n := len(tokens); n > 0; n-- {
		off, err := parseStopOffset(tokens[n-1], ctx)
		if err != nil {
			break
		}
		if second < 0 {
			second = off
		} else {
			first = off
		}
		tokens = tokens[:n-1]
	}
	c, err := ParseColor(strings.Join(tokens, " "))
	if err != nil {
		return Stop{}, false
	}
	off := first
	if second >= 0 && (first < 0 || second < first) {
		off = second
	}
	return Stop{Offset: off, Color: c}, true
}

// parseStopOffset reads one stop position: a percentage or a bare number,
// both in fractions, or a length, which is left to the percentage scale since
// the box is not known yet. Angles and the auto/none keywords do not belong.
func parseStopOffset(s string, ctx Units) (float64, error) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "%") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
		return v / 100, err
	}
	l, err := parseLengthAt(s, ctx)
	if err != nil {
		return 0, err
	}
	if l.IsPct() {
		return l.value / 100, nil
	}
	return 0, fmt.Errorf("css: unsupported gradient stop position %q", s)
}
