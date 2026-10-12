package svg

import "github.com/gabrielluizsf/antui/canvas"

// parseHexColor reads `#rgb`, `#rgba`, `#rrggbb` and `#rrggbbaa`. The three and
// four digit forms repeat each digit, which is how `#f0a` is the same colour as
// `#ff00aa` and not a shorter way of writing something else.
func parseHexColor(body string) (canvas.Color, bool) {
	r, g, b, a, ok := hexBytes(body)
	if !ok {
		return 0, false
	}
	return canvas.RGBA(r, g, b, a), true
}
