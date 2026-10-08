package svg

import (
	"math"
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// A `<marker>` is a little picture drawn at a vertex of the shape that names
// it: an arrowhead at the start of a line, a bead at the end of a trail, a dot
// on every corner of a path. The picture is written once — usually at the end
// of the drawing, beside the gradients and the clip paths — and any number of
// shapes point at it with `marker-start="url(#id)"`, `marker-mid` and
// `marker-end`, which are the first vertex, every vertex between the first and
// the last, and the last one.
//
// Like a gradient, a marker may be written after the element that names it, so
// every `<marker>` in the drawing is read before any of it is built. Unlike a
// gradient, where the marker is written says nothing about where it lands: the
// picture stays in the coordinates it was written in, and what puts it on a
// vertex is the shape that asked for it — the vertex itself, the direction the
// path is going as it passes through that vertex, and the room the marker
// asked for — so the matrix that does all three is worked out where the shape
// is built and kept beside it, the same way the dashed outline is.
//
// A marker is not drawn where its own `<marker>` stands: it is in
// [isKnownButUnpainted], and what is inside it is built only because something
// points at it, which is what forceKids says when it is built.

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

// markerPlace is where one marker sits on one shape: the marker itself and the
// matrix that takes its picture to the drawing's coordinates, with the element's
// own transform already in it under the same condition the shape's path got it,
// so that a marker follows the vertex the way the vertex follows the shape.
type markerPlace struct {
	def *markerDef
	at  canvas.Matrix
	// room is the rectangle of the room the marker asked for, taken back into
	// the coordinates its picture was written in so that it cuts the picture
	// where the room's edges land on it. It is nil where `overflow` says the
	// picture may run out of the room, which draws it whole — and nil where
	// the fit of a viewBox has no way back, which collapses the picture
	// anyway rather than cutting it the wrong way.
	room *clipRegion
	// prop is the property that asked for this one — "marker-start",
	// "marker-mid" or "marker-end" — and is only here so that a warning about
	// a marker going round for ever can say which of the three it was.
	prop string
}

// readMarkers takes every `<marker>` that named itself, under the ids the three
// marker properties may point at, and builds the picture each one holds. They
// are named before any of them is built, the way the clip paths, the masks and
// the patterns are: a marker may be written after the element that names it,
// and what is inside one may name anything else the drawing has — a gradient,
// a clip, a mask, a `<use>`, or another marker written further down the file.
func (img *Image) readMarkers(root *element) {
	type declared struct {
		id string
		e  *element
	}
	var found []declared
	seen := map[string]bool{}
	var walk func(e *element)
	walk = func(e *element) {
		if e.Name == "marker" {
			if id := e.attr("id"); id != "" && !seen[id] {
				seen[id] = true
				found = append(found, declared{id: id, e: e})
			}
		}
		for _, k := range e.Kids {
			walk(k)
		}
	}
	walk(root)

	img.markers = make(map[string]*markerDef, len(found))
	// Every marker is named with nothing in it first, and only then is any of
	// them read and built: a marker's picture may point at another by an id
	// that has not been reached yet, and a reference to a marker that is not
	// there is a warning rather than a picture.
	for _, d := range found {
		img.markers[d.id] = &markerDef{id: d.id}
	}
	for _, d := range found {
		img.readMarker(img.markers[d.id], d.e)
	}
	for _, d := range found {
		img.markers[d.id].content = img.build(d.e, baseStyle(), true)
	}
	img.cutMarkerCycles()
}

// readMarker reads one `<marker>`'s own sayings — the room it asks for, the
// point of its picture that lands on the vertex, and which way it is turned —
// into the place its id already put. What it does not say is the SVG default,
// which is a room three stroke-widths square with the reference at its origin
// and no turn on it.
//
// The style the picture is built from starts at the drawing's own rather than
// at the style of the place the marker stands, which is what [Image.readClipPath]
// says about a clipPath too. The spec has properties inherit into a marker from
// its ancestors and not from the element that names it, and the second half of
// that is the half that matters here: a marker is drawn at a vertex far from
// wherever it was written, in colours of its own, and taking the colours of the
// shape that happens to name it would make the same marker come out differently
// at every end of every line. The first half is what this leaves out — a marker
// under a group that painted everything in it does not take that paint — and a
// marker whose picture wants a colour says so itself, as they are written.
func (img *Image) readMarker(def *markerDef, e *element) {
	warn := func(format string, args ...any) { img.warn(e, format, args...) }
	if raw := strings.TrimSpace(e.attr("viewBox")); raw != "" {
		if v, ok := parseViewBox(raw); ok {
			def.view, def.hasView = v, true
		} else {
			warn("the marker viewBox %q is not four numbers, so its picture is drawn where it was written", raw)
		}
	}
	def.refX = markerLength(e, "refX", 0, warn)
	def.refY = markerLength(e, "refY", 0, warn)
	def.w = markerLength(e, "markerWidth", 3, warn)
	def.h = markerLength(e, "markerHeight", 3, warn)
	def.strokeUnits = !strings.EqualFold(strings.TrimSpace(e.attr("markerUnits")), "userspaceonuse")
	def.deg, def.auto, def.startReverse = readOrient(e.attr("orient"), warn)
	def.visible = overflowIn(e)
	if def.w <= 0 || def.h <= 0 {
		warn("the marker is %g by %g, which has no room for anything in it, so nothing is drawn where it is named", def.w, def.h)
	}
}

// markerLength is one length a marker is written with: what it does not say is
// the default the spec gives it, and what it says but this package cannot read
// is said out loud rather than quietly taken as the default. A percentage is
// left out of that — it is a fraction of the room the marker is drawn into,
// which is not known until the shape that names it is — so it comes through
// here as something that is not a length and is said.
func markerLength(e *element, name string, def float64, warn func(string, ...any)) float64 {
	if !e.hasAttr(name) {
		return def
	}
	raw := strings.TrimSpace(e.attr(name))
	v, ok := parseLength(raw)
	if !ok {
		warn("the marker %s %q is not a length, so it is taken as %g", name, raw, def)
		return def
	}
	return v
}

// readOrient is how a marker is turned: an angle, or the keyword that turns it
// to the direction of the path at the vertex it is drawn on. An angle without a
// unit is degrees, as everywhere else in SVG. What it does not say is the
// default of no turn at all, and what it says but this package cannot read
// keeps that default with one warning, because a marker pointing the wrong way
// is one an author can see.
func readOrient(raw string, warn func(string, ...any)) (deg float64, auto, startReverse bool) {
	switch s := strings.ToLower(strings.TrimSpace(raw)); s {
	case "":
		return 0, false, false
	case "auto":
		return 0, true, false
	case "auto-start-reverse":
		return 0, true, true
	}
	v, err := parseAngle(raw)
	if err != nil {
		warn("the marker orient %q is not an angle, so the marker is drawn the way it was written", raw)
		return 0, false, false
	}
	return v, false, false
}

// resolveMarkers follows the three marker references an element named, which
// only the whole drawing can say anything about, and answers the markers that
// element draws its vertices with: the `<marker>` each id names, or one
// warning and nothing where the drawing does not have it.
//
// Each reference is cleared as it is spent, the same way a fill's is, so that a
// group says what it has to say once and not once for every shape inside it.
// What the references turn into are NOT cleared, which is how this differs from
// a clip, a mask and a filter: marker-start, marker-mid and marker-end are
// inherited, so a `<g marker-end>` is what every shape under it draws a marker
// with at its last vertex, and the children take the marker itself rather than
// the id to look up again. An element that writes a marker of its own —
// including `none`, which turns it off — takes a new one over in [Style.with],
// and that is where the old one is dropped.
func (img *Image) resolveMarkers(st *Style, warn func(string, ...any)) {
	for _, r := range []struct {
		prop string
		ref  *string
		def  **markerDef
	}{
		{"marker-start", &st.markerStartRef, &st.markerStart},
		{"marker-mid", &st.markerMidRef, &st.markerMid},
		{"marker-end", &st.markerEndRef, &st.markerEnd},
	} {
		if *r.ref == "" {
			continue
		}
		id := *r.ref
		*r.ref = ""
		def := img.markers[id]
		if def == nil {
			warn("the %s names #%s, which the drawing does not have, so no marker is drawn where it points", r.prop, id)
			continue
		}
		*r.def = def
	}
}

// cutMarkerCycles takes the marker references inside markers that would go
// round for ever. A marker whose picture has a shape with a marker of its own
// is ordinary — a marker made of markers — and is drawn as it says; one whose
// picture points back at the marker it is in the middle of is a circle, and
// drawing it would never stop. The reference that closes the circle is left off
// with one warning about it, which draws that one vertex without a marker
// rather than drawing nothing at all. It does not matter which of the three
// properties made the circle: all of them are marks on a vertex and all of them
// are walked here.
//
// A circle through a mask or a pattern needs no cutting here: those are already
// followed with everything that is being followed on the way to them, and the
// second time round the one that is already being drawn is left off — see
// [paintMasked] and [patternShade], and [Image.warnMaskPatternCycles] for
// saying which reference closes one of those. What is left over once they have
// had their say is the markers that point only at markers, which is what this
// walks.
func (img *Image) cutMarkerCycles() {
	gray := map[*markerDef]bool{}
	done := map[*markerDef]bool{}
	var walk func(n *Node)
	var visit func(def *markerDef)
	walk = func(n *Node) {
		if n == nil {
			return
		}
		for i := range n.markers {
			place := &n.markers[i]
			if place.def == nil {
				continue
			}
			switch {
			case gray[place.def]:
				img.warnings.warn(n.Name,
					"the %s here points back at the marker this picture is inside, which would draw for ever, so no marker is drawn at this vertex",
					place.prop)
				place.def = nil
			case !done[place.def]:
				visit(place.def)
			}
		}
		for _, k := range n.Kids {
			walk(k)
		}
	}
	visit = func(def *markerDef) {
		if done[def] {
			return
		}
		gray[def] = true
		walk(def.content)
		gray[def] = false
		done[def] = true
	}
	for _, def := range img.markers {
		visit(def)
	}
}

// markerVertex is one vertex of a path with the direction the path is going as
// it reaches and leaves it, which is everything `orient="auto"` needs to turn a
// marker to. hasIn and hasOut say whether there is a segment on that side at
// all: the first vertex of a path has nothing coming into it, the last has
// nothing going out, and the first point of a subpath after a move is not
// arrived at — the path does not travel from the end of one subpath to the
// start of the next.
type markerVertex struct {
	pt            canvas.Point
	hasIn, hasOut bool
	inX, inY      float64
	outX, outY    float64
}

// pathVertices is every vertex of the path, subpath by subpath, in the order
// the path draws them. The point a close travelled back to is a vertex like any
// other, which is what makes the closing vertex of a rectangle its last one and
// the same place as its first — the way SVG says the last vertex of a path that
// ends closed is its initial vertex again.
func pathVertices(path *canvas.Path) []markerVertex {
	pts, closed := path.Points()
	var vs []markerVertex
	for j, sub := range pts {
		if len(sub) == 0 {
			continue
		}
		// A subpath of one point standing after a closed one is the point that
		// close travelled back to — it is what [canvas.Path.Close] leaves
		// behind — so the path arrives at it along the segment that closed the
		// subpath before it.
		remnant := len(sub) == 1 && j > 0 && closed[j-1] && len(pts[j-1]) > 0
		for k, q := range sub {
			v := markerVertex{pt: q}
			switch {
			case k > 0:
				v.hasIn, v.inX, v.inY = true, q.X-sub[k-1].X, q.Y-sub[k-1].Y
			case remnant:
				prev := pts[j-1]
				v.hasIn, v.inX, v.inY = true, q.X-prev[len(prev)-1].X, q.Y-prev[len(prev)-1].Y
			}
			switch {
			case k+1 < len(sub):
				v.hasOut, v.outX, v.outY = true, sub[k+1].X-q.X, sub[k+1].Y-q.Y
			case closed[j] && len(sub) > 1:
				v.hasOut, v.outX, v.outY = true, sub[0].X-q.X, sub[0].Y-q.Y
			}
			vs = append(vs, v)
		}
	}
	return vs
}

// angle is the direction `orient="auto"` turns a marker to at this vertex: the
// way the path is going where there is only one segment beside it, and the
// middle of the turn where there are two. The middle is the two directions
// added as unit vectors — the bisector of the angle they make — which is what
// keeps a marker on a corner pointing along the way the path is turning rather
// than down one of the two sides of it. A path that turns straight back on
// itself has no middle to stand in, and the direction it came in from is the
// answer then.
func (v markerVertex) angle() float64 {
	switch {
	case v.hasIn && v.hasOut:
		ix, iy := unit(v.inX, v.inY)
		ox, oy := unit(v.outX, v.outY)
		sx, sy := ix+ox, iy+oy
		if sx == 0 && sy == 0 {
			return math.Atan2(iy, ix)
		}
		return math.Atan2(sy, sx)
	case v.hasOut:
		return math.Atan2(v.outY, v.outX)
	default:
		return math.Atan2(v.inY, v.inX)
	}
}

// unit is a direction of no length at all answered as no direction at all, so
// that a segment of zero length adds nothing to a bisector rather than adding
// a number divided by nothing.
func unit(x, y float64) (float64, float64) {
	if x == 0 && y == 0 {
		return 0, 0
	}
	n := math.Hypot(x, y)
	return x / n, y / n
}

// placeMarkers is where the three markers an element asked for sit on its
// shape, one entry for each marker it draws in the order the vertices come:
// the first vertex for marker-start, every vertex between the first and the
// last for marker-mid, and the last one for marker-end. A path of a single
// vertex is both its first and its last, so it takes a marker-start and a
// marker-end in the same place and no marker-mid at all.
//
// It answers nothing where the shape has no path, where no marker was named —
// which is all of the shapes in most of the drawings — or where the room a
// marker asked for is no room at all, which draws nothing where it is named.
func placeMarkers(path *canvas.Path, st Style) []markerPlace {
	if path == nil || (st.markerStart == nil && st.markerMid == nil && st.markerEnd == nil) {
		return nil
	}
	vs := pathVertices(path)
	if len(vs) == 0 {
		return nil
	}
	var out []markerPlace
	if p := markerOn(st.markerStart, st, vs[0], true, "marker-start"); p != nil {
		out = append(out, *p)
	}
	if st.markerMid != nil {
		for i := 1; i+1 < len(vs); i++ {
			if p := markerOn(st.markerMid, st, vs[i], false, "marker-mid"); p != nil {
				out = append(out, *p)
			}
		}
	}
	if p := markerOn(st.markerEnd, st, vs[len(vs)-1], false, "marker-end"); p != nil {
		out = append(out, *p)
	}
	return out
}

// markerOn is one marker on one vertex of a shape: the room the picture is
// drawn into, the fit of its viewBox, the turn `orient` gave it, and the matrix
// that puts all of that down on the vertex — with the element's own transform
// in it under the same condition the path itself got it, so that a marker lands
// where the vertex landed.
//
// The picture is taken there the way it was written: fitted into the room the
// marker asked for if it has a viewBox, moved back by the reference point,
// turned, and put down on the vertex. At the first vertex a marker written
// `auto-start-reverse` is turned around, which is what lets one marker be put
// at both ends of a line and point out of each of them; everywhere else it is
// the same as `auto`, as the spec says.
func markerOn(def *markerDef, st Style, v markerVertex, atStart bool, prop string) *markerPlace {
	if def == nil {
		return nil
	}
	// The room the picture is drawn into: the size the marker asked for, and a
	// multiple of the stroke-width of this shape where markerUnits says so —
	// the default, which is what makes a marker on a thick line a thick
	// marker. Without a viewBox the room only bounds what is drawn rather than
	// scaling it, since the picture is written in the room's own coordinates
	// already; with one, the fit below is what the two numbers are used for.
	scale := 1.0
	if def.strokeUnits {
		scale = st.Width
	}
	w, h := def.w*scale, def.h*scale
	if !(w > 0) || !(h > 0) {
		return nil
	}
	// The picture onto the room, which is the fit a viewBox gets everywhere it
	// is fitted to something — the canvas, the box a `<use>` asked a `<symbol>`
	// for — kept proportions and all.
	onto := canvas.Identity()
	if def.hasView {
		onto = fitTransform(def.view, w, h)
	}
	rx, ry := onto.Map(def.refX, def.refY)
	angle := def.deg * math.Pi / 180
	if def.auto {
		angle = v.angle()
		if atStart && def.startReverse {
			angle += math.Pi
		}
	}
	at := canvas.Translate(v.pt.X, v.pt.Y).
		Mul(canvas.Rotate(angle)).
		Mul(canvas.Translate(-rx, -ry)).
		Mul(onto)
	if st.HasTransform && st.Transform != (canvas.Matrix{}) {
		at = st.Transform.Mul(at)
	}
	// The room as a cut, in the coordinates the picture is written in: its
	// rectangle taken backwards through the fit of the viewBox, so that where
	// the room's edges land on the picture follow where the picture's own
	// writing lands. A fit with no way back collapses the picture to nothing
	// anyway, so it gets no cut rather than the wrong one.
	var room *clipRegion
	if !def.visible {
		if inv, ok := onto.Inverse(); ok {
			p := canvas.NewPath()
			p.AddRect(0, 0, w, h)
			p.Transform(inv)
			room = &clipRegion{shapes: []canvas.MaskShape{{Path: p, Rule: canvas.FillNonZero}}}
		}
	}
	return &markerPlace{def: def, at: at, room: room, prop: prop}
}

// paintMarkers draws the markers an element asked for at the vertices of its
// shape, in the order the vertices come: each marker's own picture taken to the
// canvas by the matrix that was worked out where the shape was built, and then
// by the one that puts the drawing there. The pictures are drawn the same way
// everything else is — their own fill and stroke, their own groups, their own
// clip and mask — and where the room of the marker cuts its picture, the way
// any element with a clip-path on it is drawn: through a picture of its own,
// cut to the room, laid down over what was there. The masks and patterns
// already being drawn on the way to them go with them, so that a marker inside
// a pattern inside a marker finds the pattern already there rather than drawing
// it again for ever. A marker the cycle cut left off has no marker to draw and
// is passed over.
func paintMarkers(cv *canvas.Canvas, n *Node, m canvas.Matrix, current canvas.Color, masking []*maskDef, patterning []*pattern) {
	for i := range n.markers {
		place := &n.markers[i]
		if place.def == nil || place.def.content == nil {
			continue
		}
		content := place.def.content
		if place.room != nil {
			wrap := &Node{clip: place.room}
			wrap.Kids = []*Node{content}
			content = wrap
		}
		paintNodes(cv, content, m.Mul(place.at), current, masking, patterning)
	}
}
