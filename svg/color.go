package svg

import (
	"fmt"
	"math"
	"strconv"
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

// paintRef is the id inside a `url(#name)` — the thing a fill or a stroke
// points at when it is not naming a colour. The bool says whether there was
// one: `url(#)` with nothing after the hash names nothing, and a paint this
// cannot read is a shape that has to say so rather than come out blank.
//
// A `url(...)` may carry a fallback after it, as in `url(#g) red`: the fallback
// is what the shape is painted with when the reference is not there, so it is
// returned alongside the id for the caller to keep.
func paintRef(s string) (id string, ok bool) {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(strings.ToLower(s), "url(") {
		return "", false
	}
	end := strings.Index(s, ")")
	if end < 0 {
		return "", false
	}
	inside := strings.TrimSpace(s[len("url("):end])
	// Quotes around the reference are allowed and mean nothing.
	inside = strings.Trim(inside, "\"'")
	if !strings.HasPrefix(inside, "#") {
		return "", false
	}
	id = inside[1:]
	return id, id != ""
}

// paintFallback is the colour after a `url(#g) red`, which is what the shape is
// painted with when the reference is not in the drawing. It answers whether
// there was one to give: `url(#g)` with nothing after it has none, and a
// fallback spelled `currentColor` comes back as the keyword rather than as a
// colour, because the colour of the widget is only known when it is painted.
func paintFallback(s string) (c canvas.Color, current, ok bool) {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(strings.ToLower(s), "url(") {
		return 0, false, false
	}
	end := strings.Index(s, ")")
	if end < 0 {
		return 0, false, false
	}
	rest := strings.TrimSpace(s[end+1:])
	if isCurrentColor(rest) {
		return 0, true, true
	}
	c, ok = parseColor(rest)
	return c, false, ok
}

// isCurrentColor says whether a paint is the keyword that stands for the colour
// the drawing is being painted in. The comparison is case-insensitive because
// CSS keywords are, and a drawing that spells it `currentcolor` means the same
// thing as one that does not.
func isCurrentColor(s string) bool {
	return strings.EqualFold(strings.TrimSpace(s), "currentcolor")
}

// parseColor reads a colour in any of the ways a drawing may write one: the hex
// forms from three to eight digits, the functional forms with numbers or
// percentages, and the name of one of the colours CSS knows. The bool says
// whether it was one this could read, so a caller can carry on with the rest of
// the drawing rather than stopping at a colour it has never heard of.
func parseColor(s string) (canvas.Color, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	if s[0] == '#' {
		return parseHexColor(s[1:])
	}
	if name, args, ok := strings.Cut(s, "("); ok && strings.HasSuffix(args, ")") {
		return parseFuncColor(name, args[:len(args)-1])
	}
	if c, ok := canvas.ParseColor(s); ok == nil {
		return c, true
	}
	if c, ok := namedColors[strings.ToLower(s)]; ok {
		return c, true
	}
	return 0, false
}

// parseHexColor reads `#rgb`, `#rgba`, `#rrggbb` and `#rrggbbaa`. The three and
// four digit forms repeat each digit, which is how `#f0a` is the same colour as
// `#ff00aa` and not a shorter way of writing something else.
func parseHexColor(body string) (canvas.Color, bool) {
	digit := func(c byte) (uint8, bool) {
		switch {
		case c >= '0' && c <= '9':
			return c - '0', true
		case c >= 'a' && c <= 'f':
			return c - 'a' + 10, true
		case c >= 'A' && c <= 'F':
			return c - 'A' + 10, true
		}
		return 0, false
	}
	// read takes a run of digits and widens it to the eight of a colour, each
	// digit standing for itself and each byte for both of its halves.
	read := func(s string) (r, g, b, a uint8, ok bool) {
		switch len(s) {
		case 3, 4:
			v := make([]uint8, 0, 4)
			for i := range len(s) {
				d, good := digit(s[i])
				if !good {
					return 0, 0, 0, 0, false
				}
				v = append(v, d*17) // #f is #ff
			}
			r, g, b = v[0], v[1], v[2]
			a = 0xFF
			if len(v) == 4 {
				a = v[3]
			}
			return r, g, b, a, true
		case 6, 8:
			v := make([]uint8, 0, 4)
			for i := 0; i < len(s); i += 2 {
				hi, ok1 := digit(s[i])
				lo, ok2 := digit(s[i+1])
				if !ok1 || !ok2 {
					return 0, 0, 0, 0, false
				}
				v = append(v, hi<<4|lo)
			}
			r, g, b = v[0], v[1], v[2]
			a = 0xFF
			if len(v) == 4 {
				a = v[3]
			}
			return r, g, b, a, true
		}
		return 0, 0, 0, 0, false
	}
	r, g, b, a, ok := read(body)
	if !ok {
		return 0, false
	}
	return canvas.RGBA(r, g, b, a), true
}

// parseFuncColor reads the functional forms: `rgb`, `rgba` and `hsl`, with the
// numbers separated by commas, by spaces, or by a slash before the alpha, as
// SVG 2 and CSS both write them.
func parseFuncColor(name, args string) (canvas.Color, bool) {
	// split the argument list into numbers and percentages, dropping the
	// punctuation that separates them, keeping each one's unit.
	fields := strings.FieldsFunc(args, func(r rune) bool {
		return r == ',' || r == '/' || r == ' ' || r == '\t' || r == '\n' || r == '\r'
	})
	if len(fields) < 3 || len(fields) > 4 {
		return 0, false
	}
	// channel reads one number that is either 0-255 or a percentage of it.
	channel := func(s string) (float64, bool) {
		if p, ok := strings.CutSuffix(s, "%"); ok {
			v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
			if err != nil {
				return 0, false
			}
			return v * 255 / 100, true
		}
		v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
		if err != nil {
			return 0, false
		}
		return v, true
	}
	alpha := func(s string) (uint8, bool) {
		if p, ok := strings.CutSuffix(s, "%"); ok {
			v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
			if err != nil {
				return 0, false
			}
			return byte(clampFloat(v*255/100, 0, 255) + 0.5), true
		}
		v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
		if err != nil {
			return 0, false
		}
		return byte(clampFloat(v*255, 0, 255) + 0.5), true
	}
	a := 255.0
	if len(fields) == 4 {
		v, ok := alpha(fields[3])
		if !ok {
			return 0, false
		}
		a = float64(v)
	}

	switch strings.ToLower(name) {
	case "rgb", "rgba":
		r, ok1 := channel(fields[0])
		g, ok2 := channel(fields[1])
		b, ok3 := channel(fields[2])
		if !ok1 || !ok2 || !ok3 {
			return 0, false
		}
		return canvas.RGBA(toByte(r), toByte(g), toByte(b), toByte(a)), true
	case "hsl", "hsla":
		h, err1 := parseHue(fields[0])
		sv, err2 := parsePercent(fields[1])
		lv, err3 := parsePercent(fields[2])
		if err1 != nil || err2 != nil || err3 != nil {
			return 0, false
		}
		r, g, b := hslToRGB(h, sv/100, lv/100)
		return canvas.RGBA(toByte(r*255), toByte(g*255), toByte(b*255), toByte(a)), true
	}
	return 0, false
}

// parseHue reads the hue of a colour, which is an angle: parseAngle is where it
// is read, and a hue is nothing more than one in degrees.
func parseHue(s string) (float64, error) { return parseAngle(s) }

// parseAngle reads an angle in degrees, written either as a plain number or
// with any of the units a drawing may use, which is what a hue, a rotation and
// a shear all are.
func parseAngle(s string) (float64, error) {
	s = strings.TrimSpace(s)
	// The longest unit is tried first, because one of them ends in another:
	// `grad` is `rad` with a `g` in front of it, and taking the short one would
	// leave a number with a letter stuck to the end of it.
	for _, unit := range []struct {
		suffix string
		scale  float64
	}{
		{"grad", 360.0 / 400},
		{"turn", 360},
		{"rad", 180 / math.Pi},
		{"deg", 1},
	} {
		body, ok := strings.CutSuffix(s, unit.suffix)
		if !ok {
			continue
		}
		v, err := strconv.ParseFloat(strings.TrimSpace(body), 64)
		if err != nil {
			// It ends in something that only looks like this unit, so try the
			// next one rather than giving up on the number.
			continue
		}
		return v * unit.scale, nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

// parsePercent reads a number that must be written with its sign, as the
// saturation and lightness of a colour are.
func parsePercent(s string) (float64, error) {
	body, ok := strings.CutSuffix(strings.TrimSpace(s), "%")
	if !ok {
		return 0, fmt.Errorf("antui/svg: %q is not a percentage", s)
	}
	return strconv.ParseFloat(strings.TrimSpace(body), 64)
}

// hslToRGB turns a hue in degrees with a saturation and a lightness, both from
// zero to one, into the three channels it describes.
func hslToRGB(h, s, l float64) (float64, float64, float64) {
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}
	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - c/2
	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}
	return r + m, g + m, b + m
}

// toByte rounds a channel to the byte it becomes, holding it inside 0-255 so a
// value outside that range is clamped rather than wrapped around.
func toByte(v float64) uint8 { return byte(clampFloat(v, 0, 255) + 0.5) }

// clampFloat holds v between lo and hi.
func clampFloat(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
