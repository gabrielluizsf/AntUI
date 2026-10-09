package css

import (
	"fmt"
	"github.com/gabrielluizsf/antui/canvas"
	"strings"
)

// ParseColor reads a CSS colour: a name, the transparent or currentColor
// keywords, #RGB, #RGBA, #RRGGBB or #RRGGBBAA, or a functional colour —
// rgb()/rgba(), hsl()/hsla(), hwb(), lab(), lch(), oklab(), oklch() and
// color(). Both the legacy comma syntax and the modern space syntax with a
// "/  alpha" divider are accepted.
func ParseColor(raw string) (canvas.Color, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return 0, fmt.Errorf("css: empty colour")
	}
	lower := strings.ToLower(s)
	if c, ok := namedColors[lower]; ok {
		return c, nil
	}
	switch lower {
	case "transparent":
		return canvas.Transparent, nil
	case "currentcolor":
		return CurrentColor, nil
	}
	switch {
	case strings.HasPrefix(lower, "rgb("), strings.HasPrefix(lower, "rgba("):
		return parseRGB(lower)
	case strings.HasPrefix(lower, "hsl("), strings.HasPrefix(lower, "hsla("):
		return parseHSL(lower)
	case strings.HasPrefix(lower, "hwb("):
		return parseHWB(lower)
	case strings.HasPrefix(lower, "lab("):
		return parseLab(lower)
	case strings.HasPrefix(lower, "lch("):
		return parseLCH(lower)
	case strings.HasPrefix(lower, "oklab("):
		return parseOKLab(lower)
	case strings.HasPrefix(lower, "oklch("):
		return parseOKLCH(lower)
	case strings.HasPrefix(lower, "color("):
		return parseColorSpace(lower)
	}
	if strings.HasPrefix(s, "#") {
		return parseHex(s)
	}
	return 0, fmt.Errorf("css: unrecognised colour %q", raw)
}
