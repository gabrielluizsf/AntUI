package svg

import "strings"

// readBlend reads a `feBlend`: the blend mode, which is `normal` where none is
// written.
func (p *filterPrim) readBlend(e *element, warn func(string, ...any)) bool {
	mode := strings.ToLower(strings.TrimSpace(e.attr("mode")))
	if mode == "" {
		mode = "normal"
	}
	switch mode {
	case "normal":
		p.blend = blendNormal
	case "multiply":
		p.blend = blendMultiply
	case "screen":
		p.blend = blendScreen
	case "darken":
		p.blend = blendDarken
	case "lighten":
		p.blend = blendLighten
	default:
		warn("the feBlend mode %q is not one this package can read, so \"normal\" is used instead", mode)
	}
	return true
}
