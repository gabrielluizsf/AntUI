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

	// pictures are the bitmaps the `<image>`s fetched, by the address each one
	// named, so that a drawing pointing at the same file twice reads it once.
	// What is in here is never written to: an `<image>` that wants the picture
	// somewhere else or in another opacity is given where and how it fits rather
	// than changing what every other element pointing at the same address holds.
	pictures map[string]*canvas.Canvas

	// clips are the clip paths the drawing declared, by the id a `clip-path`
	// names. They are read before anything is built for the same reason as the
	// gradients: a clip is written at the end of the drawing beside them, and
	// what points at it is written first.
	clips map[string]*clipPath

	// masks are the masks the drawing declared, by the id a `mask` names, read
	// for the same reason as the clips and the gradients: a mask is written at
	// the end of the drawing beside them, and what points at it is written
	// first — sometimes inside another mask.
	masks map[string]*maskDef

	// patterns are the patterns the drawing declared, by the id a `fill` or a
	// `stroke` names, read after the masks and for the same reason once more:
	// what is inside a pattern paints with everything the drawing has, and a
	// reference to one written further down the file would find nothing there.
	patterns map[string]*pattern

	// markers are the markers the drawing declared, by the id a `marker-start`
	// names, read after the patterns and for the same reason yet again: what is
	// inside a marker is a picture painted with everything else — a gradient, a
	// clip, a mask, a `<use>`, another marker — so all of it has to be in place
	// before one of them is built.
	markers map[string]*markerDef

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
	// Pic is the picture an `<image>` draws: the bitmap its href pointed at,
	// fitted into the box the drawing gave it. Every other node leaves it nil,
	// the same way as Runs — it is what puts a photograph in a drawing that is
	// otherwise all vectors, and it is read whole while the drawing is read
	// rather than when it is painted. See [Image.imageNode] and [paintPicture].
	Pic *Picture
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
	// mask is what the element's `mask` named: the picture that says how much
	// of what this node paints is there, with the region it reaches already
	// worked out against the box of everything inside it. Like a clip it is not
	// inherited — a group is masked once as a whole, when its picture is made,
	// rather than masked again at every child as it is drawn — so it lives here
	// and not in the Style, and nil means no mask at all.
	//
	// A mask whose content the drawing cannot follow keeps nothing rather than
	// keeping everything, the same as a clip that names a clipPath that is not
	// there. See [Image.maskNamed] and [maskUnder].
	mask *masked
	// filters is what the element's `filter` put its picture through, in the
	// order they were written. Like a clip and a mask it is not inherited — a
	// group is filtered once as a whole, when its picture is made, rather than
	// filtered again at every child as it is drawn — so it lives here and not
	// in the Style, and nothing at all means no filter. What is beyond the
	// picture in a blur comes from [applyFilters].
	filters []filterOp
	// dashed is the same shape as Path, cut into the runs a dashed stroke
	// paints, and nil where the stroke takes no pattern or one that cuts
	// nothing. The cut is made where the shape was written — the one place
	// where a length of the pattern and a length of the path are the same
	// numbers — and the transform of the element and the one that puts the
	// drawing on the canvas carry it along with the shape from there, which is
	// what makes a dashed stroke grow and shrink with the drawing rather than
	// stay the same size on the screen. Path stays whole because the fill,
	// the box and everything that cuts or masks the element need the shape as
	// a shape; it is only the outline that is dashed. Paint takes this one
	// for the stroke. See [Image.build].
	dashed *canvas.Path
	// markers is where the markers this element's marker-start, marker-mid and
	// marker-end named are drawn: each entry is the marker and the matrix that
	// puts it on one vertex, in the order the vertices come — first, the
	// in-between ones, last. It is worked out where the shape is built, while
	// the path is still in the coordinates it was written in, because the
	// vertex and the direction of the path at it are numbers of the shape as
	// it was written — and the transform goes into the matrix the same way it
	// goes into the path, so that a marker follows the vertex wherever the
	// shape goes. What the children inherit is the markers themselves rather
	// than these places; each shape works out its own. See [placeMarkers].
	markers []markerPlace
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
	// NonScalingStroke is `vector-effect="non-scaling-stroke"`: the stroke is
	// drawn at the width the drawing wrote, in the pixels the drawing comes out
	// at, rather than taken through the viewBox and the transforms that carry
	// it there — a hairline that stays a hairline however much the drawing is
	// zoomed. It is not inherited, as the spec has it, so [Style.with] takes it
	// off again before reading what the element itself said, the same as it
	// does for pathLength.
	NonScalingStroke bool
	Opacity          float64
	FillOpacity      float64
	StrokeOpacity    float64
	// Dash is the stroke broken into pieces along the path, with every length
	// of it — the dashes, the gaps and the offset into them — written in the
	// drawing's own units. On a shape whose element wrote a pathLength it
	// holds those pieces stretched to the length that was declared, which is
	// what the shape is cut with; what the shapes around it inherit is the
	// pattern as it was written, since the length it was stretched against
	// belongs to this shape alone. The cut itself happens where the shape is
	// built and not where it is painted — see [Image.build] — so the pattern
	// reaches the pixels the way everything else the drawing wrote does.
	Dash         *canvas.Dash
	Hidden       bool
	Transform    canvas.Matrix
	HasTransform bool
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
	// A fill or a stroke may be a gradient or a pattern rather than a colour,
	// and the three never both hold: `fillGradient` is set only when the paint
	// named a gradient that was in the drawing, `fillPattern` the same for a
	// pattern, and `fillRef` is the id it named until the drawing has been read
	// far enough to say what that is. The fallback colour is the one written
	// after the `url(...)`, which is what the shape is painted with when the
	// reference is not there at all.
	fillGradient          *gradient
	fillPattern           *pattern
	fillRef               string
	fillFallback          canvas.Color
	fillFallbackCurrent   bool
	hasFillFallback       bool
	strokeGradient        *gradient
	strokePattern         *pattern
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
	// maskRef is the id a `mask="url(#id)"` named, kept the same way a clip
	// reference is, and cleared as early: what it turns into is not inherited
	// either, so it is taken off the style before the children are built from
	// it and followed once the whole of this element is there — the region a
	// mask reaches is measured against everything it paints. See
	// [Image.maskNamed] and [maskUnder].
	maskRef string
	// markerStart, markerMid and markerEnd are the `<marker>`s this element
	// draws at the first, the in-between and the last vertex of its shape, and
	// they are inherited, which is what sets them apart from the clip, the
	// mask and the filter above: a marker is not a cut or a picture laid over
	// one shape, it is a mark on a vertex, and a `<g marker-end>` is what every
	// shape under it puts on its own last vertex. So the references turn into
	// the markers themselves rather than into places on the node, and the
	// children take them along with the fill. See [Image.resolveMarkers].
	markerStart *markerDef
	markerMid   *markerDef
	markerEnd   *markerDef
	// markerStartRef, markerMidRef and markerEndRef are those references — the
	// ids the three `url(#id)`s named — until the whole drawing has been read
	// far enough to say what they are, and each is cleared as it is spent so
	// that a group says what it has to say once and not once for every shape
	// inside it.
	markerStartRef string
	markerMidRef   string
	markerEndRef   string
	// filters is the list a `filter="blur(1) grayscale(1)"` gave the element,
	// read where it is written rather than once the whole drawing is in hand:
	// none of the functions needs to know what any id in the file is, and what
	// is not read is said at once. It is not inherited either — a `<g filter>`
	// turns the one picture it makes of its children once — so it comes off
	// the style before the children are built from it and lives on the node;
	// see [readFilters] and [applyFilters].
	filters []filterOp
	// pathLength is the total length the drawing says the path on this element
	// has, in the drawing's own units, which is what the stroke of this shape
	// is measured out against instead of the geometry — see [applyPathLength].
	// It is not inherited, for the same reason the clip, the mask and the
	// filter are not: a length declared by a group is a length of the group's
	// own shape, and a group has none. [Style.with] takes it off again for
	// every element that does not write one of its own.
	pathLength float64
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
	img := &Image{Root: &Node{Name: "svg"}, grads: map[string]*gradient{}, pictures: map[string]*canvas.Canvas{}}
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
	// The masks come after them, for the same reason taken once more: a mask's
	// picture may be cut by a clip, filled with a gradient, drawn through a
	// `<use>` and named from inside another mask, so every id the drawing has
	// to offer has to be known before one of them is built.
	img.readMasks(root)
	// The patterns come last of the things that are pointed at: what is inside
	// one is a picture painted with everything else — a gradient, a clip, a
	// mask, a `<use>`, another pattern — so all of it has to be in place first.
	img.readPatterns(root)
	// The markers come after them, and are last for the same reason taken once
	// more: what is inside a marker is a picture painted with everything else,
	// and a marker may name another written further down the file.
	img.readMarkers(root)
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
	// The marker is followed here too, and it is the one reference that stays
	// on the style after it has been spent: a marker is inherited, so what the
	// children draw on their own vertices is this marker, looked up once here
	// rather than once each down there. A `<use>` takes it along to whatever it
	// points at for the same reason it takes the fill.
	img.resolveMarkers(&st, warn)
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
	// The mask is named where the clip is followed, and its reference comes off
	// the style for the same two reasons: a mask is not inherited either, and
	// the children would each follow it and say the same thing once each if it
	// were still there. What it turns into waits until the whole of this
	// element is built — the region a mask reaches is measured against
	// everything it paints, which is only known once everything inside it is.
	maskRef := st.maskRef
	st.maskRef = ""
	def := img.maskNamed(maskRef, warn)
	// The filter list comes off the style here for the same two reasons, and
	// it goes onto the node as it comes: what the filter is put through is the
	// whole picture this element makes, so the children must not each find the
	// same list to run again, and the list itself needs no part of the drawing
	// to be read — it was already read where it was written.
	filters := st.filters
	st.filters = nil
	n := &Node{Name: e.Name, Style: st, clip: clip, filters: filters}
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
		n.filters = filters
		n.mask = maskUnder(def, st, n)
		return n
	}
	if e.Name == "image" {
		// A photograph is not an outline with a paint over it either: it is a
		// bitmap, fetched from wherever the href points while the drawing is
		// read and fitted into the box the drawing gave it when it is painted.
		// It is built whole here for the same reason as the writing, and there
		// is nothing inside an `<image>` to build nodes for. See
		// [Image.imageNode].
		n := img.imageNode(e, st, warn)
		n.clip = clip
		n.filters = filters
		n.mask = maskUnder(def, st, n)
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
		// The length the drawing declared is spent here, on this node's own
		// style and before the transform moves the path: it is a length of the
		// shape as it was written, and the children are built from the style
		// that still holds the pattern as it was written too.
		applyPathLength(&n.Style, n.Path)
		// The pattern is spent in the same place and for the same reason: the
		// lengths of a dash pattern are lengths of the drawing, so the outline
		// can only be cut into them where it is still in the coordinates it
		// was written in. What comes out is kept beside the shape rather than
		// in it, because the fill and everything measured against the shape —
		// its box, its clip, the region of a mask — need the shape whole, and
		// only the stroke is dashed. From here the two travel together: the
		// transform below and the one that puts the drawing on the canvas move
		// the cut along with the shape, which is what makes the pattern scale
		// with the drawing instead of staying the same on the screen.
		//
		// A pattern that cuts nothing hands the path back as it is, and that
		// is the shape itself rather than a second copy of it: taking it for
		// the cut would put the transform through the path twice. Only what
		// is going to be stroked is cut at all, by the same two conditions
		// the painting checks before it widens anything.
		if st.HasStroke && st.Width > 0 && n.Style.Dash != nil {
			if cut := n.Path.Dashed(*n.Style.Dash); cut != n.Path {
				n.dashed = cut
			}
		}
		// Where the markers this element asked for sit on the shape, worked
		// out while the path is still in the coordinates it was written in:
		// the vertices they land on, the directions the path is going there and
		// the rooms the markers asked for are all numbers of the shape as it
		// was written, and the transform below takes them along with the
		// vertices the same way it takes the shape.
		n.markers = placeMarkers(n.Path, n.Style)
		if st.HasTransform && st.Transform != (canvas.Matrix{}) {
			// The transform belongs to the shape, not to the pixels around it,
			// so it is put into the path: the shape goes where the transform
			// puts it and the stroke is widened afterwards, in the coordinates
			// it landed in. The width it is widened by carries the same
			// transform — see [Style.strokeWidth] — which is what makes a shape
			// scaled twice a stroke twice as wide, the default SVG has and what
			// `vector-effect="non-scaling-stroke"` asks to leave out. A group
			// has no shape of its own to put it on, and carries the transform to
			// its children instead.
			n.Path.Transform(st.Transform)
			if n.dashed != nil {
				n.dashed.Transform(st.Transform)
			}
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
	n.mask = maskUnder(def, st, n)
	return n
}

// resolvePaints follows the `url(#id)` a fill or a stroke named, which only the
// whole drawing can say anything about. What the reference turns into is the
// gradient or the pattern that was in the drawing; the colour written after the
// `url(...)` where there was not; and nothing, with one warning, where the
// shape asked for something that is not there and gave no second choice.
//
// The reference is cleared as it is spent, so that a shape inside a group with a
// missing gradient does not say the same thing once for every shape inside it.
func (img *Image) resolvePaints(st *Style, warn func(string, ...any)) {
	if st.fillRef != "" {
		g, p := img.grads[st.fillRef], img.patterns[st.fillRef]
		st.fillRef, st.fillGradient, st.fillPattern = "", g, p
		switch {
		case g != nil || p != nil:
			st.HasFill = true
		case st.hasFillFallback:
			st.Fill, st.HasFill = st.fillFallback, true
			st.FillCurrent = st.fillFallbackCurrent
		default:
			st.HasFill = false
			warn("the fill names a gradient or a pattern that is not in the drawing and has no colour after it, so the shape is not filled")
		}
	}
	if st.strokeRef != "" {
		g, p := img.grads[st.strokeRef], img.patterns[st.strokeRef]
		st.strokeRef, st.strokeGradient, st.strokePattern = "", g, p
		switch {
		case g != nil || p != nil:
			st.HasStroke = true
		case st.hasStrokeFallback:
			st.Stroke, st.HasStroke = st.strokeFallback, true
			st.StrokeCurrent = st.strokeFallbackCurrent
		default:
			st.HasStroke = false
			warn("the stroke names a gradient or a pattern that is not in the drawing and has no colour after it, so the shape is not outlined")
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
// it fetches itself — see [Image.imageNode] — and it is built whole when its
// turn comes, the way writing is rather than the way a shape is. A `<symbol>`
// is the other way round — it draws nothing where it stands and everything it
// holds is drawn where a `<use>` points at it. A `<text>` is neither: it paints
// where it stands, and is built as writing rather than as a shape.
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
