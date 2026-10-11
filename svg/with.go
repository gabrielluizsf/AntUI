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

	// pathLength is of the path on the element it is written on, so it does not
	// travel: a group that declares how long its own shape is says nothing
	// about the shapes inside, and each of those starts from none of its own.
	s.pathLength = 0

	// vector-effect does not travel either, and for a plainer reason than the
	// length: the spec says it is not inherited, so what a `<g>` asks for is
	// about the `<g>` — which has no outline of its own to compute — and the
	// shapes inside start from the scaled stroke every drawing gets.
	s.NonScalingStroke = false

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
				s.fillGradient, s.fillPattern = nil, nil
			default:
				if c, ok := parsePaint(raw); ok {
					s.Fill, s.HasFill, s.FillCurrent = c, true, false
					s.fillRef, s.fillFallback, s.hasFillFallback, s.fillFallbackCurrent = "", 0, false, false
					s.fillGradient, s.fillPattern = nil, nil
				} else if id, ok := paintRef(raw); ok {
					// The reference is kept whole: only a caller holding the
					// whole drawing can say what `#name` is, and the colour
					// after it is the one painted where the reference turns out
					// to be nothing. Neither the fill nor its colour is decided
					// here — see [Image.resolvePaints].
					s.fillRef, s.fillGradient, s.fillPattern = id, nil, nil
					s.fillFallback, s.fillFallbackCurrent, s.hasFillFallback = paintFallback(raw)
					s.Fill, s.HasFill, s.FillCurrent = 0, false, false
				} else {
					s.Fill, s.HasFill, s.FillCurrent = 0, false, false
					s.fillRef, s.fillFallback, s.hasFillFallback, s.fillFallbackCurrent = "", 0, false, false
					s.fillGradient, s.fillPattern = nil, nil
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
		case "mask":
			// The picture that says how much of this element is there is
			// elsewhere in the drawing, the same way a clip is, so only the id
			// is kept here: only a caller holding the whole drawing can say
			// what it is. It is not inherited — see the `maskRef` field of
			// [Style] and [Image.maskNamed].
			switch id, ok := paintRef(raw); {
			case strings.EqualFold(strings.TrimSpace(raw), "none"):
				s.maskRef = ""
			case ok:
				s.maskRef = id
			default:
				s.maskRef = ""
				warn("the mask %q is not a reference to a <mask> this package can follow, so the element is drawn without it", raw)
			}
		case "filter":
			// What this element's picture is put through when it is drawn,
			// read as it is written: a function is read here, and what is not
			// one of them is said here. A reference to a `<filter>` keeps the
			// id it named for only a caller holding the whole drawing to
			// follow — see [Image.followFilterRefs]. It is not inherited —
			// see the `filters` field of [Style] and [applyFilters].
			s.filters = readFilters(raw, warn)
		case "marker-start", "marker-mid", "marker-end":
			// The markers drawn at the vertices of this element's shape, which
			// unlike the clip, the mask and the filter ARE inherited: a `<g
			// marker-end>` is what every shape under it marks its own last
			// vertex with, so the reference is kept for the whole drawing to
			// follow and what was inherited is dropped as this element writes
			// one of its own — `none` among them, which is how a shape turns
			// off the marker its group asked for. The three are read the same
			// way, and each drops only the one it writes: a shape that says
			// `marker-start="none"` still carries the marker-mid the group
			// above it asked for. See [Image.resolveMarkers].
			ref, def := &s.markerStartRef, &s.markerStart
			if name == "marker-mid" {
				ref, def = &s.markerMidRef, &s.markerMid
			} else if name == "marker-end" {
				ref, def = &s.markerEndRef, &s.markerEnd
			}
			*ref, *def = "", nil
			switch id, ok := paintRef(raw); {
			case strings.EqualFold(strings.TrimSpace(raw), "none"):
			case ok:
				*ref = id
			default:
				warn("the %s %q is not a reference to a <marker> this package can follow, so no marker is drawn at the vertices it names", name, raw)
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
				s.strokeGradient, s.strokePattern = nil, nil
			default:
				if c, ok := parsePaint(raw); ok {
					s.Stroke, s.HasStroke, s.StrokeCurrent = c, true, false
					s.strokeRef, s.strokeFallback, s.hasStrokeFallback, s.strokeFallbackCurrent = "", 0, false, false
					s.strokeGradient, s.strokePattern = nil, nil
				} else if id, ok := paintRef(raw); ok {
					s.strokeRef, s.strokeGradient, s.strokePattern = id, nil, nil
					s.strokeFallback, s.strokeFallbackCurrent, s.hasStrokeFallback = paintFallback(raw)
					s.Stroke, s.HasStroke, s.StrokeCurrent = 0, false, false
				} else {
					s.Stroke, s.HasStroke, s.StrokeCurrent = 0, false, false
					s.strokeRef, s.strokeFallback, s.hasStrokeFallback, s.strokeFallbackCurrent = "", 0, false, false
					s.strokeGradient, s.strokePattern = nil, nil
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
		case "vector-effect":
			// Whether the stroke is taken through everything that scaled the
			// shape — which is what SVG does by default — or drawn at the width
			// it says, in the pixels the drawing comes out at whatever scaled
			// it. Any other value is a vector effect this package does not do,
			// and saying so beats quietly painting a stroke nobody asked for.
			switch strings.ToLower(strings.TrimSpace(raw)) {
			case "none":
				s.NonScalingStroke = false
			case "non-scaling-stroke":
				s.NonScalingStroke = true
			default:
				s.NonScalingStroke = false
				if strings.TrimSpace(raw) != "" {
					warn("the vector-effect %q is not none or non-scaling-stroke, so the stroke is scaled the way SVG scales it", raw)
				}
			}
		case "stroke-dasharray":
			dash = raw
		case "stroke-dashoffset":
			// Read together with the pattern below, since it means nothing on its
			// own and the order the two are found in is not one to depend on.
			dashOffset = raw
		case "pathlength":
			// The length is of the path on this element alone, which is what
			// taking it off the inherited style above is about. An attribute
			// name is lowercased as it is read, so `pathLength` is found here
			// spelled the way every other name in this switch is.
			s.pathLength = readPathLength(raw, warn)
		case "font-size":
			s.FontSize = readFontSize(s.FontSize, raw, warn)
		case "text-anchor":
			if a, ok := readAnchor(raw); ok {
				s.Anchor = a
			} else if strings.TrimSpace(raw) != "" {
				warn("the text-anchor %q is not start, middle or end, so the writing stays where the pen is", raw)
			}
		case "xml:space":
			// Only two words mean anything here, and anything else leaves the
			// inherited answer alone rather than switching the writing off its
			// own whitespace by accident.
			switch strings.ToLower(strings.TrimSpace(raw)) {
			case "preserve":
				s.PreserveSpace = true
			case "default":
				s.PreserveSpace = false
			}
		case "font-family":
			// The list is taken whole, the way CSS takes one: the names stay
			// as the drawing wrote them, and which of them the canvas answers
			// to is a question for when the writing is drawn. See [Style].
			s.FontFamily = strings.TrimSpace(raw)
		case "font-weight":
			s.FontWeight = readFontWeight(s.FontWeight, raw, warn)
		case "letter-spacing":
			s.LetterSpacing = readLetterSpacing(raw, warn)
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
