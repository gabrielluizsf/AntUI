package svg

import (
	"strings"
	"sync"

	"github.com/gabrielluizsf/antui/canvas"
)

// Image is a drawing read from SVG, ready to be painted at whatever size it is
// asked for. One is read once and painted as often as it is needed, which is
// what makes an icon cheaper than the file it came from after the first frame.
//
// An Image is safe to paint from more than one goroutine as long as nothing
// writes to it: painting caches the last size it was drawn at, and that cache
// is guarded.
type Image struct {
	// Width and Height are the size the drawing says it is, in the units of its
	// own viewBox, which is what a caller asking for an icon at twenty-four
	// pixels wants to know it was drawn to be.
	Width, Height float64
	// ViewBox is the rectangle of the drawing's own coordinates that maps onto
	// the size above, as minX, minY, width, height. A drawing with no viewBox
	// has one covering the whole of it.
	ViewBox [4]float64
	// Root is the top of the drawing: the `<svg>` and everything inside it.
	Root *Node

	warnings Warnings

	// grads are the gradients the drawing declared, by the id a `url(#name)`
	// names. They are read before anything is built, because a shape may point
	// at a gradient written after it — which is the order a drawing is usually
	// written in, with the shapes first and the `<defs>` at the end.
	grads map[string]*gradient

	// clips are the clip paths the drawing declared, by the id a `clip-path`
	// names. They are read before anything is built for the same reason as the
	// gradients: a clip is written at the end of the drawing beside them, and
	// what points at it is written first.
	clips map[string]*clipPath

	// ids are the elements that wrote an id, under that id, so that a `<use>`
	// can be followed to the element it names — which may be written after the
	// `<use>`, and usually is, at the end of the drawing next to the gradients.
	ids map[string]*element

	// uses are the ids being drawn through a `<use>` right now, outermost
	// first, while the drawing is being read. A reference that names one of
	// them closes a circle rather than drawing a shape, so it is stopped at
	// the first repeat.
	uses []string

	// The last size and colour painted and the canvas it was painted on, so that
	// a frame which asks for the same icon again costs nothing. Guarded by mu,
	// because two goroutines may ask for the same icon at once and the worst
	// that can happen is that one of them draws it.
	mu          sync.Mutex
	drawn       *canvas.Canvas
	drawnWidth  int
	drawnHeight int
	// drawnColour is the `currentColor` the cached canvas was painted in, which
	// is part of what makes it the canvas asked for: an icon painted red is not
	// the icon asked for a moment later in blue.
	drawnColour canvas.Color
}

// Node is one element of a drawing, with the shape it draws and the style it
// draws it in. A group is a node with nothing of its own to draw, which is what
// carries the attributes its children inherit.
type Node struct {
	// Name is the tag: `rect`, `circle`, `path`, `g`, `svg`.
	Name string
	// Path is the shape the node draws, in the drawing's own coordinates. A
	// group and anything this package cannot draw leave it nil.
	Path *canvas.Path
	// Style is everything the node was told about how to paint itself.
	Style Style
	// Runs are the pieces of writing a `<text>` is made of, in the order they
	// were written. Every other node leaves it empty: writing is the one thing
	// in a drawing that carries on across the elements inside it, so it is kept
	// as the whole of it rather than as a node for each `<tspan>`.
	Runs []TextRun
	// Kids are the nodes written inside this one.
	Kids []*Node
	// box is the shape as the drawing wrote it, measured before the node's
	// transform moved it: x, y, width, height. It is the box a gradient in
	// objectBoundingBox units is a fraction of, because SVG takes the box of the
	// geometry and moves the box and the paint laid over it together afterwards
	// rather than stretching the paint over the moved shape.
	box [4]float64
	// clip is what the element's `clip-path` named: the shapes its painting has
	// to stay inside, with the element's own transform already put on them so
	// that a clip follows the element it cuts. A clip is not inherited — a group
	// is cut once as a whole, when its picture is made, rather than cut again at
	// every child as it is drawn — so it lives here and not in the Style, and
	// nil means no clip at all.
	//
	// A clip with no shapes in it is a different thing from no clip: it keeps
	// nothing, which is what a reference that the drawing cannot follow means
	// everywhere else here too. See [Image.resolveClip].
	clip *clipRegion
}

// find is the first node of that name anywhere at or below this one, and nil
// where there is none. It is how a caller reaches one shape of a drawing to
// measure it or to see what style the file gave it.
func (n *Node) find(name string) *Node {
	if n == nil {
		return nil
	}
	if n.Name == name {
		return n
	}
	for _, k := range n.Kids {
		if got := k.find(name); got != nil {
			return got
		}
	}
	return nil
}

// Style is the painting a node was given: how thick its stroke is, how its ends
// and corners are treated, how transparent it is, and whether it is drawn at
// all. It is kept whole on the node so a caller that wants to restyle a drawing
// can see what the file asked for before deciding what to change.
type Style struct {
	Fill      canvas.Color
	HasFill   bool
	FillRule  canvas.FillRule
	Stroke    canvas.Color
	HasStroke bool
	// ClipRule is the rule the shapes of a clipPath are read with when this
	// element is the one they cut — the clip-rule of the clipPath, inherited
	// down to the shapes inside it the way fill-rule is — and it says which of
	// the insides of a shape that crosses itself counts.
	ClipRule canvas.FillRule
	// FillCurrent and StrokeCurrent say that the paint was the keyword
	// `currentColor`, so the colour the node was given is not the one it is
	// drawn in: it takes the colour the painting was asked for. One drawing
	// written once therefore paints in any colour, which is how a single icon
	// serves a button in every state without being written again.
	FillCurrent   bool
	StrokeCurrent bool
	Width         float64
	Cap           canvas.LineCap
	Join          canvas.LineJoin
	MiterLimit    float64
	Opacity       float64
	FillOpacity   float64
	StrokeOpacity float64
	Dash          *canvas.Dash
	Hidden        bool
	Transform     canvas.Matrix
	HasTransform  bool
	// FontSize is how tall the writing is, in the drawing's own units, and
	// Anchor is where along the pen a piece of it hangs. Both only matter to a
	// `<text>`, but both are inherited like the rest, so a size written on a
	// `<g>` is the size everything inside it is written at.
	//
	// A FontSize of zero is no size written at all rather than no writing: the
	// reader of the whole drawing decides what that means, which is what lets
	// a size written once on the root govern every text under it. See
	// [readFontSize].
	FontSize float64
	Anchor   TextAnchor
	// A fill or a stroke may be a gradient rather than a colour, and the two
	// never both hold: `fillGradient` is set only when the paint named a
	// gradient that was in the drawing, and `fillRef` is the id it named until
	// the drawing has been read far enough to say what that is. The fallback
	// colour is the one written after the `url(...)`, which is what the shape is
	// painted with when the reference is not there at all.
	fillGradient          *gradient
	fillRef               string
	fillFallback          canvas.Color
	fillFallbackCurrent   bool
	hasFillFallback       bool
	strokeGradient        *gradient
	strokeRef             string
	strokeFallback        canvas.Color
	strokeFallbackCurrent bool
	hasStrokeFallback     bool
	// clipRef is the id a `clip-path="url(#id)"` named, kept the same way a
	// fill reference is: only the whole drawing can say what it is, and the
	// clip it turns into is not inherited by the children — they are drawn
	// into this element's picture first and the picture is cut once. It is
	// cleared as it is spent, see [Image.resolveClip].
	clipRef string
}

// Parse reads a drawing and answers it, or an error saying why it could not be
// read at all. The error is only ever structural: a drawing that does not
// close what it opens has no shape to paint, while anything inside it that
// cannot be drawn is left to the warnings, which are asked for afterwards.
//
// A drawing that is not an `<svg>` is an error, because there is nothing to
// paint; one whose size cannot be read falls back to what its viewBox says, and
// a drawing with neither is a hundred by a hundred.
func Parse(src string) (*Image, error) {
	root, err := parseDocument(src)
	if err != nil {
		return nil, err
	}
	if root.Name != "svg" {
		return nil, &Error{What: "the root tag is <" + root.Name + ">, not <svg>"}
	}
	img := &Image{Root: &Node{Name: "svg"}, grads: map[string]*gradient{}}
	img.readSize(root)
	// The gradients come before the shapes, because a shape may name one that is
	// only written later in the file, and every gradient that names another has
	// to be followed before either is asked about a colour.
	img.readGradients(root)
	// The ids come with them, for the same reason: a `<use>` may name an element
	// that is written after it, and only the whole drawing can say what an id is.
	img.readIds(root)
	// The clip paths come last of the things that are pointed at, since a
	// `clip-path` may name one written after the element it cuts, and building
	// them is building shapes that may themselves be pointed at.
	img.readClipPaths(root)
	img.read(root)
	return img, nil
}

// readIds takes every element that wrote an id, under the id it wrote, for a
// `<use>` to be followed to what it names. The first element to an id is the
// one it answers to, the same as the gradients: a drawing that writes the same
// id twice is broken rather than undecided, and guessing which one was meant
// would be worse than picking the first.
func (img *Image) readIds(root *element) {
	img.ids = map[string]*element{}
	collectIds(root, img.ids)
}

// collectIds walks the drawing once, taking every element that wrote an id.
// An id nothing points at is still kept: it costs nothing and a drawing may
// point at one later, from a `<use>` this has not read yet.
func collectIds(e *element, into map[string]*element) {
	if e == nil {
		return
	}
	if id := e.attr("id"); id != "" {
		if _, taken := into[id]; !taken {
			into[id] = e
		}
	}
	for _, k := range e.Kids {
		collectIds(k, into)
	}
}

// readGradients walks the drawing for the gradients it declares, keyed by the id
// a `url(#id)` names, and then resolves each one against the rest: an href takes
// what it does not say from another gradient, and only the whole drawing can say
// what that other one had.
func (img *Image) readGradients(root *element) {
	raw := map[string]*gradient{}
	collectGradients(root, img, raw)
	for id, g := range raw {
		img.grads[id] = g.resolved(raw)
	}
}

// collectGradients takes every gradient element in the drawing, once each, under
// the id it answers to. An id nothing can name is still kept: a drawing may
// point at a gradient by an id it reuses, and a gradient declared without one is
// still part of the drawing even though no `url()` can reach it.
func collectGradients(e *element, img *Image, into map[string]*gradient) {
	if e == nil {
		return
	}
	if e.Name == "linearGradient" || e.Name == "radialGradient" {
		warn := func(format string, args ...any) { img.warn(e, format, args...) }
		g := parseGradient(e, warn)
		// A percentage written here is a fraction of the drawing, which is the
		// viewBox rather than the pixels it will be drawn into.
		g.viewport = [2]float64{img.ViewBox[2], img.ViewBox[3]}
		if id := e.attr("id"); id != "" {
			if _, taken := into[id]; !taken {
				into[id] = g
			}
		}
	}
	for _, k := range e.Kids {
		collectGradients(k, img, into)
	}
}

// read turns the element tree into the node tree, carrying the style down as it
// goes so that a node's fill is the one it was given or, failing that, the one
// it inherited from the element above.
func (img *Image) read(e *element) {
	img.Root.Kids = append(img.Root.Kids, img.node(e, baseStyle()))
}

// baseStyle is where the style of everything starts: the fill a drawing that
// says nothing about colour paints with, a stroke the width the spec gives it,
// the opacities that let anything be seen through, and the transform nothing
// has moved from. A `<clipPath>` is built from this too rather than from the
// style where it stands, because its shapes are written in the coordinates of
// whoever points at it and take only what was written on the way down into
// them.
func baseStyle() Style {
	return Style{
		Fill:        canvas.RGB(0, 0, 0),
		HasFill:     true,
		FillRule:    canvas.FillNonZero,
		Width:       1,
		Cap:         canvas.CapButt,
		Join:        canvas.JoinMiter,
		MiterLimit:  canvas.DefaultMiterLimit,
		Opacity:     1,
		FillOpacity: 1, StrokeOpacity: 1,
		Transform: canvas.Identity(), HasTransform: true,
	}
}

// node builds one node and the nodes inside it.
func (img *Image) node(e *element, inherited Style) *Node {
	return img.build(e, inherited, false)
}

// build builds one node and the nodes inside it. forceKids says that what is
// inside this element is drawn after all, which is the one case of a
// definition being taken out of where it stands and put into the picture: a
// `<use>` naming a `<symbol>`. Everywhere else a definition keeps its contents
// to itself, because they are material for pointing at rather than part of the
// picture.
func (img *Image) build(e *element, inherited Style, forceKids bool) *Node {
	warn := func(format string, args ...any) { img.warn(e, format, args...) }
	st := inherited.with(e, warn)
	// The gradients are read before any node is built, so a `url(#id)` named
	// here can be followed now and the children inherit whatever it turned out to
	// be, the same way they inherit a colour.
	img.resolvePaints(&st, warn)
	if e.Name == "use" {
		// A `<use>` follows its own clip itself, after its x and y have moved
		// it and before it builds what it points at — see [Image.useNode].
		return img.useNode(e, st, warn)
	}
	// The clip is followed here too, and spent: what it turns into is this
	// element's alone — the children are drawn into the picture this element
	// makes and the picture is cut once, so a clip is not something they can
	// inherit. Only the reference is in the style, and it is cleared on the way
	// out so that a `<g clip-path>` says what it has to say once and not once
	// for every shape inside it.
	clip := img.resolveClip(&st, warn)
	n := &Node{Name: e.Name, Style: st, clip: clip}
	if n.Style.Hidden {
		return n
	}
	if e.Name == "text" || e.Name == "tspan" {
		// Writing is not an outline with a fill and a stroke over it: it is a
		// piece of writing with a place and a size, so it is kept as the pieces
		// it was written in rather than as a shape, and there is nothing inside
		// it to build nodes for — a `<tspan>` is flattened into the writing
		// around it, because it means nothing on its own.
		n := img.textNode(e, st, warn)
		n.clip = clip
		return n
	}
	n.Path = img.shape(e, st)
	if n.Path != nil {
		if minX, minY, maxX, maxY, ok := n.Path.Bounds(); ok {
			// The box as the drawing wrote it, before the transform moved it: a
			// gradient in fractions of the shape's own box is a fraction of
			// where the shape was, and the transform moves the paint over it
			// along with the shape rather than stretching it again.
			n.box = [4]float64{minX, minY, maxX - minX, maxY - minY}
		}
		if st.HasTransform && st.Transform != (canvas.Matrix{}) {
			// The transform belongs to the shape, not to the pixels around it, so
			// it is applied to the path: a stroke then keeps its own width instead
			// of being scaled with the drawing. A group has no shape of its own to
			// put it on, and carries the transform to its children instead.
			n.Path.Transform(st.Transform)
		}
	}
	if !forceKids && isKnownButUnpainted(e.Name) {
		// What is written inside a `<defs>`, a `<clipPath>` or a `<symbol>` is
		// kept out of the picture the way SVG keeps it: as though it were
		// hidden, present in the drawing for something to point at and painted
		// nowhere. Building the nodes would paint it where it stands, and a
		// `<use>` would then paint it a second time.
		return n
	}
	for _, k := range e.Kids {
		n.Kids = append(n.Kids, img.node(k, st))
	}
	return n
}

// resolvePaints follows the `url(#id)` a fill or a stroke named, which only the
// whole drawing can say anything about. What the reference turns into is the
// gradient itself, where there was one in the drawing; the colour written after
// the `url(...)` where there was not; and nothing, with one warning, where the
// shape asked for something that is not there and gave no second choice.
//
// The reference is cleared as it is spent, so that a shape inside a group with a
// missing gradient does not say the same thing once for every shape inside it.
func (img *Image) resolvePaints(st *Style, warn func(string, ...any)) {
	if st.fillRef != "" {
		g := img.grads[st.fillRef]
		st.fillRef, st.fillGradient = "", g
		switch {
		case g != nil:
			st.HasFill = true
		case st.hasFillFallback:
			st.Fill, st.HasFill = st.fillFallback, true
			st.FillCurrent = st.fillFallbackCurrent
		default:
			st.HasFill = false
			warn("the fill names a gradient that is not in the drawing and has no colour after it, so the shape is not filled")
		}
	}
	if st.strokeRef != "" {
		g := img.grads[st.strokeRef]
		st.strokeRef, st.strokeGradient = "", g
		switch {
		case g != nil:
			st.HasStroke = true
		case st.hasStrokeFallback:
			st.Stroke, st.HasStroke = st.strokeFallback, true
			st.StrokeCurrent = st.strokeFallbackCurrent
		default:
			st.HasStroke = false
			warn("the stroke names a gradient that is not in the drawing and has no colour after it, so the shape is not outlined")
		}
	}
}

// shape is the outline a node draws, or nil when it draws none: a group, a
// `<defs>` of things nothing points at, or a tag this package does not know.
func (img *Image) shape(e *element, st Style) *canvas.Path {
	warn := func(format string, args ...any) { img.warn(e, format, args...) }
	switch e.Name {
	case "path":
		return pathData(e.attr("d"), warn)
	case "rect":
		return rectShape(e, warn)
	case "circle":
		return circleShape(e, warn)
	case "ellipse":
		return ellipseShape(e, warn)
	case "line":
		return lineShape(e)
	case "polyline":
		return pointsShape(e.attr("points"), false, warn)
	case "polygon":
		return pointsShape(e.attr("points"), true, warn)
	}
	switch e.Name {
	case "g", "svg", "a", "switch":
		return nil
	}
	if isKnownButUnpainted(e.Name) {
		return nil
	}
	warn("%q is not a shape this package draws, so it is left out", e.Name)
	return nil
}

// isKnownButUnpainted says whether a tag is one SVG defines that paints nothing
// where it stands. They are left out quietly: a `<defs>` or a `<title>` in the
// middle of an icon is normal, and warning about each of them would bury the
// warnings that matter.
//
// A tag that paints where it stands is not in here: `<image>` draws a picture
// this package has no way to fetch or decode, and a drawing that asked for one
// has asked for a picture that will not be there, which is worth saying. A
// `<symbol>` is the other way round — it draws nothing where it stands and
// everything it holds is drawn where a `<use>` points at it. A `<text>` is
// neither: it paints where it stands, and is built as writing rather than as a
// shape.
func isKnownButUnpainted(name string) bool {
	switch name {
	case "defs", "title", "desc", "metadata", "style", "styleSheet", "script",
		"view", "font", "fontFace", "clipPath", "mask", "pattern", "symbol",
		"linearGradient", "radialGradient", "stop", "filter", "marker", "set",
		"animate", "animateTransform", "animateMotion":
		return true
	}
	return false
}

// readSize takes the size of the drawing from its width and height, falling
// back to the viewBox and then to a hundred square, and makes sure the viewBox
// is one that describes the drawing: an absent or unreadable one is taken from
// the size. There is always one of the three, so this never fails — a drawing
// that says nothing about its size is still a drawing, and it is painted at
// the size a caller asks for either way.
func (img *Image) readSize(e *element) {
	img.Width, img.Height = 100, 100
	hasW, hasH := e.hasAttr("width"), e.hasAttr("height")
	if w, ok := parseLength(e.attr("width")); ok {
		img.Width = w
	} else {
		hasW = false
	}
	if h, ok := parseLength(e.attr("height")); ok {
		img.Height = h
	} else {
		hasH = false
	}
	if v, ok := parseViewBox(e.attr("viewBox")); ok {
		img.ViewBox = v
		if !hasW && !hasH {
			// A drawing with only a viewBox is sized by it, which is how a file
			// meant to scale says how big it is without saying.
			img.Width, img.Height = v[2], v[3]
		}
		return
	}
	// No viewBox: the drawing's coordinates are the size, starting at the
	// origin, so that is what it is.
	img.ViewBox = [4]float64{0, 0, img.Width, img.Height}
}

// parseViewBox reads the four numbers of a viewBox, answering false for one
// that is not four numbers or that has nothing to show.
func parseViewBox(s string) ([4]float64, bool) {
	fields := strings.FieldsFunc(s, func(r rune) bool { return r == ' ' || r == ',' || r == '\t' || r == '\n' || r == '\r' })
	if len(fields) != 4 {
		return [4]float64{}, false
	}
	var v [4]float64
	for i, f := range fields {
		n, err := parseNumber(f)
		if err != nil {
			return [4]float64{}, false
		}
		v[i] = n
	}
	if v[2] <= 0 || v[3] <= 0 {
		return [4]float64{}, false
	}
	return v, true
}

// warn records one thing the drawing asked for that could not be done.
func (img *Image) warn(e *element, format string, args ...any) {
	img.warnings.warn(e.Name, format, args...)
}

// Warnings is everything the drawing asked for that this package could not do,
// read after it is painted. A drawing that came out as it should has none.
func (img *Image) Warnings() Warnings { return img.warnings }
