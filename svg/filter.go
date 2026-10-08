package svg

import (
	"math"
	"strconv"
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// A `filter` says what an element's picture goes through on its way into the
// drawing: a list of functions — `blur(1) grayscale(1) hue-rotate(90deg)` —
// each of which changes the picture of everything to its left, all of them run
// in the order they were written. The picture is the whole of what the element
// draws, its children and its writing included, which is why a filter is not
// inherited: a `<g filter>` turns the one picture it makes of its children
// once, rather than each child being turned as it is drawn, and the same is
// true of a clip and a mask beside it.
//
// The functions here are the ones the canvas already knows how to apply —
// grayscale, sepia, invert, brightness, contrast, hue-rotate and blur, the
// same set a stylesheet's `filter` maps onto — and a `drop-shadow`, which is
// not one of them and is said out loud and left out of the list rather than
// taken to mean nothing: what the list can still do, it still does.
//
// The other way a file asks for a filter is `url(#name)`, a reference to a
// `<filter>` element written elsewhere in the drawing. The id is kept as it
// was written, and what it names is read with everything else the drawing
// points at — see [Image.readFilterDefs] — so the primitives inside one run
// where the reference stood in the list.
//
// The list reaches the canvas in [applyFilters], which runs each step over
// the picture the element made — the whole of the picture, since that is what
// the element is asking to have put through it, and a blur that reaches past
// where the element stopped has room to reach into.

// filterOp is one step of the list the `filter` attribute holds: a function
// with what it was given — a fraction for grayscale, sepia, invert,
// brightness and contrast, degrees for hue-rotate, and a length in the
// drawing's own units for blur — or a reference to a `<filter>` element,
// which stands as the id it named until the whole drawing has been read and
// as the element itself after that. See [readFilters] and
// [Image.followFilterRefs].
type filterOp struct {
	kind   canvas.FilterKind
	amount float64
	ref    string
	def    *filterDef
}

// readFilters reads the list a `filter` attribute holds: `none` is no filter
// at all and says so without complaint, and every function is read on its own,
// so a list this package cannot wholly read still does the part of it that it
// can — one function left out is a warning, not a whole filter dropped.
func readFilters(raw string, warn func(string, ...any)) []filterOp {
	v := strings.TrimSpace(raw)
	if v == "" {
		warn("the filter %q is not one this package can read, so it is left out", raw)
		return nil
	}
	if strings.EqualFold(v, "none") {
		return nil
	}
	var out []filterOp
	for _, tok := range filterTokens(v) {
		if isFilterRef(tok) {
			// A reference to a `<filter>` element elsewhere in the drawing is
			// the other way a file asks for one of these. The id is kept as
			// written — only a caller holding the whole drawing can say what
			// an id in the file is — and the rest of the list goes on. What
			// names nothing the drawing has is said once the drawing is all
			// read; see [Image.followFilterRefs].
			if id, ok := paintRef(tok); ok {
				out = append(out, filterOp{ref: id})
				continue
			}
			warn("the filter reference %q is not a reference to a <filter> this package can follow, so it is left out", tok)
			continue
		}
		if op, ok := readFilter(tok); ok {
			out = append(out, op)
			continue
		}
		warn("the filter %q is not one this package can read, so it is left out", tok)
	}
	return out
}

// filterTokens breaks a filter list into its functions, keeping the arguments
// of each one with it: `drop-shadow(0 1px 2px black)` is one thing to try to
// read and not five, because a space inside the parentheses is part of what it
// says. Only a space at the top level, where nothing is open, comes between
// two functions.
func filterTokens(s string) []string {
	var out []string
	depth, start := 0, 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == '(':
			depth++
		case c == ')':
			if depth > 0 {
				depth--
			}
		case depth == 0 && (c == ' ' || c == '\t' || c == '\n' || c == '\r'):
			if start < i {
				out = append(out, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}

// isFilterRef says the token is a `url(...)` — the reference to a `<filter>`
// element whose id is kept here for the whole drawing to follow.
func isFilterRef(tok string) bool {
	t := strings.TrimSpace(tok)
	return strings.HasPrefix(strings.ToLower(t), "url(") ||
		strings.EqualFold(t, "url")
}

// readFilter reads one function of the list, such as `blur(2px)` or
// `grayscale(50%)`. What is not one of the functions this package can apply,
// does not open and close its arguments, or gives a number that is not one
// answers false, and the caller says so rather than guessing what was meant.
func readFilter(tok string) (filterOp, bool) {
	open := strings.IndexByte(tok, '(')
	if open < 0 || !strings.HasSuffix(tok, ")") {
		return filterOp{}, false
	}
	name := strings.ToLower(strings.TrimSpace(tok[:open]))
	args := strings.TrimSpace(tok[open+1 : len(tok)-1])
	switch name {
	case "grayscale", "sepia", "invert", "brightness", "contrast":
		v, ok := filterAmount(args)
		if !ok {
			return filterOp{}, false
		}
		kind := canvas.FilterGrayscale
		switch name {
		case "sepia":
			kind = canvas.FilterSepia
		case "invert":
			kind = canvas.FilterInvert
		case "brightness":
			kind = canvas.FilterBrightness
		case "contrast":
			kind = canvas.FilterContrast
		}
		return filterOp{kind: kind, amount: v}, true
	case "hue-rotate":
		deg, err := parseAngle(args)
		if err != nil {
			return filterOp{}, false
		}
		return filterOp{kind: canvas.FilterHueRotate, amount: deg}, true
	case "blur":
		if args == "" {
			// A blur of nothing is no blur, which is what the function with no
			// argument says in CSS, and it is read rather than refused so that
			// a list containing it is not warning about something harmless.
			return filterOp{kind: canvas.FilterBlur, amount: 0}, true
		}
		v, ok := parseLength(args)
		if !ok || v < 0 {
			return filterOp{}, false
		}
		return filterOp{kind: canvas.FilterBlur, amount: v}, true
	}
	return filterOp{}, false
}

// filterAmount reads one function's argument: a number or a percentage, where
// an empty argument is the identity of one. Something else is not a number at
// all, and the caller says so instead of the filter quietly becoming full
// strength or none.
func filterAmount(args string) (float64, bool) {
	if args == "" {
		return 1, true
	}
	if p, ok := strings.CutSuffix(args, "%"); ok {
		v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
		if err != nil {
			return 0, false
		}
		return v / 100, true
	}
	v, err := strconv.ParseFloat(args, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

// applyFilters runs an element's filter list over the picture it drew, one
// step after another in the order they were written: a grayscale and then a
// blur is not the same picture as a blur and then a grayscale, and the list
// says which one it wants. A step that stands for a `<filter>` element runs
// the primitives inside it over the picture so far, and the steps beside one
// in the same list keep the order they were written in.
//
// Every step goes over the whole picture the layer holds, which is exactly
// what the element asked to have put through it. A smaller rectangle would
// have to be what the element painted, widened before a blur by as far as three
// box passes can carry a picture past where the element stopped — and the
// canvas only remembers where it was written when it is watching for changes,
// which this one is not, so the painted rectangle is not something to be had
// here. The whole of the layer is honest about that, gives a blur room to
// reach in every direction, and costs a walk over transparent pixels that come
// out of a colour filter untouched and of a blur still carrying nothing.
//
// A blur is a size and is written in the drawing's own units, so it is turned
// into pixels here the same as a stroke's width is: one unit of the drawing is
// scale pixels of the canvas, and the scale has the element's transform in it
// because what a `<g transform="scale(2)">` draws is twice as big and its blur
// has to be too. The rest of the functions are fractions and angles, which a
// scale does not touch, and the lengths a primitive writes for itself are
// turned into pixels where the primitive is applied.
//
// The picture may come back as a different one than it went in as, because
// every primitive leaves a picture of its own behind, so the caller takes
// what comes out rather than the layer it put in.
func applyFilters(cv *canvas.Canvas, filters []filterOp, scale float64, current canvas.Color) *canvas.Canvas {
	for _, f := range filters {
		if f.def != nil {
			cv = applyFilterDef(cv, f.def, scale, current)
			continue
		}
		amount := f.amount
		if f.kind == canvas.FilterBlur {
			amount = f.amount * scale
		}
		cv.FilterRegion(0, 0, cv.Width, cv.Height, f.kind, amount)
	}
	return cv
}

// A `<filter>` is what an element's picture goes through, written as one
// primitive after another inside the element itself: each takes its inputs —
// the picture before it, the picture the element drew, a picture an earlier
// primitive named — and leaves a picture behind under its own name for what
// follows to pick up. What is read here are the six primitives this package
// can apply: feGaussianBlur, feColorMatrix, feOffset, feFlood, feComposite and
// feBlend. One written otherwise is said and left out, the same as a function
// in the attribute that is not one this package can do.
//
// Every primitive works over the whole picture, the way [applyFilters] does:
// there is no filter region here, so a flood floods everything and a
// composite takes its two pictures whole. ponytail: the filter region — the
// x/y/width/height on the element and on the primitives — is the upgrade
// path; the whole layer is what today's CSS functions already answer to.

// filterDef is a `<filter>` element: the primitives it holds, in the order
// they were written, which is the order they run in.
type filterDef struct {
	prims []filterPrim
}

// primKind is one of the filter primitives this package can apply.
type primKind uint8

const (
	primGaussianBlur primKind = iota
	primColorMatrix
	primOffset
	primFlood
	primComposite
	primBlend
)

// compositeOp is feComposite's operator; the zero is "over", which is what the
// primitive means when the operator is left off.
type compositeOp uint8

const (
	compOver compositeOp = iota
	compIn
	compOut
	compAtop
	compXor
	compLighter
	compArithmetic
)

// blendMode is feBlend's mode; the zero is "normal", which is what the
// primitive means when the mode is left off.
type blendMode uint8

const (
	blendNormal blendMode = iota
	blendMultiply
	blendScreen
	blendDarken
	blendLighten
)

// filterPrim is one primitive of a `<filter>`: what it takes, what it leaves
// its result under, and the handful of numbers its own kind needs. The inputs
// are the names as written; what they name is resolved while the picture is
// being made, where the pictures are to hand.
type filterPrim struct {
	kind            primKind
	in, in2, result string
	stdDev          float64      // feGaussianBlur
	dx, dy          float64      // feOffset, in the drawing's own units
	flood           canvas.Color // feFlood, its colour and opacity apart
	floodOpacity    float64      // feFlood, 1 when none was written
	floodCurrent    bool         // feFlood's colour was currentColor
	matrix          [20]float64  // feColorMatrix
	comp            compositeOp  // feComposite
	blend           blendMode    // feBlend
	k               [4]float64   // feComposite's arithmetic
}

// readFilterDefs reads every `<filter>` in the drawing, by the id an element's
// `filter="url(#name)"` names. They are read while the drawing's ids are all
// in place and before anything is built, for the same reason as the masks and
// the gradients: a filter may be written after the element that points at it.
// What is inside a filter is only what its primitives say — no picture painted
// with the rest of the drawing — so the ids being in place is all one needs to
// be read.
func (img *Image) readFilterDefs(root *element) {
	type declared struct {
		id string
		e  *element
	}
	var found []declared
	seen := map[string]bool{}
	var walk func(e *element)
	walk = func(e *element) {
		if e.Name == "filter" {
			if id := e.attr("id"); id != "" && !seen[id] {
				seen[id] = true
				found = append(found, declared{id: id, e: e})
			}
		}
		for _, kid := range e.Kids {
			walk(kid)
		}
	}
	walk(root)
	img.filters = make(map[string]*filterDef, len(found))
	for _, d := range found {
		def := &filterDef{}
		warn := func(format string, args ...any) { img.warn(d.e, format, args...) }
		// What a primitive may name: the two pictures the element drew, the
		// background this package has no way to reach, and every result an
		// earlier primitive in this filter left behind. Naming anything else
		// is said now rather than at paint time.
		known := map[string]bool{
			"SourceGraphic":   true,
			"SourceAlpha":     true,
			"BackgroundImage": true,
			"BackgroundAlpha": true,
		}
		for _, kid := range d.e.Kids {
			p, ok := readFilterPrim(kid, warn)
			if !ok {
				continue
			}
			for _, name := range []string{p.in, p.in2} {
				if name == "" || known[name] {
					if name == "BackgroundImage" || name == "BackgroundAlpha" {
						warn("the background %q a filter reads is not where a filter can see it, so it is taken as nothing", name)
					}
					continue
				}
				warn("the filter input %q names nothing this filter has, so it is taken as nothing", name)
			}
			def.prims = append(def.prims, p)
			if p.result != "" {
				known[p.result] = true
			}
		}
		img.filters[d.id] = def
	}
}

// readFilterPrim reads one primitive of a `<filter>` — the six this package
// can apply, with the inputs it takes and the name it leaves its result under.
// Anything else inside a `<filter>`, and a primitive whose own numbers are not
// the ones it needs, is said and left out: the rest of the filter still runs,
// the same way one unreadable function in the attribute's list does.
func readFilterPrim(e *element, warn func(string, ...any)) (filterPrim, bool) {
	p := filterPrim{in: e.attr("in"), in2: e.attr("in2"), result: e.attr("result")}
	switch e.Name {
	case "feGaussianBlur":
		p.kind = primGaussianBlur
		if raw := e.attr("stdDeviation"); raw != "" {
			ns, ok := numbers(raw)
			if !ok || len(ns) == 0 {
				warn("the feGaussianBlur stdDeviation %q is not a number this package can read, so the primitive is left out", raw)
				return filterPrim{}, false
			}
			p.stdDev = ns[0]
			if p.stdDev < 0 {
				warn("the feGaussianBlur stdDeviation %q is negative, so the primitive is left out", raw)
				return filterPrim{}, false
			}
			if len(ns) > 1 && ns[1] != ns[0] {
				// ponytail: one radius both ways; a per-axis blur is the
				// upgrade path when a drawing asks for one.
				warn("the second number of the feGaussianBlur stdDeviation %q is not read here, so the blur is the same both ways", raw)
			}
		}
	case "feColorMatrix":
		p.kind = primColorMatrix
		typ := strings.ToLower(strings.TrimSpace(e.attr("type")))
		if typ == "" {
			typ = "matrix"
		}
		raw := e.attr("values")
		switch typ {
		case "matrix":
			if raw == "" {
				p.matrix = identityMatrix()
				break
			}
			ns, ok := numbers(raw)
			if !ok || len(ns) != 20 {
				warn("the feColorMatrix values %q are not the twenty numbers the matrix needs, so the primitive is left out", raw)
				return filterPrim{}, false
			}
			copy(p.matrix[:], ns)
		case "saturate":
			v := 1.0
			if raw != "" {
				var ok bool
				if v, ok = filterAmount(raw); !ok {
					warn("the feColorMatrix saturate %q is not a number this package can read, so the primitive is left out", raw)
					return filterPrim{}, false
				}
			}
			p.matrix = saturateMatrix(v)
		case "huerotate":
			v := 0.0
			if raw != "" {
				var err error
				if v, err = parseAngle(raw); err != nil {
					warn("the feColorMatrix hueRotate %q is not an angle this package can read, so the primitive is left out", raw)
					return filterPrim{}, false
				}
			}
			p.matrix = hueRotateMatrix(v)
		case "luminancetoalpha":
			p.matrix = luminanceMatrix()
		default:
			warn("the feColorMatrix type %q is not one this package can read, so the primitive is left out", typ)
			return filterPrim{}, false
		}
	case "feOffset":
		p.kind = primOffset
		for _, at := range []struct {
			attr string
			slot *float64
		}{{"dx", &p.dx}, {"dy", &p.dy}} {
			if raw := e.attr(at.attr); raw != "" {
				v, ok := parseLength(raw)
				if !ok {
					warn("the feOffset %s %q is not a length this package can read, so the primitive is left out", at.attr, raw)
					return filterPrim{}, false
				}
				*at.slot = v
			}
		}
	case "feFlood":
		p.kind = primFlood
		p.flood = canvas.Color(0xFF000000)
		p.floodOpacity = 1
		if raw := e.attr("flood-color"); raw != "" {
			if isCurrentColor(raw) {
				p.floodCurrent = true
			} else if c, ok := parsePaint(raw); ok {
				p.flood = c
			} else {
				warn("the feFlood flood-color %q is not a colour this package can read, so black is flooded instead", raw)
			}
		}
		if raw := e.attr("flood-opacity"); raw != "" {
			v, ok := filterAmount(raw)
			if !ok {
				warn("the feFlood flood-opacity %q is not a number this package can read, so it is left out", raw)
			} else {
				p.floodOpacity = v
			}
		}
	case "feComposite":
		p.kind = primComposite
		op := strings.ToLower(strings.TrimSpace(e.attr("operator")))
		if op == "" {
			op = "over"
		}
		switch op {
		case "over":
			p.comp = compOver
		case "in":
			p.comp = compIn
		case "out":
			p.comp = compOut
		case "atop":
			p.comp = compAtop
		case "xor":
			p.comp = compXor
		case "lighter":
			p.comp = compLighter
		case "arithmetic":
			p.comp = compArithmetic
		default:
			warn("the feComposite operator %q is not one this package can read, so \"over\" is used instead", op)
		}
		for i, at := range []string{"k1", "k2", "k3", "k4"} {
			if raw := e.attr(at); raw != "" {
				v, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
				if err != nil {
					warn("the feComposite %s %q is not a number this package can read, so the primitive is left out", at, raw)
					return filterPrim{}, false
				}
				p.k[i] = v
			}
		}
	case "feBlend":
		p.kind = primBlend
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
	default:
		warn("the filter primitive %q is not one this package can read, so it is left out", e.Name)
		return filterPrim{}, false
	}
	return p, true
}

// numbers reads a list of numbers the way a filter writes one: separated by
// spaces, commas or both, and all of them the same kind of thing. A piece that
// is not a number answers false for the whole list — the caller wanted a list
// and what it has is not one.
func numbers(s string) ([]float64, bool) {
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == ','
	})
	out := make([]float64, 0, len(fields))
	for _, f := range fields {
		v, err := strconv.ParseFloat(f, 64)
		if err != nil {
			return nil, false
		}
		out = append(out, v)
	}
	return out, true
}

func identityMatrix() [20]float64 {
	return [20]float64{
		1, 0, 0, 0, 0,
		0, 1, 0, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 0, 1, 0,
	}
}

// saturateMatrix, hueRotateMatrix and luminanceMatrix are the three shortcuts
// of feColorMatrix — `saturate`, `hueRotate` and `luminanceToAlpha` — written
// out as the same twenty numbers `matrix` would take, since that is all a
// matrix ever is here.
func saturateMatrix(s float64) [20]float64 {
	return [20]float64{
		0.213 + 0.787*s, 0.715 - 0.715*s, 0.072 - 0.072*s, 0, 0,
		0.213 - 0.213*s, 0.715 + 0.285*s, 0.072 - 0.072*s, 0, 0,
		0.213 - 0.213*s, 0.715 - 0.715*s, 0.072 + 0.928*s, 0, 0,
		0, 0, 0, 1, 0,
	}
}

func hueRotateMatrix(deg float64) [20]float64 {
	a := deg * math.Pi / 180
	c, s := math.Cos(a), math.Sin(a)
	at := func(base, cosc, sinc float64) float64 { return base + cosc*c + sinc*s }
	return [20]float64{
		at(0.213, 0.787, -0.213), at(0.715, -0.715, -0.715), at(0.072, -0.072, 0.928), 0, 0,
		at(0.213, -0.213, 0.143), at(0.715, 0.285, 0.140), at(0.072, -0.072, -0.283), 0, 0,
		at(0.213, -0.213, -0.787), at(0.715, -0.715, 0.715), at(0.072, 0.928, 0.072), 0, 0,
		0, 0, 0, 1, 0,
	}
}

func luminanceMatrix() [20]float64 {
	return [20]float64{
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0.2125, 0.7154, 0.0721, 0, 0,
	}
}

// followFilterRefs follows what each `url(#name)` in an element's filter list
// names, once the whole drawing is read: the reference was kept as written
// where the list was read, because only a caller holding the drawing can say
// what an id is — the same as a fill, a clip or a mask. What names nothing the
// drawing has is said here, once, and left out of the list the same way a
// function the package cannot read is.
func (img *Image) followFilterRefs(filters []filterOp, warn func(string, ...any)) []filterOp {
	if len(filters) == 0 {
		return filters
	}
	out := make([]filterOp, 0, len(filters))
	for _, f := range filters {
		if f.ref == "" {
			out = append(out, f)
			continue
		}
		def := img.filters[f.ref]
		if def == nil {
			warn("the filter names #%s, which the drawing does not have, so it is left out", f.ref)
			continue
		}
		out = append(out, filterOp{def: def})
	}
	return out
}

// applyFilterDef runs a `<filter>` element's primitives over the picture it
// was pointed at, one after another: each takes the pictures its inputs name
// — the picture so far when it names none, the picture the element drew for
// SourceGraphic, and for anything an earlier primitive left behind — and what
// it leaves under its own name is what the rest may take. What comes out is a
// picture of the whole layer, the same as one CSS function's own picture.
func applyFilterDef(dst *canvas.Canvas, def *filterDef, scale float64, current canvas.Color) *canvas.Canvas {
	if dst == nil || len(def.prims) == 0 {
		return dst
	}
	src := dst
	named := map[string]*canvas.Canvas{}
	var withAlpha *canvas.Canvas
	prev := dst
	for _, p := range def.prims {
		in := filterInput(p.in, prev, src, &withAlpha, named)
		in2 := filterInput(p.in2, src, src, &withAlpha, named)
		out := primPicture(p, in, in2, src, scale, current)
		if out == nil {
			continue
		}
		prev = out
		if p.result != "" {
			named[p.result] = out
		}
	}
	return prev
}

// filterInput is the picture a primitive's `in` or `in2` names: nothing at
// all means what the primitive stands beside — the picture so far for `in`,
// the picture the element drew for `in2` — and SourceGraphic, SourceAlpha and
// the results an earlier primitive left are the rest. A name nothing holds
// was said while the filter was being read; the blank picture is what it was
// said to mean.
func filterInput(name string, deflt, src *canvas.Canvas, withAlpha **canvas.Canvas, named map[string]*canvas.Canvas) *canvas.Canvas {
	switch {
	case name == "":
		return deflt
	case name == "SourceGraphic":
		return src
	case name == "SourceAlpha":
		if *withAlpha == nil {
			*withAlpha = alphaPicture(src)
		}
		return *withAlpha
	case name == "BackgroundImage" || name == "BackgroundAlpha":
		return blankPicture(src.Width, src.Height)
	default:
		if c := named[name]; c != nil {
			return c
		}
		return blankPicture(src.Width, src.Height)
	}
}

// primPicture is what one primitive leaves behind: its own picture of the
// pictures it took, never one of them itself, because the same picture may be
// named again by what follows and a primitive must not change a picture it
// does not own. The lengths a primitive writes for itself — an offset, a
// blur's radius — are turned into pixels here, the way a blur's in the
// attribute is turned in [applyFilters]: one unit of the drawing is scale
// pixels of the canvas.
func primPicture(p filterPrim, in, in2, src *canvas.Canvas, scale float64, current canvas.Color) *canvas.Canvas {
	switch p.kind {
	case primGaussianBlur:
		out := pictureOf(in)
		if out == nil {
			return nil
		}
		out.FilterRegion(0, 0, out.Width, out.Height, canvas.FilterBlur, p.stdDev*scale)
		return out
	case primColorMatrix:
		out := pictureOf(in)
		colorMatrix(out, p.matrix)
		return out
	case primOffset:
		out := blankPicture(src.Width, src.Height)
		if out == nil || in == nil {
			return out
		}
		out.BlitOver(int(math.Round(p.dx*scale)), int(math.Round(p.dy*scale)), in)
		return out
	case primFlood:
		out := blankPicture(src.Width, src.Height)
		if out == nil {
			return nil
		}
		c := p.flood
		if p.floodCurrent {
			c = current
		}
		out.FillRect(0, 0, out.Width, out.Height, fade(c, p.floodOpacity))
		return out
	case primComposite:
		out := pictureOf(in)
		pairMap(out, in2, func(a, b canvas.Color) canvas.Color { return compositePx(a, b, p.comp, p.k) })
		return out
	case primBlend:
		out := pictureOf(in)
		pairMap(out, in2, func(a, b canvas.Color) canvas.Color { return blendPx(a, b, p.blend) })
		return out
	}
	return nil
}

// pictureOf is a picture of the one it is given: the primitives draw into a
// copy of their input so that what they took may be named again by what
// follows and still be as it was.
func pictureOf(src *canvas.Canvas) *canvas.Canvas {
	if src == nil {
		return nil
	}
	out := blankPicture(src.Width, src.Height)
	if out == nil {
		return nil
	}
	if src.Pixels != nil && out.Pixels != nil {
		for y := 0; y < src.Height; y++ {
			copy(out.Pixels[y*out.Stride:y*out.Stride+src.Width], src.Pixels[y*src.Stride:y*src.Stride+src.Width])
		}
		return out
	}
	out.BlitOver(0, 0, src)
	return out
}

// blankPicture is a cleared picture the whole layer big, which is what a
// primitive that paints rather than changes starts from.
func blankPicture(w, h int) *canvas.Canvas {
	out, err := canvas.NewLayer(w, h)
	if err != nil {
		return nil
	}
	return out
}

// alphaPicture is a picture of only where there is something — the colour
// taken off, keeping the shape a filter can put its own colour into, which is
// what SourceAlpha names.
func alphaPicture(src *canvas.Canvas) *canvas.Canvas {
	out := pictureOf(src)
	if out == nil || out.Pixels == nil {
		return out
	}
	for y := 0; y < out.Height; y++ {
		row := out.Pixels[y*out.Stride : y*out.Stride+out.Width]
		for x := range row {
			row[x] &= 0xFF000000
		}
	}
	return out
}

// pairMap runs a two-input primitive's colour over every pixel the two
// pictures have, the one it is given first and the one it takes second, and
// writes the answer into the first. A pixel of one picture the other does not
// reach keeps what it had.
func pairMap(dst, src *canvas.Canvas, f func(a, b canvas.Color) canvas.Color) {
	if dst == nil || src == nil || dst.Pixels == nil || src.Pixels == nil {
		return
	}
	for y := 0; y < dst.Height && y < src.Height; y++ {
		drow := dst.Pixels[y*dst.Stride : y*dst.Stride+dst.Width]
		srow := src.Pixels[y*src.Stride : y*src.Stride+src.Width]
		for x := 0; x < dst.Width && x < src.Width; x++ {
			drow[x] = f(drow[x], srow[x])
		}
	}
}

// colorMatrix runs feColorMatrix over a picture. The matrix works on the
// colour as it stands — not on it pressed against the background — so the
// four channels of every pixel go in straight and the four that come out are
// clamped and written back.
func colorMatrix(cv *canvas.Canvas, m [20]float64) {
	if cv == nil || cv.Pixels == nil {
		return
	}
	for y := 0; y < cv.Height; y++ {
		row := cv.Pixels[y*cv.Stride : y*cv.Stride+cv.Width]
		for x, c := range row {
			r, g, b, a := frac(c.R()), frac(c.G()), frac(c.B()), frac(c.A())
			rp := clamp01(m[0]*r + m[1]*g + m[2]*b + m[3]*a + m[4])
			gp := clamp01(m[5]*r + m[6]*g + m[7]*b + m[8]*a + m[9])
			bp := clamp01(m[10]*r + m[11]*g + m[12]*b + m[13]*a + m[14])
			ap := clamp01(m[15]*r + m[16]*g + m[17]*b + m[18]*a + m[19])
			row[x] = rgbaColor(rp, gp, bp, ap)
		}
	}
}

// compositePx is one pixel of feComposite: in is the source — `in` — and out
// the destination — `in2` — with the operations the filter effects spec
// names, and the arithmetic the same on every channel including the alpha.
// The two are held against their own alpha while the operation runs and come
// out of it against theirs, which is what the spec's formulas write.
func compositePx(in, out canvas.Color, op compositeOp, k [4]float64) canvas.Color {
	if op == compOver {
		return canvas.BlendOver(out, in)
	}
	qa, qb := frac(in.A()), frac(out.A())
	ca := [3]float64{frac(in.R()) * qa, frac(in.G()) * qa, frac(in.B()) * qa}
	cb := [3]float64{frac(out.R()) * qb, frac(out.G()) * qb, frac(out.B()) * qb}
	var cr [3]float64
	var qr float64
	switch op {
	case compIn:
		qr = qa * qb
		for i := range cr {
			cr[i] = ca[i] * qb
		}
	case compOut:
		qr = qa * (1 - qb)
		for i := range cr {
			cr[i] = ca[i] * (1 - qb)
		}
	case compAtop:
		qr = qb
		for i := range cr {
			cr[i] = ca[i]*qb + cb[i]*(1-qa)
		}
	case compXor:
		qr = qa*(1-qb) + qb*(1-qa)
		for i := range cr {
			cr[i] = ca[i]*(1-qb) + cb[i]*(1-qa)
		}
	case compLighter:
		qr = math.Min(1, qa+qb)
		for i := range cr {
			cr[i] = math.Min(1, ca[i]+cb[i])
		}
	case compArithmetic:
		qr = clamp01(k[0]*qa*qb + k[1]*qa + k[2]*qb + k[3])
		for i := range cr {
			cr[i] = clamp01(k[0]*ca[i]*cb[i] + k[1]*ca[i] + k[2]*cb[i] + k[3])
		}
	}
	return unPremultiply(cr, qr)
}

// blendPx is one pixel of feBlend, with the formulas the filter effects spec
// writes: A is `in` — the top of the two — and B is `in2`, the bottom. The
// colour is held against each picture's own alpha while the mode runs and
// comes out against the alpha of what the mode gave back.
func blendPx(a, b canvas.Color, mode blendMode) canvas.Color {
	if mode == blendNormal {
		return canvas.BlendOver(b, a)
	}
	qa, qb := frac(a.A()), frac(b.A())
	ca := [3]float64{frac(a.R()) * qa, frac(a.G()) * qa, frac(a.B()) * qa}
	cb := [3]float64{frac(b.R()) * qb, frac(b.G()) * qb, frac(b.B()) * qb}
	var cr [3]float64
	switch mode {
	case blendMultiply:
		for i := range cr {
			cr[i] = (1-qa)*cb[i] + (1-qb)*ca[i] + ca[i]*cb[i]
		}
	case blendScreen:
		for i := range cr {
			cr[i] = cb[i] + ca[i] - ca[i]*cb[i]
		}
	case blendDarken:
		for i := range cr {
			cr[i] = math.Min((1-qa)*cb[i]+ca[i], (1-qb)*ca[i]+cb[i])
		}
	case blendLighten:
		for i := range cr {
			cr[i] = math.Max((1-qa)*cb[i]+ca[i], (1-qb)*ca[i]+cb[i])
		}
	}
	return unPremultiply(cr, 1-(1-qa)*(1-qb))
}

// unPremultiply takes a colour out of the form the operations ran it in —
// held against the alpha that came out with it — and back to the straight
// colour a pixel is written as.
func unPremultiply(cr [3]float64, qr float64) canvas.Color {
	if qr <= 0 {
		return 0
	}
	return rgbaColor(math.Min(1, cr[0]/qr), math.Min(1, cr[1]/qr), math.Min(1, cr[2]/qr), qr)
}

// frac is a colour's channel as the number an operation works on: between
// zero and one, where a byte's two hundred and fifty-five steps are.
func frac(c uint8) float64 { return float64(c) / 255 }

// rgbaColor is four numbers between zero and one as a pixel's colour.
func rgbaColor(r, g, b, a float64) canvas.Color {
	return canvas.Color(uint32(a*255+0.5)<<24 |
		uint32(r*255+0.5)<<16 |
		uint32(g*255+0.5)<<8 |
		uint32(b*255+0.5))
}

// clamp01 keeps a number a channel can be: between zero and one, where the
// operations may push it past what a colour holds.
func clamp01(v float64) float64 { return math.Max(0, math.Min(1, v)) }
