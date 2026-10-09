package css

import (
	"github.com/gabrielluizsf/antui/canvas"
	"strings"
)

// parseFilter reads one function such as blur(4px) or drop-shadow(0 2px 3px
// black).
func parseFilter(tk string, ctx Units) (Filter, bool) {
	open := strings.IndexByte(tk, '(')
	if open < 0 || !strings.HasSuffix(tk, ")") {
		return Filter{}, false
	}
	name := strings.TrimSpace(tk[:open])
	args := strings.TrimSpace(tk[open+1 : len(tk)-1])
	switch name {
	case "grayscale":
		return Filter{Kind: canvas.FilterGrayscale, Amount: filterAmount(args)}, true
	case "sepia":
		return Filter{Kind: canvas.FilterSepia, Amount: filterAmount(args)}, true
	case "invert":
		return Filter{Kind: canvas.FilterInvert, Amount: filterAmount(args)}, true
	case "brightness":
		return Filter{Kind: canvas.FilterBrightness, Amount: filterAmount(args)}, true
	case "contrast":
		return Filter{Kind: canvas.FilterContrast, Amount: filterAmount(args)}, true
	case "hue-rotate":
		deg, err := parseAngle(args)
		if err != nil {
			return Filter{}, false
		}
		return Filter{Kind: canvas.FilterHueRotate, Amount: deg}, true
	case "blur":
		l, err := parseLengthAt(args, ctx)
		if err != nil || l.IsPct() {
			return Filter{}, false
		}
		return Filter{Kind: canvas.FilterBlur, Amount: float64(l.Resolve(ctx))}, true
	case "drop-shadow":
		sh, ok := parseShadow(args, ctx, false, 3)
		if !ok {
			return Filter{}, false
		}
		return Filter{Drop: &sh}, true
	}
	return Filter{}, false
}
