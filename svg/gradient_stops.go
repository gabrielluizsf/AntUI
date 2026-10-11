package svg

import "github.com/gabrielluizsf/antui/canvas"

// gradientStops reads the `<stop>` children of a gradient. An offset is a number
// or a percentage of the line, held at 0 when the stop does not say — which is
// not the same as an offset of zero, and is why the two are told apart. A stop
// with no colour of its own is black, which is what SVG says rather than what
// looks right.
func gradientStops(kids []*element, warn func(string, ...any)) []gradientStop {
	var stops []gradientStop
	for _, k := range kids {
		if k.Name != "stop" {
			continue
		}
		s := gradientStop{}
		if raw, ok := k.Attr["offset"]; ok {
			v, ok := gradientLength(raw, 100)
			if !ok {
				warn("the stop offset %q is not a number, so it lands at the start of the gradient", raw)
			}
			s.offset, s.hasOffset = clampFloat(v, 0, 1), ok
		}
		if isCurrentColor(k.attr("stop-color")) {
			// The keyword is kept rather than resolved: the colour of the widget
			// is only known when the drawing is painted, and a gradient written in
			// it should follow the widget like any other paint does.
			s.current = true
		} else if c, ok := parsePaint(k.attr("stop-color")); ok {
			s.color = c
		} else {
			s.color = canvas.Black
		}
		if raw, ok := k.Attr["stop-opacity"]; ok {
			v, ok := parseAlpha(raw)
			if !ok {
				warn("the stop opacity %q is not a number, so the stop is fully opaque", raw)
				v = 1
			}
			if !s.current {
				s.color = fade(s.color, v)
			}
		}
		stops = append(stops, s)
	}
	return stops
}
