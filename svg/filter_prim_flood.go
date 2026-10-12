package svg

import "github.com/gabrielluizsf/antui/canvas"

// readFlood reads a `feFlood`: the colour it floods with and how opaque, both
// defaulting to opaque black.
func (p *filterPrim) readFlood(e *element, warn func(string, ...any)) bool {
	p.flood = canvas.Color(0xFF000000)
	p.floodOpacity = 1
	if raw := e.attr("flood-color"); raw != "" {
		if isCurrentColor(raw) {
			p.floodCurrent = true
		} else if c, ok := parsePaint(raw); ok {
			p.flood = c
		} else {
			warn("the feFlood flood-color %q is not a colour this package can read, so black is flooded instead", raw)
		}
	}
	if raw := e.attr("flood-opacity"); raw != "" {
		if v, ok := filterAmount(raw); ok {
			p.floodOpacity = v
		} else {
			warn("the feFlood flood-opacity %q is not a number this package can read, so it is left out", raw)
		}
	}
	return true
}
