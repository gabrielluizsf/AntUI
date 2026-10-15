package svg

// markerDef is one `<marker>` of a drawing: the picture it holds and everything
// that says where that picture sits on a vertex.
type markerDef struct {
	// id is the name the marker was written under, which is only here so a
	// warning about it can say which one it was.
	id string
	// content is the marker's own tree, built in the coordinates it was written
	// in — the transform on the marker and on anything inside it counts, since
	// that is part of how it was written — and drawn through the matrix that
	// puts it on a vertex. It is nil only where the drawing has nothing in the
	// marker at all.
	content *Node
	// view is the viewBox the picture was drawn with, and hasView says there
	// was one: without it the picture is written in the coordinates of its own
	// room already, and with it the two are fitted together the way a viewBox
	// fits the canvas, which is the same fit a `<symbol>` gets.
	view    [4]float64
	hasView bool
	// fit and align are the marker's own preserveAspectRatio, which decides
	// how that viewBox is fitted onto the room. Without a viewBox there is
	// nothing to fit and the two hold the centred meeting that a drawing
	// saying nothing gets — which is why they are read even where the
	// attribute is absent or passed over.
	fit   PictureFit
	align PictureAlign
	// refX and refY are the point of the picture that lands on the vertex: the
	// marker is moved back by it before it is put down, which is what makes an
	// arrowhead whose tip is written at 0 0 sit that tip on the vertex rather
	// than its top-left corner beside it.
	refX, refY float64
	// w and h are the room the marker asks for, in whichever units markerUnits
	// says: a length of the drawing with `userSpaceOnUse`, and — the default —
	// a multiple of the stroke-width of the shape that names it, which is what
	// makes a marker on a thick line a thick marker.
	w, h float64
	// strokeUnits is that default of markerUnits rather than a length of the
	// drawing's own.
	strokeUnits bool
	// deg is the angle the marker is turned by where orient says a number,
	// while auto says to turn it to the direction of the path instead, which is
	// the whole of what an arrowhead is for. startReverse is
	// `auto-start-reverse`: the same as auto everywhere but at the start of a
	// path, where the marker is turned around — which is what lets one marker
	// be put at both ends of a line and point out of each of them.
	deg          float64
	auto         bool
	startReverse bool
	// visible is `overflow`: whether the picture may run out of the room the
	// marker asked for. What it does not say — which is what the spec's own
	// stylesheet says a marker stands for — is that it may not, so the room
	// cuts what runs past its edge unless `overflow="visible"` is written.
	visible bool
}
