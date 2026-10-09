package css

import (
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
