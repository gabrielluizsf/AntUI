package svg

import (
	"strconv"
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// startOffset is how far along the outline the writing begins. A percentage is
// a share of how long the outline is, which only the outline can answer, and
// any other length is a number of the drawing's own units.
func startOffset(e *element, along *canvas.Path, warn func(string, ...any)) float64 {
	raw := strings.TrimSpace(e.attr("startOffset"))
	if raw == "" {
		return 0
	}
	if body, ok := strings.CutSuffix(raw, "%"); ok {
		share, err := strconv.ParseFloat(body, 64)
		if err != nil {
			warn("the startOffset %q is not a number, so the writing starts at the beginning of the outline", raw)
			return 0
		}
		return along.Length() * share / 100
	}
	return readLength(warn, e, "startOffset", "startOffset")
}
