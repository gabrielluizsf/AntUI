package css

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// CurrentColor is the sentinel a style field holds while its declaration said
// currentColor. The cascade resolves it to the computed color property once
// every declaration has applied, because a browser reads currentColor after
// the whole cascade rather than mid-way. The value happens to spell
// "transparent black" if a sheet literally writes rgba(0 0 0 / 0); the engine
// accepts that quirk.
const CurrentColor canvas.Color = 0x00000001

// namedColors is every CSS extended colour keyword, the full CSS Color 4
// table glued to the sixteen HTML names. Each maps to the exact sRGB value a
// browser paints, so red really is the red of the rainbow and tomato really
// is a tomato. transparent and currentColor are handled as keywords by
// [ParseColor], not by this table.
var namedColors = func() map[string]canvas.Color {
	m := make(map[string]canvas.Color, len(cssColorNames)+4)
	for _, c := range cssColorNames {
		var v uint64
		fmt.Sscanf(c.hex, "%x", &v)
		m[c.name] = canvas.Color(v | 0xFF000000)
	}
	return m
}()

// cssColorNames is name/hex for every CSS extended colour keyword. The grey
// spellings are aliases of the gray spellings and share their rows' value.
var cssColorNames = []struct{ name, hex string }{
	{"aliceblue", "F0F8FF"},
	{"antiquewhite", "FAEBD7"},
	{"aqua", "00FFFF"},
	{"aquamarine", "7FFFD4"},
	{"azure", "F0FFFF"},
	{"beige", "F5F5DC"},
	{"bisque", "FFE4C4"},
	{"black", "000000"},
	{"blanchedalmond", "FFEBCD"},
	{"blue", "0000FF"},
	{"blueviolet", "8A2BE2"},
	{"brown", "A52A2A"},
	{"burlywood", "DEB887"},
	{"cadetblue", "5F9EA0"},
	{"chartreuse", "7FFF00"},
	{"chocolate", "D2691E"},
	{"coral", "FF7F50"},
	{"cornflowerblue", "6495ED"},
	{"cornsilk", "FFF8DC"},
	{"crimson", "DC143C"},
	{"cyan", "00FFFF"},
	{"darkblue", "00008B"},
	{"darkcyan", "008B8B"},
	{"darkgoldenrod", "B8860B"},
	{"darkgray", "A9A9A9"},
	{"darkgreen", "006400"},
	{"darkgrey", "A9A9A9"},
	{"darkkhaki", "BDB76B"},
	{"darkmagenta", "8B008B"},
	{"darkolivegreen", "556B2F"},
	{"darkorange", "FF8C00"},
	{"darkorchid", "9932CC"},
	{"darkred", "8B0000"},
	{"darksalmon", "E9967A"},
	{"darkseagreen", "8FBC8F"},
	{"darkslateblue", "483D8B"},
	{"darkslategray", "2F4F4F"},
	{"darkslategrey", "2F4F4F"},
	{"darkturquoise", "00CED1"},
	{"darkviolet", "9400D3"},
	{"deeppink", "FF1493"},
	{"deepskyblue", "00BFFF"},
	{"dimgray", "696969"},
	{"dimgrey", "696969"},
	{"dodgerblue", "1E90FF"},
	{"firebrick", "B22222"},
	{"floralwhite", "FFFAF0"},
	{"forestgreen", "228B22"},
	{"fuchsia", "FF00FF"},
	{"gainsboro", "DCDCDC"},
	{"ghostwhite", "F8F8FF"},
	{"gold", "FFD700"},
	{"goldenrod", "DAA520"},
	{"gray", "808080"},
	{"green", "008000"},
	{"greenyellow", "ADFF2F"},
	{"grey", "808080"},
	{"honeydew", "F0FFF0"},
	{"hotpink", "FF69B4"},
	{"indianred", "CD5C5C"},
	{"indigo", "4B0082"},
	{"ivory", "FFFFF0"},
	{"khaki", "F0E68C"},
	{"lavender", "E6E6FA"},
	{"lavenderblush", "FFF0F5"},
	{"lawngreen", "7CFC00"},
	{"lemonchiffon", "FFFACD"},
	{"lightblue", "ADD8E6"},
	{"lightcoral", "F08080"},
	{"lightcyan", "E0FFFF"},
	{"lightgoldenrodyellow", "FAFAD2"},
	{"lightgray", "D3D3D3"},
	{"lightgreen", "90EE90"},
	{"lightgrey", "D3D3D3"},
	{"lightpink", "FFB6C1"},
	{"lightsalmon", "FFA07A"},
	{"lightseagreen", "20B2AA"},
	{"lightskyblue", "87CEFA"},
	{"lightslategray", "778899"},
	{"lightslategrey", "778899"},
	{"lightsteelblue", "B0C4DE"},
	{"lightyellow", "FFFFE0"},
	{"lime", "00FF00"},
	{"limegreen", "32CD32"},
	{"linen", "FAF0E6"},
	{"magenta", "FF00FF"},
	{"maroon", "800000"},
	{"mediumaquamarine", "66CDAA"},
	{"mediumblue", "0000CD"},
	{"mediumorchid", "BA55D3"},
	{"mediumpurple", "9370DB"},
	{"mediumseagreen", "3CB371"},
	{"mediumslateblue", "7B68EE"},
	{"mediumspringgreen", "00FA9A"},
	{"mediumturquoise", "48D1CC"},
	{"mediumvioletred", "C71585"},
	{"midnightblue", "191970"},
	{"mintcream", "F5FFFA"},
	{"mistyrose", "FFE4E1"},
	{"moccasin", "FFE4B5"},
	{"navajowhite", "FFDEAD"},
	{"navy", "000080"},
	{"oldlace", "FDF5E6"},
	{"olive", "808000"},
	{"olivedrab", "6B8E23"},
	{"orange", "FFA500"},
	{"orangered", "FF4500"},
	{"orchid", "DA70D6"},
	{"palegoldenrod", "EEE8AA"},
	{"palegreen", "98FB98"},
	{"paleturquoise", "AFEEEE"},
	{"palevioletred", "DB7093"},
	{"papayawhip", "FFEFD5"},
	{"peachpuff", "FFDAB9"},
	{"peru", "CD853F"},
	{"pink", "FFC0CB"},
	{"plum", "DDA0DD"},
	{"powderblue", "B0E0E6"},
	{"purple", "800080"},
	{"rebeccapurple", "663399"},
	{"red", "FF0000"},
	{"rosybrown", "BC8F8F"},
	{"royalblue", "4169E1"},
	{"saddlebrown", "8B4513"},
	{"salmon", "FA8072"},
	{"sandybrown", "F4A460"},
	{"seagreen", "2E8B57"},
	{"seashell", "FFF5EE"},
	{"sienna", "A0522D"},
	{"silver", "C0C0C0"},
	{"skyblue", "87CEEB"},
	{"slateblue", "6A5ACD"},
	{"slategray", "708090"},
	{"slategrey", "708090"},
	{"snow", "FFFAFA"},
	{"springgreen", "00FF7F"},
	{"steelblue", "4682B4"},
	{"tan", "D2B48C"},
	{"teal", "008080"},
	{"thistle", "D8BFD8"},
	{"tomato", "FF6347"},
	{"turquoise", "40E0D0"},
	{"violet", "EE82EE"},
	{"wheat", "F5DEB3"},
	{"white", "FFFFFF"},
	{"whitesmoke", "F5F5F5"},
	{"yellow", "FFFF00"},
	{"yellowgreen", "9ACD32"},
}

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

func parseHex(s string) (canvas.Color, error) {
	body := strings.TrimPrefix(s, "#")
	valid := len(body) == 3 || len(body) == 4 || len(body) == 6 || len(body) == 8
	if !valid {
		return 0, fmt.Errorf("css: %q is not a hex colour", s)
	}
	nibble := func(c byte) (byte, bool) {
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
	pair := func(a, b byte) (byte, bool) {
		ha, ok := nibble(a)
		if !ok {
			return 0, false
		}
		hb, ok := nibble(b)
		if !ok {
			return 0, false
		}
		return ha<<4 | hb, true
	}
	var r, g, b, a uint8 = 0, 0, 0, 0xFF
	var err bool
	switch len(body) {
	case 3, 4:
		var ok bool
		if r, ok = pair(body[0], body[0]); !ok {
			err = true
		}
		if g, ok = pair(body[1], body[1]); !ok {
			err = true
		}
		if b, ok = pair(body[2], body[2]); !ok {
			err = true
		}
		if len(body) == 4 {
			if a, ok = pair(body[3], body[3]); !ok {
				err = true
			}
		}
	case 6, 8:
		var ok bool
		if r, ok = pair(body[0], body[1]); !ok {
			err = true
		}
		if g, ok = pair(body[2], body[3]); !ok {
			err = true
		}
		if b, ok = pair(body[4], body[5]); !ok {
			err = true
		}
		if len(body) == 8 {
			if a, ok = pair(body[6], body[7]); !ok {
				err = true
			}
		}
	}
	if err {
		return 0, fmt.Errorf("css: %q is not a hex colour", s)
	}
	return canvas.RGBA(r, g, b, a), nil
}

// parseRGB reads rgb()/rgba() in both the legacy comma form and the modern
// space-separated form with an optional "/ alpha" divider. Channels are
// integers 0-255 or percentages; the alpha is a number 0-1 or a percentage.
func parseRGB(s string) (canvas.Color, error) {
	slots, alpha, hasAlpha, ok := splitColorArgs(s)
	if !ok {
		return 0, fmt.Errorf("css: %q is not an rgb() colour", s)
	}
	a := uint8(255)
	switch {
	case hasAlpha:
		var err error
		if a, err = parseColorAlpha(alpha); err != nil {
			return 0, err
		}
	case len(slots) == 4:
		var err error
		if a, err = parseColorAlpha(slots[3]); err != nil {
			return 0, err
		}
		slots = slots[:3]
	}
	if len(slots) != 3 {
		return 0, fmt.Errorf("css: %q has the wrong channel count", s)
	}
	ch := func(f string) (uint8, error) {
		f = strings.TrimSpace(f)
		if strings.HasSuffix(f, "%") {
			v, err := strconv.ParseFloat(strings.TrimSuffix(f, "%"), 64)
			if err != nil {
				return 0, err
			}
			return clampByte(v * 255 / 100), nil
		}
		v, err := strconv.ParseFloat(f, 64)
		if err != nil {
			return 0, err
		}
		return clampByte(v), nil
	}
	r, err := ch(slots[0])
	if err != nil {
		return 0, err
	}
	g, err := ch(slots[1])
	if err != nil {
		return 0, err
	}
	b, err := ch(slots[2])
	if err != nil {
		return 0, err
	}
	return canvas.RGBA(r, g, b, a), nil
}

// parseHSL reads hsl()/hsla() in either syntax. The hue is an angle unit
// (deg/rad/grad/turn, or a bare number meaning degrees); saturation and
// lightness are percentages.
func parseHSL(s string) (canvas.Color, error) {
	slots, alpha, hasAlpha, ok := splitColorArgs(s)
	if !ok {
		return 0, fmt.Errorf("css: %q is not an hsl() colour", s)
	}
	a := uint8(255)
	switch {
	case hasAlpha:
		var err error
		if a, err = parseColorAlpha(alpha); err != nil {
			return 0, err
		}
	case len(slots) == 4:
		var err error
		if a, err = parseColorAlpha(slots[3]); err != nil {
			return 0, err
		}
		slots = slots[:3]
	}
	if len(slots) != 3 {
		return 0, fmt.Errorf("css: %q has the wrong channel count", s)
	}
	h, err := parseAngle(slots[0])
	if err != nil {
		return 0, err
	}
	sp, err := parsePercent(slots[1])
	if err != nil {
		return 0, err
	}
	lp, err := parsePercent(slots[2])
	if err != nil {
		return 0, err
	}
	c := hslToRGB(h, sp, lp)
	return canvas.RGBA(c[0], c[1], c[2], a), nil
}

// parseHWB reads hwb(hue white blackness / alpha): the hue angle, then two
// percentages that mix white and black into the hue.
func parseHWB(s string) (canvas.Color, error) {
	slots, alpha, hasAlpha, ok := splitColorArgs(s)
	if !ok {
		return 0, fmt.Errorf("css: %q is not an hwb() colour", s)
	}
	if len(slots) != 3 {
		return 0, fmt.Errorf("css: %q has the wrong channel count", s)
	}
	h, err := parseAngle(slots[0])
	if err != nil {
		return 0, err
	}
	w, err := parsePercent(slots[1])
	if err != nil {
		return 0, err
	}
	b, err := parsePercent(slots[2])
	if err != nil {
		return 0, err
	}
	a := uint8(255)
	if hasAlpha {
		if a, err = parseColorAlpha(alpha); err != nil {
			return 0, err
		}
	}
	c := hwbToRGB(h, w, b)
	return canvas.RGBA(c[0], c[1], c[2], a), nil
}

// parseLab reads lab(L a b / alpha) and parseLCH/parseOKLab/parseOKLCH the
// spaces alongside it. The CIE spaces are mapped to sRGB and clamped, which
// is the honest canvas behaviour for out-of-gamut colours.
func parseLab(s string) (canvas.Color, error) {
	slots, alpha, hasAlpha, ok := splitColorArgs(s)
	if !ok {
		return 0, fmt.Errorf("css: %q is not a lab() colour", s)
	}
	if len(slots) != 3 {
		return 0, fmt.Errorf("css: %q has the wrong channel count", s)
	}
	l, err := parseLabLight(slots[0])
	if err != nil {
		return 0, err
	}
	a, err := parseOpponent(slots[1])
	if err != nil {
		return 0, err
	}
	b, err := parseOpponent(slots[2])
	if err != nil {
		return 0, err
	}
	x, y, z := labToXYZ(l, a, b)
	return finishColor(x, y, z, alpha, hasAlpha)
}

func parseLCH(s string) (canvas.Color, error) {
	slots, alpha, hasAlpha, ok := splitColorArgs(s)
	if !ok {
		return 0, fmt.Errorf("css: %q is not an lch() colour", s)
	}
	if len(slots) != 3 {
		return 0, fmt.Errorf("css: %q has the wrong channel count", s)
	}
	l, err := parseLabLight(slots[0])
	if err != nil {
		return 0, err
	}
	c, err := parseLCHChroma(slots[1])
	if err != nil {
		return 0, err
	}
	h, err := parseAngle(slots[2])
	if err != nil {
		return 0, err
	}
	a, b := polarToLab(c, h)
	x, y, z := labToXYZ(l, a, b)
	return finishColor(x, y, z, alpha, hasAlpha)
}

func parseOKLab(s string) (canvas.Color, error) {
	slots, alpha, hasAlpha, ok := splitColorArgs(s)
	if !ok {
		return 0, fmt.Errorf("css: %q is not an oklab() colour", s)
	}
	if len(slots) != 3 {
		return 0, fmt.Errorf("css: %q has the wrong channel count", s)
	}
	l, err := parseOKLabLight(slots[0])
	if err != nil {
		return 0, err
	}
	a, err := parseOKLabChroma(slots[1])
	if err != nil {
		return 0, err
	}
	b, err := parseOKLabChroma(slots[2])
	if err != nil {
		return 0, err
	}
	r, g, b := okLabToSRGB(l, a, b)
	return finishOK(r, g, b, alpha, hasAlpha)
}

func parseOKLCH(s string) (canvas.Color, error) {
	slots, alpha, hasAlpha, ok := splitColorArgs(s)
	if !ok {
		return 0, fmt.Errorf("css: %q is not an oklch() colour", s)
	}
	if len(slots) != 3 {
		return 0, fmt.Errorf("css: %q has the wrong channel count", s)
	}
	l, err := parseOKLabLight(slots[0])
	if err != nil {
		return 0, err
	}
	c, err := parseOKLCHChroma(slots[1])
	if err != nil {
		return 0, err
	}
	h, err := parseAngle(slots[2])
	if err != nil {
		return 0, err
	}
	a, b := polarToLab(c, h)
	r, g, b := okLabToSRGB(l, a, b)
	return finishOK(r, g, b, alpha, hasAlpha)
}

// parseColorSpace reads color(space r g b / alpha), mapping the predefined
// spaces this engine can honestly paint: srgb, srgb-linear, display-p3, xyz
// and xyz-d65, plus xyz-d50 through a Bradford adaptation.
func parseColorSpace(s string) (canvas.Color, error) {
	slots, alpha, hasAlpha, ok := splitColorArgs(s)
	if !ok {
		return 0, fmt.Errorf("css: %q is not a color() colour", s)
	}
	if len(slots) != 4 {
		return 0, fmt.Errorf("css: %q has the wrong channel count", s)
	}
	space := strings.ToLower(strings.TrimSpace(slots[0]))
	r, g, b, err := parseColorComponents(slots[1:])
	if err != nil {
		return 0, err
	}
	var lr, lg, lb float64
	switch space {
	case "srgb":
		return finishOK(r, g, b, alpha, hasAlpha)
	case "srgb-linear":
		lr, lg, lb = r, g, b
	default:
		var x, y, z float64
		switch space {
		case "display-p3":
			lr, lg, lb = displayP3ToSRGB(r, g, b)
		case "xyz", "xyz-d65":
			x, y, z = r, g, b
			lr, lg, lb = xyzToLinearSRGB(x, y, z)
		case "xyz-d50":
			lr, lg, lb = xyzD50ToSRGB(r, g, b)
		default:
			return 0, fmt.Errorf("css: unsupported colour space %q", space)
		}
	}
	return finishOK(linearToSRGB(lr), linearToSRGB(lg), linearToSRGB(lb), alpha, hasAlpha)
}

// splitColorArgs cuts a functional colour's body ("1 2 3 / 0.5",
// "255, 0, 0, 0.5") into channel slots and, when a "/" introduced it, the
// trailing alpha. The comma form spreads its slots across the legacy syntax;
// the spaces form is the modern CSS4 way. ok is false when the trailing ")"
// is missing or the body is empty.
func splitColorArgs(s string) (slots []string, alpha string, hasAlpha, ok bool) {
	open := strings.IndexByte(s, '(')
	if open < 0 || !strings.HasSuffix(s, ")") {
		return nil, "", false, false
	}
	body := s[open+1 : len(s)-1]
	head, tail, slash := splitSlash(body)
	if strings.Contains(head, ",") {
		slots = splitFields(head, ',')
	} else {
		slots = strings.Fields(head)
	}
	for i := range slots {
		slots[i] = strings.TrimSpace(slots[i])
	}
	if slash {
		return slots, strings.TrimSpace(tail), true, true
	}
	return slots, "", false, true
}

// splitSlash splits on the "/" that sits at the top level of a colour
// function, the divider between channels and alpha. A slash inside nested
// parentheses belongs to something else and is ignored.
func splitSlash(s string) (head, tail string, has bool) {
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
		case '/':
			if depth == 0 {
				return s[:i], s[i+1:], true
			}
		}
	}
	return s, "", false
}

// splitFields splits a comma-joined list, ignoring commas nested in
// parentheses.
func splitFields(s string, sep byte) []string {
	var out []string
	start, depth := 0, 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
		case sep:
			if depth == 0 {
				out = append(out, s[start:i])
				start = i + 1
			}
		}
	}
	out = append(out, s[start:])
	return out
}

// parseColorAlpha reads the "/ 0.5" or "/ 50%" part of a colour: a number
// 0-1 or a percentage, clamped and scaled to an alpha byte.
func parseColorAlpha(f string) (uint8, error) {
	f = strings.TrimSpace(f)
	if strings.HasSuffix(f, "%") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(f, "%"), 64)
		if err != nil {
			return 0, err
		}
		return uint8(clamp01(v/100)*255 + 0.5), nil
	}
	v, err := strconv.ParseFloat(f, 64)
	if err != nil {
		return 0, err
	}
	return uint8(clamp01(v)*255 + 0.5), nil
}

// parseAngle reads a hue: a bare number means degrees; deg, rad, grad and
// turn convert to degrees.
func parseAngle(f string) (float64, error) {
	f = strings.TrimSpace(f)
	for _, c := range []struct {
		name string
		mult float64
	}{
		{"turn", 360},
		{"rad", 180 / math.Pi},
		{"grad", 0.9},
		{"deg", 1},
	} {
		if strings.HasSuffix(f, c.name) {
			v, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(f, c.name)), 64)
			if err != nil {
				return 0, err
			}
			return v * c.mult, nil
		}
	}
	return strconv.ParseFloat(f, 64)
}

// parsePercent reads a percentage or a bare number, both scaled in
// percent-units.
func parsePercent(f string) (float64, error) {
	f = strings.TrimSpace(f)
	return strconv.ParseFloat(strings.TrimSuffix(f, "%"), 64)
}

// parseLabLight reads lab/lch lightness: 0-100 as a number or a percentage.
func parseLabLight(f string) (float64, error) {
	return parsePercent(f)
}

// parseOpponent reads lab a/b: numbers, or percentages where 100% is 125.
func parseOpponent(f string) (float64, error) {
	if strings.HasSuffix(strings.TrimSpace(f), "%") {
		v, err := parsePercent(f)
		return v * 1.25, err
	}
	return parsePercent(f)
}

// parseLCHChroma reads lch chroma: a number, or a percentage where 100% is
// 150 (the CIE reference saturation).
func parseLCHChroma(f string) (float64, error) {
	if strings.HasSuffix(strings.TrimSpace(f), "%") {
		v, err := parsePercent(f)
		return v * 1.5, err
	}
	return parsePercent(f)
}

// parseOKLabLight reads oklab/oklch lightness: 0-1, with 100% meaning 1.
func parseOKLabLight(f string) (float64, error) {
	if strings.HasSuffix(strings.TrimSpace(f), "%") {
		v, err := parsePercent(f)
		return v / 100, err
	}
	return parsePercent(f)
}

// parseOKLabChroma reads oklab a/b: a number, or a percentage where 100% is
// 0.4 (the OKLab reference extent).
func parseOKLabChroma(f string) (float64, error) {
	if strings.HasSuffix(strings.TrimSpace(f), "%") {
		v, err := parsePercent(f)
		return v * 0.4 / 100, err
	}
	return parsePercent(f)
}

// parseOKLCHChroma reads oklch chroma: a number, or a percentage where 100%
// is 0.4.
func parseOKLCHChroma(f string) (float64, error) {
	return parseOKLabChroma(f)
}

// parseColorComponents reads the three channels of color(): numbers 0-1 or
// percentages, scaled to 0-1.
func parseColorComponents(slots []string) (r, g, b float64, err error) {
	vals := [3]float64{}
	for i, f := range slots {
		f = strings.TrimSpace(f)
		if strings.HasSuffix(f, "%") {
			v, perr := strconv.ParseFloat(strings.TrimSuffix(f, "%"), 64)
			if perr != nil {
				return 0, 0, 0, perr
			}
			vals[i] = v / 100
			continue
		}
		v, perr := strconv.ParseFloat(f, 64)
		if perr != nil {
			return 0, 0, 0, perr
		}
		vals[i] = v
	}
	return vals[0], vals[1], vals[2], nil
}

// polarToLab converts polar lch coordinates into lab a/b.
func polarToLab(c, h float64) (a, b float64) {
	h = h * math.Pi / 180
	return c * math.Cos(h), c * math.Sin(h)
}

// finishColor turns CIE XYZ into a canvas colour through linear sRGB,
// honouring an optional alpha. Out-of-gamut components are clamped.
func finishColor(x, y, z float64, alpha string, hasAlpha bool) (canvas.Color, error) {
	r, g, b := xyzToLinearSRGB(x, y, z)
	return finishOK(linearToSRGB(r), linearToSRGB(g), linearToSRGB(b), alpha, hasAlpha)
}

// finishOK turns linear sRGB components into a canvas colour, applying the
// transfer curve, clamping and an optional alpha.
func finishOK(r, g, b float64, alpha string, hasAlpha bool) (canvas.Color, error) {
	var a uint8 = 255
	if hasAlpha {
		var err error
		if a, err = parseColorAlpha(alpha); err != nil {
			return 0, err
		}
	}
	return canvas.RGBA(toByte(r), toByte(g), toByte(b), a), nil
}

// hwbToRGB turns hue/whiteness/blackness percentages into sRGB. When white
// and black fill the circle a neutral grey emerges; otherwise the pure hue is
// scaled and mixed toward white and black.
func hwbToRGB(h, w, b float64) [3]uint8 {
	w /= 100
	b /= 100
	if w+b >= 1 {
		g := uint8(w/(w+b)*255 + 0.5)
		return [3]uint8{g, g, g}
	}
	// The hue at maximum chroma, then constricted by the white and black.
	hue := hslToRGB(h, 100, 50)
	f := 1 - w - b
	return [3]uint8{
		uint8(float64(hue[0])*f + w*255 + 0.5),
		uint8(float64(hue[1])*f + w*255 + 0.5),
		uint8(float64(hue[2])*f + w*255 + 0.5),
	}
}

func hslToRGB(h, s, l float64) [3]uint8 {
	h = h / 360
	s = s / 100
	l = l / 100
	var r, g, b float64
	if s == 0 {
		r, g, b = l, l, l
	} else {
		q := l + s - l*s
		if l < 0.5 {
			q = l * (1 + s)
		}
		p := 2*l - q
		hu := func(t float64) float64 {
			if t < 0 {
				t++
			}
			if t > 1 {
				t--
			}
			switch {
			case t < 1.0/6.0:
				return p + (q-p)*6*t
			case t < 1.0/2.0:
				return q
			case t < 2.0/3.0:
				return p + (q-p)*(2.0/3.0-t)*6
			}
			return p
		}
		r = hu(h + 1.0/3.0)
		g = hu(h)
		b = hu(h - 1.0/3.0)
	}
	return [3]uint8{uint8(r*255 + 0.5), uint8(g*255 + 0.5), uint8(b*255 + 0.5)}
}

// The colour-space conversions below move a colour from a wide-gamut space
// into the canvas's sRGB. Reference whites are D65: X=0.95047, Y=1, Z=1.08883.
const (
	colorEpsilon = 216.0 / 24389.0
	colorKappa   = 24389.0 / 27.0
)

// labToXYZ maps CIE Lab into CIE XYZ under the D65 white point.
func labToXYZ(l, a, b float64) (x, y, z float64) {
	fy := (l + 16) / 116
	fx := fy + a/500
	fz := fy - b/200
	if v := fx * fx * fx; v > colorEpsilon {
		x = v
	} else {
		x = (116*fx - 16) / colorKappa
	}
	if v := fy * fy * fy; v > colorEpsilon {
		y = v
	} else {
		y = (116*fy - 16) / colorKappa
	}
	if v := fz * fz * fz; v > colorEpsilon {
		z = v
	} else {
		z = (116*fz - 16) / colorKappa
	}
	return x * 0.95047, y, z * 1.08883
}

// xyzToLinearSRGB maps CIE XYZ into linear sRGB; negative or out-of-range
// components are clamped by the caller's byte conversion.
func xyzToLinearSRGB(x, y, z float64) (r, g, b float64) {
	return 3.2404542*x - 1.5371385*y - 0.4985314*z,
		-0.9692660*x + 1.8760108*y + 0.0415560*z,
		0.0556434*x - 0.2040259*y + 1.0572252*z
}

// linearToSRGB runs the sRGB transfer curve on one linear channel.
func linearToSRGB(v float64) float64 {
	if v <= 0.0031308 {
		return 12.92 * v
	}
	return 1.055*math.Pow(v, 1.0/2.4) - 0.055
}

// srgbToLinear inverts the sRGB transfer curve.
func srgbToLinear(v float64) float64 {
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

// okLabToSRGB maps OKLab into linear sRGB, the closed-form matrix from the
// CSS Color 4 spec.
func okLabToSRGB(l, a, c float64) (r, g, b float64) {
	l_, m_, s_ := l+0.3963377774*a+0.2158037573*c,
		l-0.1055613458*a-0.0638541728*c,
		l-0.0894841775*a-1.2914855480*c
	l3, m3, s3 := l_*l_*l_, m_*m_*m_, s_*s_*s_
	return 4.0767416621*l3 - 3.3077115913*m3 + 0.2309699292*s3,
		-1.2684380046*l3 + 2.6097574011*m3 - 0.3413193965*s3,
		-0.0041960863*l3 - 0.7034186147*m3 + 1.7076147010*s3
}

// displayP3ToSRGB maps a gamma-encoded display-p3 component triple into
// linear sRGB: linearise with the shared sRGB curve, cross the primaries into
// XYZ, and come back through the sRGB matrix.
func displayP3ToSRGB(r, g, b float64) (sr, sg, sb float64) {
	r, g, b = srgbToLinear(r), srgbToLinear(g), srgbToLinear(b)
	x := 0.4865709486482162*r + 0.2656676931690931*g + 0.1982172852343625*b
	y := 0.2289745640697488*r + 0.6917385218365064*g + 0.0792869140937449*b
	z := 0.0451133818589026*g + 1.0439443689009757*b
	return xyzToLinearSRGB(x, y, z)
}

// xyzD50ToSRGB adapts a D50-adapted XYZ triple to D65 and maps it into
// linear sRGB, via the Bradford cone-response matrix.
func xyzD50ToSRGB(x, y, z float64) (r, g, b float64) {
	x = 0.9554734527*x - 0.0230985363*y + 0.0632593087*z
	y = -0.0283697069*x + 1.0099954581*y + 0.0210413984*z
	z = 0.0123140013*x - 0.0205076964*y + 1.3303659366*z
	return xyzToLinearSRGB(x, y, z)
}

// clampByte rounds a 0-1 component to a byte.
func toByte(v float64) uint8 {
	return uint8(clamp01(v)*255 + 0.5)
}

// clampByte clamps a 0-255 channel to a byte.
func clampByte(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v + 0.5)
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
