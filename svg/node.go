package svg

import "github.com/gabrielluizsf/antui/canvas"

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
