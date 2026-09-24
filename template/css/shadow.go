package css

import (
	"strconv"
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// Shadow is one box-shadow or text-shadow value: an offset, a blur radius, an
// optional spread and colour, and whether it falls inside the box. Lengths are
// reference pixels; the template scales them with the window. A colour that
// was never written stays [CurrentColor] until the cascade resolves it.
type Shadow struct {
	X, Y   int
	Blur   int
	Spread int
	Color  canvas.Color
	Inset  bool
}

// Filter is one filter function. Kind and Amount carry the ones the canvas can
// apply; Drop is set for drop-shadow(), whose shape comes from the element
// rather than a matrix. A nil Drop means Kind and Amount are the filter.
type Filter struct {
	Kind   canvas.FilterKind
	Amount float64
	Drop   *Shadow
}

// parseShadows reads a comma-separated shadow list. allowInset admits the
// inset keyword (box-shadow) and max is the number of lengths the property
// takes: four for box-shadow, three for text-shadow. "none" clears the list.
func parseShadows(raw string, ctx Units, allowInset bool, max int) ([]Shadow, bool) {
	if raw == "none" {
		return nil, true
	}
	var out []Shadow
	for _, part := range splitFields(raw, ',') {
		if sh, ok := parseShadow(strings.TrimSpace(part), ctx, allowInset, max); ok {
			out = append(out, sh)
		}
	}
	return out, len(out) > 0
}

// parseShadow reads one shadow value: an optional inset, two to four lengths,
// and a colour in any position.
func parseShadow(part string, ctx Units, allowInset bool, max int) (Shadow, bool) {
	var lens []int
	sh := Shadow{Color: CurrentColor}
	for _, tk := range splitTokens(part) {
		if allowInset && tk == "inset" {
			sh.Inset = true
			continue
		}
		if c, err := ParseColor(tk); err == nil {
			sh.Color = c
			continue
		}
		l, err := parseLengthAt(tk, ctx)
		if err != nil || l.IsPct() || l.Auto() || l.None() {
			return Shadow{}, false
		}
		lens = append(lens, l.Resolve(ctx))
	}
	if len(lens) < 2 || len(lens) > max {
		return Shadow{}, false
	}
	sh.X, sh.Y = lens[0], lens[1]
	if len(lens) > 2 {
		sh.Blur = lens[2]
	}
	if len(lens) > 3 {
		sh.Spread = lens[3]
	}
	return sh, true
}

// parseFilters reads a space-separated filter list. "none" clears it.
func parseFilters(raw string, ctx Units) ([]Filter, bool) {
	if raw == "none" {
		return nil, true
	}
	var out []Filter
	for _, tk := range splitTokens(raw) {
		if f, ok := parseFilter(tk, ctx); ok {
			out = append(out, f)
		}
	}
	return out, len(out) > 0
}

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

// filterAmount reads a filter's argument: a percentage or a bare number,
// where the empty argument is the identity of 1.
func filterAmount(args string) float64 {
	if args == "" {
		return 1
	}
	if strings.HasSuffix(args, "%") {
		if v, err := strconv.ParseFloat(strings.TrimSuffix(args, "%"), 64); err == nil {
			return v / 100
		}
		return 1
	}
	v, err := strconv.ParseFloat(args, 64)
	if err != nil {
		return 1
	}
	return v
}

// splitTokens splits a value on top-level whitespace, keeping a parenthesised
// group such as rgb(0 0 0) or drop-shadow(0 1px 2px black) in one piece.
func splitTokens(s string) []string {
	var out []string
	depth, start := 0, -1
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case ' ', '\t', '\n', '\r':
			if depth == 0 {
				if start >= 0 {
					out = append(out, s[start:i])
					start = -1
				}
				continue
			}
		}
		if start < 0 {
			start = i
		}
	}
	if start >= 0 {
		out = append(out, s[start:])
	}
	return out
}
