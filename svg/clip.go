package svg

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// A `<clipPath>` is a set of shapes that says which part of an element is
// drawn and which part of it is not: what the shapes cover keeps, and the rest
// of it is nowhere. The shapes are written once — usually at the end of the
// drawing, beside the gradients — and any number of elements point at them
// with `clip-path="url(#id)"`.
//
// Like a gradient, a clip may be written after the element that names it, so
// the whole drawing is read for them before any of it is built. Unlike a
// gradient, where the clip is written says nothing about where it cuts: the
// shapes are kept in the coordinates they were written in, and it is the
// element that points at them that moves them, because a clip cuts in the
// coordinates of whoever is being clipped.

// clipPath is one `<clipPath>` of a drawing: the picture it held, the shapes
// that picture came out as and how it said they were to be measured.
type clipPath struct {
	// shapes are what the clip is made of, each with the rule its own
	// `clip-rule` asked for. The clip is the union of them, so every shape is
	// measured on its own and what any one of them covers is inside — two
	// shapes overlapping leave their overlap standing rather than cutting it
	// away again.
	shapes []canvas.MaskShape
	// boxUnits is whether the shapes were written as fractions of the box of
	// the element being clipped rather than in the drawing's own coordinates,
	// which is the second of the two things `clipPathUnits` may say. See
	// [Image.resolveClip].
	boxUnits bool
	// content is the tree the shapes are read from, kept because a `clip-path`
	// written inside a `<clipPath>` cuts the shape it is written on, and the
	// shapes of the outer one can only be worked out once every clipPath in
	// the drawing is known — one of them may be named further down the file
	// than the place it is pointed at. See [Image.readClipPaths].
	content *Node
	// built says the shapes above are finished: the picture in content was
	// walked and every `clip-path` inside it was followed into them. A
	// reference met while this is still false is one from inside a picture
	// still being read, and waits to be followed — see [Image.resolveClip]
	// and [Image.resolveClipsIn].
	built bool
}

// readClipPaths takes every `<clipPath>` that named itself, under the id a
// `clip-path` may point at, and works out the shapes it holds. One written
// without an id cannot be pointed at, so it is left unread the way a gradient
// with no id is: nothing in the drawing can reach it.
func (img *Image) readClipPaths(root *element) {
	type declared struct {
		id   string
		e    *element
		rule canvas.FillRule
	}
	var found []declared
	seen := map[string]bool{}
	// The clip-rule in force walks down with the elements, because it is
	// inherited and a clipPath may stand anywhere in the drawing — inside a
	// group that wrote it once for everything below it. It is read on the way
	// down rather than by building the elements above the clipPath, because
	// those are part of the picture too and every one of them is built and
	// read once already: reading them again here would say everything twice.
	var walk func(e *element, rule canvas.FillRule)
	walk = func(e *element, rule canvas.FillRule) {
		rule = clipRuleIn(e, rule)
		if e.Name == "clipPath" {
			if id := e.attr("id"); id != "" && !seen[id] {
				seen[id] = true
				found = append(found, declared{id: id, e: e, rule: rule})
			}
		}
		for _, k := range e.Kids {
			walk(k, rule)
		}
	}
	walk(root, canvas.FillNonZero)

	img.clips = make(map[string]*clipPath, len(found))
	// Every clipPath is named before any of them is built, so that one
	// written inside another — or pointed at from inside one — finds a place
	// to land rather than finding nothing there at all.
	for _, d := range found {
		img.clips[d.id] = &clipPath{}
	}
	for _, d := range found {
		img.clips[d.id] = img.readClipPath(d.e, d.rule)
	}
	// Every clipPath holds its picture now, and every `clip-path` written
	// inside one of them is still a name pointing at a clip whose shapes are
	// not worked out yet. They are worked out here, one id at a time,
	// following what is inside each first: a clip written inside a clipPath
	// cuts the shape it is written on, and it may name one written further
	// down the file, which can only be followed now that every id is known
	// and every picture is built. The walk goes round a circle at the first
	// repeat rather than for ever — see [Image.measureClip].
	gray := map[string]bool{}
	for _, d := range found {
		img.measureClip(d.id, gray)
	}
}

// measureClip works out the shapes one clipPath cuts with, following the
// `clip-path`s written inside it first, and hands back the clipPath with
// them in it. The gray map is the walk: an id in it is a clipPath whose
// shapes are being worked out right now, so a reference back at one of them
// is a circle the walk would go round for ever without. A clipPath whose
// shapes are already worked out is handed back as it stands, which is what
// stops the walk from reading the same picture twice — and saying everything
// inside it twice with it.
func (img *Image) measureClip(id string, gray map[string]bool) *clipPath {
	cp := img.clips[id]
	if cp == nil || cp.built {
		return cp
	}
	gray[id] = true
	img.resolveClipsIn(cp.content, gray)
	cp.shapes = img.clipShapesOf(cp.content)
	cp.built = true
	delete(gray, id)
	return cp
}

// resolveClipsIn follows every `clip-path` still waiting in a picture: they
// were named while the picture was built, when the clipPath they point at
// had a name but no shapes yet, and now that it has, the shapes are moved
// into the place of whoever named them the same way [Image.resolveClip] does
// for a reference anywhere else in the drawing.
func (img *Image) resolveClipsIn(n *Node, gray map[string]bool) {
	if n == nil {
		return
	}
	if n.clip != nil && n.clip.pending != "" {
		n.clip = img.clipNow(n, n.clip.pending, gray)
	}
	for _, k := range n.Kids {
		img.resolveClipsIn(k, gray)
	}
}

// clipNow is the clip one node inside a `<clipPath>` takes while the shapes
// of the clipPath it is inside are being worked out: the shapes of the
// clipPath it names, moved into the node's own place the same way
// [Image.resolveClip] moves them. A name that points back at a clipPath
// whose shapes are being worked out right now is a circle, and it is left
// off with one warning saying so — a clip left off keeps everything, so the
// shape it was written on stands whole inside the clipPath, which is what
// lets the walk come back rather than going round for ever.
func (img *Image) clipNow(n *Node, id string, gray map[string]bool) *clipRegion {
	warn := func(format string, args ...any) { img.warnings.warn(n.Name, format, args...) }
	if gray[id] {
		warn("the clip-path here points back at the <clipPath> this picture is inside, which would go round for ever, so the clip is left off")
		return nil
	}
	cp := img.measureClip(id, gray)
	if cp == nil {
		// Not there: resolveClip would have said so and kept nothing while
		// the picture was built, so a name only reaches here from a clipPath
		// that was standing when it was read.
		return &clipRegion{}
	}
	return clipUnder(img.clipRegionOf(cp, n.Style), n.Style, n, warn)
}

// readClipPath reads one `<clipPath>`: how it said its shapes were to be
// measured, and the picture they are to be read out of. The picture is built
// in the coordinates it was written in — the transform on the clipPath and on
// anything inside it counts, since that is part of how it was written — and
// the shapes themselves are worked out afterwards, once every clipPath in the
// drawing is named; see [Image.readClipPaths]. The style it is built from
// starts at the drawing's own rather than at the style of the place the
// clipPath stands, because a clipPath takes nothing from the elements around
// it except the clip-rule it was written under.
//
// The children are built rather than left out, which is what forceKids says:
// what is inside a definition is not part of the picture where it stands —
// that is what being a definition means — but here it is not being painted
// either, it is being read as the shape of a cut.
func (img *Image) readClipPath(e *element, inherited canvas.FillRule) *clipPath {
	base := baseStyle()
	base.ClipRule = inherited
	units := strings.ToLower(strings.TrimSpace(e.attr("clipPathUnits")))
	return &clipPath{
		boxUnits: units != "" && units != "userspaceonuse",
		content:  img.build(e, base, true),
	}
}

// clipRuleIn is the clip-rule one element stands under: the one written on it
// — as an attribute or as a declaration in its `style` — or the one it
// inherited from above, which is what makes clip-rule an inherited property
// like the rest of the style. It is read without a warn because there is
// nothing here that cannot be read: a rule that is not `evenodd` is simply
// the nonzero one, which is also what it starts as.
func clipRuleIn(e *element, inherited canvas.FillRule) canvas.FillRule {
	if e.hasAttr("clip-rule") {
		return clipRuleOf(e.attr("clip-rule"))
	}
	if s := e.attr("style"); s != "" {
		if raw, ok := parseDeclarations(s)["clip-rule"]; ok {
			return clipRuleOf(raw)
		}
	}
	return inherited
}

// clipRuleOf is one value of it.
func clipRuleOf(raw string) canvas.FillRule {
	if strings.EqualFold(strings.TrimSpace(raw), "evenodd") {
		return canvas.FillEvenOdd
	}
	return canvas.FillNonZero
}

// overflowIn is whether an element lets what it draws run past the room it
// asked for — the room of a `<marker>` or the tile of a `<pattern>`, which are
// the two places the spec says `overflow` stands. Only `visible` says so:
// anything else, including nothing at all, says the room cuts what runs past
// its edge, which is what the spec's own stylesheet writes on both of them.
// It is read the way clip-rule is — an attribute first, then a declaration in
// the `style` — and, like clip-rule, read without a warn: a value that is not
// `visible` is simply the cut.
func overflowIn(e *element) bool {
	if e.hasAttr("overflow") {
		return strings.EqualFold(strings.TrimSpace(e.attr("overflow")), "visible")
	}
	if s := e.attr("style"); s != "" {
		if raw, ok := parseDeclarations(s)["overflow"]; ok {
			return strings.EqualFold(strings.TrimSpace(raw), "visible")
		}
	}
	return false
}

// clipShapesOf is what a built `<clipPath>` adds to the clip: the outline of
// every shape inside it, read with the rule it asked for. Everything that
// draws nothing adds nothing — a group is its children, a definition is
// nothing until something points at it, and what is hidden is nowhere — and
// the one thing that has an outline but no shape to cut with says so, because
// a drawing that asked to cut with writing is a drawing that is not going to
// get what it asked for.
//
// A `clip-path` written on the way down to a shape cuts that shape: each one
// found on the way is carried the rest of the way and put on the shape, so
// what comes out is a shape that keeps where it keeps and where every clip
// over it keeps too. The list of them is copied as it is carried down, so
// that two shapes side by side do not append to the same list behind each
// other's back and end up cutting each other's clips away.
func (img *Image) clipShapesOf(n *Node) []canvas.MaskShape {
	var shapes []canvas.MaskShape
	var take func(n *Node, cut [][]canvas.MaskShape)
	take = func(n *Node, cut [][]canvas.MaskShape) {
		if n.clip != nil {
			// The clip is carried even when it holds nothing: a reference the
			// drawing could not follow keeps nothing, so the shape it was
			// written on falls out of the clipPath, the same as a clip of
			// nothing anywhere else in the drawing.
			cut = append(append([][]canvas.MaskShape(nil), cut...), append([]canvas.MaskShape(nil), n.clip.shapes...))
		}
		switch {
		case n.Path != nil && !n.Path.Empty():
			shapes = append(shapes, canvas.MaskShape{Path: n.Path, Rule: n.Style.ClipRule, And: cut})
		case len(n.Runs) > 0:
			img.warnings.warn(n.Name,
				"writing inside a <clipPath> has no outline to cut with here, so it is left out of the clip")
		case n.Pic != nil:
			img.warnings.warn(n.Name,
				"a picture inside a <clipPath> has no outline to cut with here, so it is left out of the clip")
		default:
			for _, k := range n.Kids {
				take(k, cut)
			}
		}
	}
	take(n, nil)
	return shapes
}

// resolveClip follows the `clip-path` an element named, which only the whole
// drawing can say anything about, and answers the clip that element is cut by:
// the shapes of the clipPath it names, moved into the element's own place
// before anyone paints, since that is where the element is going. The shapes
// in the drawing are shared by everything that points at them, so each takes a
// copy rather than moving the ones everyone else is holding.
//
// A reference the drawing cannot follow keeps nothing and says so, which is
// what the rest of this package does with a reference to nowhere: a `<use>`
// that names nothing draws nothing, a fill that names a gradient that is not
// there is not filled, and a clip that names a clipPath that is not there cuts
// everything away.
//
// A reference from inside a `<clipPath>` cannot be followed here at all: every
// clipPath is being built at this point and none of them has its shapes yet —
// one of them may be named further down the file than the place it is pointed
// at. It waits as the name alone, on the node it was written on, until the
// picture it is written in is walked for its shapes and the name can be
// followed for real — see [Image.resolveClipsIn].
//
// The reference is cleared as it is spent, the same way a fill's is, so that a
// `<g clip-path>` says what it has to say once and not once for every shape
// inside it, and so that the children never see a clip of their own to
// inherit — the picture they are drawn into is what gets cut.
func (img *Image) resolveClip(st *Style, warn func(string, ...any)) *clipRegion {
	if st.clipRef == "" {
		return nil
	}
	id := st.clipRef
	st.clipRef = ""
	cp := img.clips[id]
	switch {
	case cp == nil:
		warn("the clip-path names #%s, which the drawing does not have, so nothing of what it cuts is drawn", id)
		return &clipRegion{}
	case !cp.built:
		return &clipRegion{pending: id}
	}
	return img.clipRegionOf(cp, *st)
}

// clipRegionOf is the shapes of one clipPath moved into the place of the
// element that named it: with the element's own transform put on them, so
// that the cut follows the element it cuts, or left as they were written
// where the shapes are written in fractions of the element's own box and it
// is [clipUnder] that puts them where the box is. Every copy carries the
// clips cut into the shapes along with it — see [mapMaskShapes].
func (img *Image) clipRegionOf(cp *clipPath, st Style) *clipRegion {
	if cp.boxUnits {
		// The clipPath uses objectBoundingBox units: the shapes are fractions
		// of the element's box and are moved by clipUnder, once the element is
		// built and the box is known — not by the element's own transform,
		// which would put them somewhere before they were ever mapped.
		return &clipRegion{shapes: mapMaskShapes(cp.shapes, canvas.Identity()), boxUnits: true}
	}
	m := canvas.Identity()
	if st.HasTransform && st.Transform != (canvas.Matrix{}) {
		// The element's own transform is on the shapes now, the same way it is
		// put on the element's own path while it is built, so that what cuts
		// the element moves with the element. Everything that points at one
		// clipPath takes its own copy for exactly this reason: the drawing is
		// read once and may be painted at many sizes and in many places.
		m = st.Transform
	}
	return &clipRegion{shapes: mapMaskShapes(cp.shapes, m)}
}

// clipUnder puts a clip written in fractions of the box of the element it
// cuts where that box is: one matrix takes every shape of it — and every clip
// a shape carries with it — from fractions into the drawing's own coordinates,
// the same matrix because the shapes were all written against one box and
// moving each to fit wherever its own outline happens to reach would pull
// them apart from each other and from what they cut. A box that cannot be
// measured or has a side of nothing leaves the element drawn whole, with one
// warning: an element drawn whole is one an author can see, and a warning is
// what tells them it was not clipped.
func clipUnder(clip *clipRegion, st Style, n *Node, warn func(string, ...any)) *clipRegion {
	if clip == nil {
		return nil
	}
	if !clip.boxUnits {
		return clip
	}
	box, ok := paintedBox(n)
	if !ok || box[2] <= 0 || box[3] <= 0 {
		warn("the clipPath is written in fractions of the box of this element, which has no box to measure, so the element is drawn without it")
		return nil
	}
	m := canvas.Translate(box[0], box[1]).Mul(canvas.Scale(box[2], box[3]))
	return &clipRegion{shapes: mapMaskShapes(clip.shapes, m)}
}

// mapMaskShapes is every shape of a clip taken through m and back out as a
// shape of its own, the cuts a shape carries taken through it alongside: a
// shape written inside a `<clipPath>` may carry a clip of its own, and the
// two only stay in the same places if the same move takes both. Nothing is
// shared with the shapes it came from, which is the other half of it: the
// shapes in the drawing are read once and moved into the place of every
// element that points at them.
func mapMaskShapes(shapes []canvas.MaskShape, m canvas.Matrix) []canvas.MaskShape {
	out := make([]canvas.MaskShape, len(shapes))
	for i, s := range shapes {
		out[i] = canvas.MaskShape{Rule: s.Rule}
		if s.Path != nil {
			p := clonePath(s.Path)
			p.Transform(m)
			out[i].Path = p
		}
		if len(s.And) > 0 {
			out[i].And = make([][]canvas.MaskShape, len(s.And))
			for j, e := range s.And {
				out[i].And[j] = mapMaskShapes(e, m)
			}
		}
	}
	return out
}

// clipRegion is the clip one element carries: the shapes of the clipPath it
// named, with its own transform already put on them so that the cut follows
// the element it cuts. It hangs off the node rather than off the style because
// a clip is not inherited (see the `clip` field of [Node]), and it can be
// empty, which is a clip that keeps nothing rather than no clip at all.
type clipRegion struct {
	// boxUnits is whether the shapes were written as fractions of the box of
	// the element being clipped rather than in the drawing's own coordinates.
	boxUnits bool
	shapes   []canvas.MaskShape
	// pending is the id of a clipPath whose shapes are not worked out yet: a
	// `clip-path` written inside a `<clipPath>`, named while every clipPath in
	// the drawing was being built and before any of them was measured. It
	// carries the name alone and no shapes; [Image.resolveClipsIn] follows it
	// when the picture it was written in is walked for its own shapes. See
	// [Image.resolveClip].
	pending string
}

// measured is the region of a node against the canvas it is being painted on:
// the shapes carry the element's own transform already, and the matrix that
// maps the drawing onto the canvas is what is left to put on them.
func (r *clipRegion) measured(m canvas.Matrix) []canvas.MaskShape {
	return mapMaskShapes(r.shapes, m)
}
