package svg

import "github.com/gabrielluizsf/antui/canvas"

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
	// same list to run again. The functions in it were already read where the
	// list was written; only a reference still needs the drawing to say what
	// its id names, which is said here, once — the same way a fill, a clip and
	// a mask follow their own ids.
	filters := st.filters
	st.filters = nil
	filters = img.followFilterRefs(filters, warn)
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
		n.clip = clipUnder(n.clip, st, n, warn)
		n.filters = filters
		n.mask = maskUnder(def, st, n, warn)
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
		n.clip = clipUnder(n.clip, st, n, warn)
		n.filters = filters
		n.mask = maskUnder(def, st, n, warn)
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
	n.mask = maskUnder(def, st, n, warn)
	n.clip = clipUnder(n.clip, st, n, warn)
	return n
}
