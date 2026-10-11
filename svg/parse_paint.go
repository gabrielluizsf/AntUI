package svg

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// parsePaint reads one of the values a fill or a stroke can be given: a colour,
// the word `none` for nothing at all, or `url(#name)` for something taken from
// elsewhere in the drawing.
//
// The bool it answers says whether there is a colour here to paint with. A paint
// that names a gradient answers false, because it is not a colour — it is the
// reference itself that is wanted, and [Style.with] keeps that so the caller who
// can see the whole drawing can follow it to the gradient it names.
func parsePaint(s string) (canvas.Color, bool) {
	s = strings.TrimSpace(s)
	if s == "" || s == "none" {
		return 0, false
	}
	if strings.HasPrefix(s, "url(") {
		return 0, false
	}
	// `currentColor` is a colour, but not one the file says: it says "whatever
	// colour this is being drawn in", which is only known at painting time. It
	// reads as a colour here so the shape is not left unpainted, and the flag on
	// the style is what makes the painter use the colour it was handed instead.
	if isCurrentColor(s) {
		return 0, false
	}
	c, ok := parseColor(s)
	return c, ok
}
