package svg

import (
	"math"
	"strconv"
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// A gradient is a paint rather than a colour: `fill="url(#g)"` names one
// elsewhere in the drawing, and the shape is painted by asking it what colour it
// is at each point rather than being painted one flat colour.
//
// The geometry is kept in the coordinates the gradient was written in and is not
// turned into a matrix at parse time, because SVG's own model is worth keeping:
// each number is a length in either user units or fractions of the shape's own
// box, and the two differ by whether x is divided by the box's width before it is
// used. A gradient written in a unit square and stretched to the shape it fills is
// exactly what a circular gradient in fractions of the bounding box is.
type gradient struct {
	// radial says a radial gradient; the others are linear. The kind is fixed by
	// the tag, so it is never inherited from the gradient an href points at.
	radial bool
	// userSpace says the geometry below is in user units rather than in fractions
	// of the shape's own bounding box.
	userSpace bool
	// stops are the colour changes along the gradient line, in the order they
	// were written.
	stops []gradientStop
	// transform is gradientTransform: it moves the gradient itself without
	// moving the shape it paints.
	transform canvas.Matrix
	hasXform  bool
	// spread is what happens past either end of the line: clamped, mirrored or
	// repeated.
	spread spread
	// set is which attributes the element actually wrote, which is what an href
	// consults: everything absent here is taken from the gradient it points at.
	set map[string]bool
	// href is the id of another gradient this one takes what it does not say from.
	href string
	// pct holds the attributes that were written as percentages, kept as they
	// were written until the units are settled. Which units they are is only
	// known once an href has been followed, and against a percentage the two are
	// worlds apart: a fraction of the shape, or a fraction of the drawing.
	pct map[string]pctRead

	// The linear line, from x1,y1 to x2,y2.
	x1, y1, x2, y2 float64
	// The radial circle, centred at cx,cy with radius r, and the focus fx,fy the
	// gradient is brightest at.
	cx, cy, r float64
	fx, fy    float64
	// viewport is the width and height a percentage written in user units is a
	// fraction of, which is the drawing's own and nothing to do with the shape
	// the gradient paints.
	viewport [2]float64
}

// pctRead is one number written as a percentage, kept until the units it is a
// percentage of are known: what was written, and which axis of the drawing it is
// measured along. The axis is -1 for the one a percentage has no side to it,
// which is a radius.
type pctRead struct {
	raw  string
	axis int
}

// spread is what a gradient does past the ends of its line.
type spread uint8

const (
	spreadPad spread = iota
	spreadReflect
	spreadRepeat
)

// gradientStop is one colour change: a colour and where along the line it
// happens, from 0 to 1. A stop written as `currentColor` keeps the keyword
// rather than a colour, because the colour of the widget is only known at
// painting time.
type gradientStop struct {
	offset    float64
	hasOffset bool
	color     canvas.Color
	current   bool
}

// gradientDefaults is what a gradient is before it says anything, and what it
// falls back to for every attribute it leaves out. The default is a linear
// gradient running left to right across the shape and, in fractions of that shape,
// a circle at the middle of it half as wide as it is.
func gradientDefaults(radial bool) *gradient {
	g := &gradient{radial: radial, spread: spreadPad, set: map[string]bool{}, pct: map[string]pctRead{}}
	if radial {
		g.cx, g.cy, g.r = 0.5, 0.5, 0.5
		return g
	}
	g.x2 = 1
	return g
}

// parseGradient reads one `<linearGradient>` or `<radialGradient>`. What it does
// not say is not decided here: an href names another gradient, and only a caller
// holding the whole drawing can say what that one had, so the reference is kept
// and followed by [gradient.resolved].
func parseGradient(e *element, warn func(string, ...any)) *gradient {
	g := gradientDefaults(e.Name == "radialGradient")
	g.href = gradientHref(e.Attr)
	g.stops = gradientStops(e.Kids, warn)
	g.read(e.Attr, warn)
	return g
}

// read takes the attributes written on the gradient. The defaults are already in
// place, so an attribute that is not there is left alone: that is what lets an
// href decide it afterwards.
func (g *gradient) read(attrs map[string]string, warn func(string, ...any)) {
	units := strings.TrimSpace(attrs["gradientunits"])
	if units != "" {
		g.userSpace = strings.EqualFold(units, "userspaceonuse")
		g.set["gradientunits"] = true
	}
	if v := strings.TrimSpace(attrs["spreadmethod"]); v != "" {
		g.set["spreadmethod"] = true
		switch strings.ToLower(v) {
		case "reflect":
			g.spread = spreadReflect
		case "repeat":
			g.spread = spreadRepeat
		default:
			// An unreadable spread method is the default one, which is what the
			// browser draws rather than leaving the ends of the line blank.
			g.spread = spreadPad
		}
	}
	if v := strings.TrimSpace(attrs["gradienttransform"]); v != "" {
		g.set["gradienttransform"] = true
		if m, ok := parseTransform(v); ok {
			g.transform, g.hasXform = m, true
		} else {
			warn("the gradient transform %q is not one this package can read, so the gradient is not moved", v)
		}
	}
	num := func(name string, dst *float64, axis int) {
		raw, ok := attrs[name]
		if !ok {
			return
		}
		g.set[name] = true
		if body, isPct := strings.CutSuffix(strings.TrimSpace(raw), "%"); isPct {
			// A percentage is kept as it was written, because which of the two
			// things it is a percentage of is decided by the units, and the units
			// may still be coming from an href that has not been followed yet.
			// Only that it is a number is worth saying now.
			if _, err := strconv.ParseFloat(strings.TrimSpace(body), 64); err != nil {
				warn("%s %q is not a number, so it stays at %g", name, raw, *dst)
				return
			}
			g.pct[name] = pctRead{raw: strings.TrimSpace(raw), axis: axis}
			return
		}
		v, ok := gradientLength(raw, 0)
		if !ok {
			warn("%s %q is not a length, so it stays at %g", name, raw, *dst)
			return
		}
		*dst = v
	}
	if g.radial {
		num("cx", &g.cx, 0)
		num("cy", &g.cy, 1)
		num("r", &g.r, -1)
		// The focus follows the centre where it does not say otherwise, which is
		// decided once the centre is settled — and it may only be settled by an
		// href, so it waits for [gradient.resolved].
		num("fx", &g.fx, 0)
		num("fy", &g.fy, 1)
		if g.r < 0 {
			warn("the radius %g is negative, so the gradient has none and is painted as one colour", g.r)
			g.r = 0
		}
		return
	}
	num("x1", &g.x1, 0)
	num("y1", &g.y1, 1)
	num("x2", &g.x2, 0)
	num("y2", &g.y2, 1)
}

// gradientLength reads one number of a gradient's geometry, which is a number
// with an optional unit or a percentage. A percentage is a fraction of whatever
// the units say the numbers are measured against, which is the shape's own box —
// a viewport of one — or, in user units, the side of the drawing the number is
// on. The axis is which of those it is, and -1 is the one axis a percentage has
// no meaning for.
func gradientLength(raw string, viewport float64) (float64, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, false
	}
	if p, ok := strings.CutSuffix(raw, "%"); ok {
		v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
		if err != nil {
			return 0, false
		}
		return v * viewport / 100, true
	}
	return parseLength(raw)
}

// gradientStops reads the `<stop>` children of a gradient. An offset is a number
// or a percentage of the line, held at 0 when the stop does not say — which is
// not the same as an offset of zero, and is why the two are told apart. A stop
// with no colour of its own is black, which is what SVG says rather than what
// looks right.
func gradientStops(kids []*element, warn func(string, ...any)) []gradientStop {
	var stops []gradientStop
	for _, k := range kids {
		if k.Name != "stop" {
			continue
		}
		s := gradientStop{}
		if raw, ok := k.Attr["offset"]; ok {
			v, ok := gradientLength(raw, 100)
			if !ok {
				warn("the stop offset %q is not a number, so it lands at the start of the gradient", raw)
			}
			s.offset, s.hasOffset = clampFloat(v, 0, 1), ok
		}
		if isCurrentColor(k.attr("stop-color")) {
			// The keyword is kept rather than resolved: the colour of the widget
			// is only known when the drawing is painted, and a gradient written in
			// it should follow the widget like any other paint does.
			s.current = true
		} else if c, ok := parsePaint(k.attr("stop-color")); ok {
			s.color = c
		} else {
			s.color = canvas.Black
		}
		if raw, ok := k.Attr["stop-opacity"]; ok {
			v, ok := parseAlpha(raw)
			if !ok {
				warn("the stop opacity %q is not a number, so the stop is fully opaque", raw)
				v = 1
			}
			if !s.current {
				s.color = fade(s.color, v)
			}
		}
		stops = append(stops, s)
	}
	return stops
}

// gradientHref is the id another gradient's attributes come from. An href is
// written as a plain `#id`, and files in the wild spell it `url(#id)` the way a
// paint is spelled, since the two say the same thing to a reader of the drawing.
func gradientHref(attrs map[string]string) string {
	for _, name := range []string{"href", "xlink:href"} {
		v, ok := attrs[name]
		if !ok {
			continue
		}
		if id, ok := hashRef(v); ok {
			return id
		}
	}
	return ""
}

// hashRef reads the `#id` out of either spelling of a reference: the plain one
// an href is written in, and the `url(#id)` a paint takes.
func hashRef(s string) (string, bool) {
	if id, ok := paintRef(s); ok {
		return id, true
	}
	s = strings.Trim(strings.TrimSpace(s), "\"'")
	if !strings.HasPrefix(s, "#") {
		return "", false
	}
	return s[1:], s[1:] != ""
}

// resolved is the gradient as it will be painted: every attribute it does not
// say taken from the gradient its href points at, which may itself point at
// another. The chain is walked once with the whole drawing in hand, which is why
// an href cannot be followed by [parseGradient] on its own.
//
// A gradient that points at itself, or in a circle, stops where it would start
// repeating: everything such a chain has to say is already in it by then, and
// following it further would loop.
//
// A radial gradient that never said where its focus is has it at the centre,
// which is only known once the chain above has settled the centre itself.
func (g *gradient) resolved(all map[string]*gradient) *gradient {
	out := gradientDefaults(g.radial)
	out.viewport = g.viewport
	var chain []*gradient
	for at, seen := g, map[*gradient]bool{}; !seen[at]; {
		seen[at] = true
		chain = append(chain, at)
		if at.href == "" {
			break
		}
		next, ok := all[at.href]
		if !ok {
			break
		}
		at = next
	}
	// The chain runs nearest first, and the first gradient to have said something
	// is the one that counts: an href fills in what was left out and nothing more.
	for _, at := range chain {
		out.inherit(at)
	}
	if g.radial {
		if !out.set["fx"] {
			out.fx = out.cx
		}
		if !out.set["fy"] {
			out.fy = out.cy
		}
	}
	out.resolvePercentages()
	return out
}

// inherit takes from another gradient everything it wrote itself, which is what
// an href means: the stops, and any geometry left out. An attribute written here
// is kept, so a gradient can take the stops of another and move its own line.
func (out *gradient) inherit(from *gradient) {
	if len(out.stops) == 0 {
		out.stops = from.stops
	}
	for name := range from.set {
		if out.set[name] {
			continue
		}
		switch name {
		case "gradientunits":
			out.userSpace = from.userSpace
		case "spreadmethod":
			out.spread = from.spread
		case "gradienttransform":
			out.transform, out.hasXform = from.transform, from.hasXform
		case "x1":
			out.x1 = from.x1
		case "y1":
			out.y1 = from.y1
		case "x2":
			out.x2 = from.x2
		case "y2":
			out.y2 = from.y2
		case "cx":
			out.cx = from.cx
		case "cy":
			out.cy = from.cy
		case "r":
			out.r = from.r
		case "fx":
			out.fx = from.fx
		case "fy":
			out.fy = from.fy
		}
		if p, ok := from.pct[name]; ok {
			out.pct[name] = p
		}
		out.set[name] = true
	}
}

// resolvePercentages is where the percentages kept as written become numbers,
// now that the units they are a percentage of are finally known. A radius lands
// on the normalized diagonal of the drawing in user units, which is what SVG
// measures it against where it has no side of its own to take.
func (out *gradient) resolvePercentages() {
	for name, p := range out.pct {
		v, ok := gradientLength(p.raw, out.percentOf(p.axis))
		if !ok {
			continue
		}
		switch name {
		case "x1":
			out.x1 = v
		case "y1":
			out.y1 = v
		case "x2":
			out.x2 = v
		case "y2":
			out.y2 = v
		case "cx":
			out.cx = v
		case "cy":
			out.cy = v
		case "r":
			out.r = clampFloat(v, 0, math.Inf(1))
		case "fx":
			out.fx = v
		case "fy":
			out.fy = v
		}
	}
}

// percentOf is what one axis of a percentage is a fraction of. In fractions of
// the shape every axis is the whole of one; in user units they are the drawing's
// own width and height, and an axis with no side to it — a radius — takes the
// diagonal.
func (out *gradient) percentOf(axis int) float64 {
	if !out.userSpace {
		return 1
	}
	switch axis {
	case 0:
		return out.viewport[0]
	case 1:
		return out.viewport[1]
	}
	return math.Hypot(out.viewport[0], out.viewport[1]) / math.Sqrt2
}
