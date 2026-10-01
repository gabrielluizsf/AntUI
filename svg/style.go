package svg

import (
	"strconv"
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// with reads the attributes written on one element over the style it inherited
// from the one above, and answers the style its children inherit in turn. Only
// what the element actually says changes: an attribute that is not there leaves
// the inherited value alone, which is what lets a `<g fill="red">` colour
// everything inside it and a single `<rect fill="blue">` be the one blue thing.
//
// The warn is called for each attribute that is written but cannot be read, so
// that a drawing with a number that is not one says so rather than being painted
// as though it had never asked.
func (s Style) with(e *element, warn func(string, ...any)) Style {
	// A presentation attribute is also what a `style="…"` declaration says, and
	// the declaration wins where both name the same thing.
	attrs := map[string]string{}
	for k, v := range e.Attr {
		if k != "" && k != "style" {
			attrs[k] = v
		}
	}
	for k, v := range parseDeclarations(e.attr("style")) {
		attrs[k] = v
	}

	// dash is the pattern as it was written, and dashOffset where along it the
	// stroke starts; the two are read together once every other attribute has
	// been, because the order a map is walked in is not one to depend on.
	var dash, dashOffset string
	_, hasDash := attrs["stroke-dasharray"]
	_, hasOffset := attrs["stroke-dashoffset"]

	for name, raw := range attrs {
		switch name {
		case "fill":
			// A colour to paint with, or a reference to something elsewhere in
			// the drawing, or nothing at all: `none` says so, and so does a
			// paint this cannot read at all, which is not a reason to paint the
			// shape black.
			switch {
			case isCurrentColor(raw):
				s.Fill, s.HasFill, s.FillCurrent = 0, true, true
				s.fillRef, s.fillFallback, s.hasFillFallback, s.fillFallbackCurrent = "", 0, false, false
				s.fillGradient = nil
			default:
				if c, ok := parsePaint(raw); ok {
					s.Fill, s.HasFill, s.FillCurrent = c, true, false
					s.fillRef, s.fillFallback, s.hasFillFallback, s.fillFallbackCurrent = "", 0, false, false
					s.fillGradient = nil
				} else if id, ok := paintRef(raw); ok {
					// The reference is kept whole: only a caller holding the
					// whole drawing can say what `#name` is, and the colour
					// after it is the one painted where the reference turns out
					// to be nothing. Neither the fill nor its colour is decided
					// here — see [Image.resolvePaints].
					s.fillRef, s.fillGradient = id, nil
					s.fillFallback, s.fillFallbackCurrent, s.hasFillFallback = paintFallback(raw)
					s.Fill, s.HasFill, s.FillCurrent = 0, false, false
				} else {
					s.Fill, s.HasFill, s.FillCurrent = 0, false, false
					s.fillRef, s.fillFallback, s.hasFillFallback, s.fillFallbackCurrent = "", 0, false, false
					s.fillGradient = nil
				}
			}
		case "fill-rule":
			if strings.EqualFold(strings.TrimSpace(raw), "evenodd") {
				s.FillRule = canvas.FillEvenOdd
			} else {
				s.FillRule = canvas.FillNonZero
			}
		case "clip-path":
			// The shapes that cut this element are elsewhere in the drawing,
			// the same way a gradient is, so only the id is kept here: only a
			// caller holding the whole drawing can say what it is. See
			// [Image.resolveClip].
			switch id, ok := paintRef(raw); {
			case strings.EqualFold(strings.TrimSpace(raw), "none"):
				s.clipRef = ""
			case ok:
				s.clipRef = id
			default:
				s.clipRef = ""
				warn("the clip-path %q is not a reference to a <clipPath> this package can follow, so the element is drawn without it", raw)
			}
		case "clip-rule":
			// The rule the shapes of a clipPath are read with, which is
			// inherited like the rest of the style and means the same thing
			// here that fill-rule means for a fill.
			if strings.EqualFold(strings.TrimSpace(raw), "evenodd") {
				s.ClipRule = canvas.FillEvenOdd
			} else {
				s.ClipRule = canvas.FillNonZero
			}
		case "fill-opacity":
			if v, ok := parseAlpha(raw); ok {
				s.FillOpacity = v
			}
		case "stroke":
			switch {
			case isCurrentColor(raw):
				s.Stroke, s.HasStroke, s.StrokeCurrent = 0, true, true
				s.strokeRef, s.strokeFallback, s.hasStrokeFallback, s.strokeFallbackCurrent = "", 0, false, false
				s.strokeGradient = nil
			default:
				if c, ok := parsePaint(raw); ok {
					s.Stroke, s.HasStroke, s.StrokeCurrent = c, true, false
					s.strokeRef, s.strokeFallback, s.hasStrokeFallback, s.strokeFallbackCurrent = "", 0, false, false
					s.strokeGradient = nil
				} else if id, ok := paintRef(raw); ok {
					s.strokeRef, s.strokeGradient = id, nil
					s.strokeFallback, s.strokeFallbackCurrent, s.hasStrokeFallback = paintFallback(raw)
					s.Stroke, s.HasStroke, s.StrokeCurrent = 0, false, false
				} else {
					s.Stroke, s.HasStroke, s.StrokeCurrent = 0, false, false
					s.strokeRef, s.strokeFallback, s.hasStrokeFallback, s.strokeFallbackCurrent = "", 0, false, false
					s.strokeGradient = nil
				}
			}
		case "stroke-width":
			if v, ok := parseLength(raw); ok {
				s.Width = v
			} else if strings.TrimSpace(raw) != "" {
				warn("the stroke width %q is not a length, so the stroke keeps the width it had", raw)
			}
		case "stroke-opacity":
			if v, ok := parseAlpha(raw); ok {
				s.StrokeOpacity = v
			}
		case "opacity":
			if v, ok := parseAlpha(raw); ok {
				s.Opacity = v
			}
		case "stroke-linecap":
			switch strings.ToLower(strings.TrimSpace(raw)) {
			case "round":
				s.Cap = canvas.CapRound
			case "square":
				s.Cap = canvas.CapSquare
			default:
				s.Cap = canvas.CapButt
			}
		case "stroke-linejoin":
			switch strings.ToLower(strings.TrimSpace(raw)) {
			case "round":
				s.Join = canvas.JoinRound
			case "bevel":
				s.Join = canvas.JoinBevel
			default:
				s.Join = canvas.JoinMiter
			}
		case "stroke-miterlimit":
			if v, err := strconv.ParseFloat(strings.TrimSpace(raw), 64); err == nil && v > 0 {
				s.MiterLimit = v
			}
		case "stroke-dasharray":
			dash = raw
		case "stroke-dashoffset":
			// Read together with the pattern below, since it means nothing on its
			// own and the order the two are found in is not one to depend on.
			dashOffset = raw
		case "font-size":
			s.FontSize = readFontSize(s.FontSize, raw, warn)
		case "text-anchor":
			if a, ok := readAnchor(raw); ok {
				s.Anchor = a
			} else if strings.TrimSpace(raw) != "" {
				warn("the text-anchor %q is not start, middle or end, so the writing stays where the pen is", raw)
			}
		case "display", "visibility":
			if strings.EqualFold(strings.TrimSpace(raw), "none") ||
				strings.EqualFold(strings.TrimSpace(raw), "hidden") {
				s.Hidden = true
			}
		case "transform":
			m, ok := parseTransform(raw)
			switch {
			case ok:
				if s.HasTransform {
					// Transforms compose outwards: the one on the child happens
					// to the shape before the one above it.
					s.Transform = s.Transform.Mul(m)
				} else {
					s.Transform, s.HasTransform = m, true
				}
			case strings.TrimSpace(raw) == "":
			default:
				warn("the transform %q is not one this package can read, so the shape stays where it is", raw)
			}
		}
	}
	if hasDash {
		s.Dash = parseDashArray(dash)
	}
	// The offset moves the pattern along the line, and an offset on its own says
	// nothing, so it is only read where there is a pattern for it to move.
	if hasOffset && s.Dash != nil {
		if v, ok := parseLength(dashOffset); ok {
			s.Dash.Offset = v
		} else {
			warn("the dash offset %q is not a length, so the pattern starts where it would have", dashOffset)
		}
	}
	return s
}

// readFontSize is how tall the writing is: a length in the drawing's own units,
// a percentage of the size it already has, one of the sizes CSS gives a name
// to, or a step up or down from the size it already has. Anything else leaves
// that size alone, since a size this cannot read is not a size to write at.
func readFontSize(was float64, raw string, warn func(string, ...any)) float64 {
	if was <= 0 {
		was = defaultFontSize
	}
	s := strings.ToLower(strings.TrimSpace(raw))
	switch s {
	case "larger":
		return was * 1.2
	case "smaller":
		return was / 1.2
	}
	if named, ok := namedFontSizes[s]; ok {
		return named
	}
	if p, ok := strings.CutSuffix(s, "%"); ok {
		if v, err := strconv.ParseFloat(strings.TrimSpace(p), 64); err == nil && v > 0 {
			return was * v / 100
		}
		warn("the font size %q is not a percentage, so the writing keeps the size it had", raw)
		return was
	}
	if v, ok := parseLength(s); ok && v > 0 {
		return v
	}
	warn("the font size %q is not a length, so the writing keeps the size it had", raw)
	return was
}

// namedFontSizes are the sizes CSS gives a name to, as how many pixels tall the
// writing is at the size the names step through.
var namedFontSizes = map[string]float64{
	"xx-small": 9,
	"x-small":  10,
	"small":    13,
	"medium":   16,
	"large":    18,
	"x-large":  24,
	"xx-large": 32,
}

// readAnchor reads where along the pen a piece of writing hangs: at it, centred
// on it, or ending at it. The third of those is what makes a label sit to the
// left of the point it names rather than to the right.
func readAnchor(raw string) (TextAnchor, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "start":
		return AnchorStart, true
	case "middle":
		return AnchorMiddle, true
	case "end":
		return AnchorEnd, true
	}
	return AnchorStart, false
}

// parseDeclarations reads the inside of a `style="…"` attribute, which is a list
// of `name: value` pairs separated by semicolons. The last one of a name is the
// one that counts, as in CSS.
func parseDeclarations(s string) map[string]string {
	out := map[string]string{}
	for _, decl := range strings.Split(s, ";") {
		name, value, ok := strings.Cut(decl, ":")
		if !ok {
			continue
		}
		name = strings.ToLower(strings.TrimSpace(name))
		if name == "" {
			continue
		}
		out[name] = strings.TrimSpace(value)
	}
	return out
}

// parseAlpha reads an opacity, which is a number from zero to one, a percentage,
// or nothing at all, in which case it is one.
func parseAlpha(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	if p, ok := strings.CutSuffix(s, "%"); ok {
		v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
		if err != nil {
			return 0, false
		}
		return clampFloat(v/100, 0, 1), true
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return clampFloat(v, 0, 1), true
}

// parseDashArray reads a `stroke-dasharray`, which is a list of lengths and
// their gaps, or the word `none` for a stroke drawn whole. A list that is not a
// pattern draws no dashes rather than drawing something arbitrary, which is what
// a browser does with one it cannot read.
func parseDashArray(s string) *canvas.Dash {
	s = strings.TrimSpace(s)
	if s == "" || strings.EqualFold(s, "none") {
		return nil
	}
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t' || r == '\n' || r == '\r'
	})
	lengths := make([]float64, 0, len(fields))
	for _, f := range fields {
		v, ok := parseLength(f)
		if !ok {
			return nil
		}
		lengths = append(lengths, v)
	}
	d, ok := canvas.NewDash(lengths, 0)
	if !ok {
		return nil
	}
	return &d
}

// parseLength reads a length, which is a number with an optional unit. The units
// a drawing may use are the pixels and the unitless numbers of a viewBox, and
// the absolute ones, which are turned into pixels at ninety-six to the inch, as
// CSS measures them.
func parseLength(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	for _, unit := range lengthUnits {
		if body, ok := strings.CutSuffix(s, unit.suffix); ok {
			v, err := strconv.ParseFloat(body, 64)
			if err != nil {
				// It ends in something that only looks like this unit, so the
				// next one is worth trying before giving up on the number.
				continue
			}
			return v * unit.scale, true
		}
	}
	// No unit at all is a number of the drawing's own units, which is what a
	// viewBox is written in.
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

// lengthUnits are the units a length may be written in, as how many pixels one
// of them is. A unitless number is not in here: it is the drawing's own, and
// only the caller knows how big that is.
//
// The longest unit comes first, because one of them ends in another: `cm` and
// `mm` share their last letter, and `in` is the tail of nothing else here, but a
// number is read by the first unit it genuinely ends in rather than the first
// whose last letter happens to match.
var lengthUnits = []struct {
	suffix string
	scale  float64
}{
	{"px", 1},
	{"pt", 96.0 / 72},
	{"pc", 16},
	{"cm", 96 / 2.54},
	{"mm", 96 / 25.4},
	{"in", 96},
	{"q", 96 / 101.6},
}

// parseNumber reads one plain number, which is what the path data and the point
// lists are written in.
func parseNumber(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}
