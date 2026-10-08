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

// clipPath is one `<clipPath>` of a drawing: the shapes it held and how it said
// they were to be measured.
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
}

// readClipPaths takes every `<clipPath>` that named itself, under the id a
// `clip-path` may point at, and builds the shapes it holds. One written
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
}

// readClipPath builds the shapes one `<clipPath>` holds. They are built in the
// coordinates they were written in — the transform on the clipPath and on
// anything inside it counts, since that is part of how they were written —
// and it is the element that points at them that moves them afterwards. The
// style they are built from starts at the drawing's own rather than at the
// style of the place the clipPath stands, because a clipPath takes nothing
// from the elements around it except the clip-rule it was written under.
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
		shapes:   img.clipShapesOf(img.build(e, base, true)),
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

// clipShapesOf is what a built `<clipPath>` adds to the clip: the outline of
// every shape inside it, read with the rule it asked for. Everything that
// draws nothing adds nothing — a group is its children, a definition is
// nothing until something points at it, and what is hidden is nowhere — and
// the one thing that has an outline but no shape to cut with says so, because
// a drawing that asked to cut with writing is a drawing that is not going to
// get what it asked for.
func (img *Image) clipShapesOf(n *Node) []canvas.MaskShape {
	var shapes []canvas.MaskShape
	var take func(n *Node)
	take = func(n *Node) {
		switch {
		case n.Path != nil && !n.Path.Empty():
			shapes = append(shapes, canvas.MaskShape{Path: n.Path, Rule: n.Style.ClipRule})
		case len(n.Runs) > 0:
			img.warnings.warn(n.Name,
				"writing inside a <clipPath> has no outline to cut with here, so it is left out of the clip")
		case n.Pic != nil:
			img.warnings.warn(n.Name,
				"a picture inside a <clipPath> has no outline to cut with here, so it is left out of the clip")
		default:
			for _, k := range n.Kids {
				take(k)
			}
		}
	}
	take(n)
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
// everything away. A clip that asks for something this package cannot measure
// — the fractions of a box rather than the drawing's own coordinates — is left
// off altogether, with the same said about it: an element drawn whole is one
// an author can see, and a warning is what tells them it was not clipped.
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
	case cp.boxUnits:
		// The clipPath uses objectBoundingBox units: shapes are written as
		// fractions of the element's box. We store the flag and defer the
		// mapping to clipUnder, which runs after the node is built and the
		// element's box is known. The shapes here are in the coordinates they
		// were written in (before the element's own transform).
		shapes := make([]canvas.MaskShape, len(cp.shapes))
		for i, s := range cp.shapes {
			p := clonePath(s.Path)
			// Do NOT apply st.Transform here — the element transform will
			// be applied by clipUnder after we know the painted box, so the
			// fractions map consistently onto the element's own bounding box.
			shapes[i] = canvas.MaskShape{Path: p, Rule: s.Rule}
		}
		return &clipRegion{shapes: shapes, boxUnits: true}
	}
	shapes := make([]canvas.MaskShape, len(cp.shapes))
	for i, s := range cp.shapes {
		// The element's own transform is on the shapes now, the same way it is
		// put on the element's own path while it is built, so that what cuts
		// the element moves with the element. Everything that points at one
		// clipPath takes its own copy for exactly this reason: the drawing is
		// read once and may be painted at many sizes and in many places.
		p := clonePath(s.Path)
		if st.HasTransform && st.Transform != (canvas.Matrix{}) {
			p.Transform(st.Transform)
		}
		shapes[i] = canvas.MaskShape{Path: p, Rule: s.Rule}
	}
	return &clipRegion{shapes: shapes, boxUnits: false}
}

// clipUnder maps a clipRegion with boxUnits onto the element's painted box.
// If the box cannot be measured, it warns and returns nil (the element is drawn whole).
func clipUnder(clip *clipRegion, st Style, n *Node, warn func(string, ...any)) *clipRegion {
	if clip == nil {
		return nil
	}
	if !clip.boxUnits {
		return clip
	}
	box, ok := paintedBox(n)
	if !ok {
		warn("the clipPath is written in objectBoundingBox units, which are not read here for this element, so the element is drawn without it")
		return nil
	}
	// Map each shape's fractions onto the element's painted box in parent space.
	// paintedBox gives the box in parent coordinates (including the element's own transform).
	// Fractions of the box are: x = box.x + frac.x * box.w, etc.
	shapes := make([]canvas.MaskShape, len(clip.shapes))
	for i, s := range clip.shapes {
		p := clonePath(s.Path)
		// Transform the shape by translating/scaling to the painted box.
		// The shape's original coords are in the element's user space; we map
		// them so that 0→box.x, 1→box.x+box.w, etc. This is equivalent to:
		// Translate(box.x, box.y).Scale(box.w, box.h).
		// Use a matrix that scales by the box dimensions and translates by the box origin.
		minX, minY, maxX, maxY, okBox := s.Path.Bounds()
		if !okBox {
			// Fallback: keep the shape as-is (should not happen for valid paths).
			shapes[i] = canvas.MaskShape{Path: p, Rule: s.Rule}
			continue
		}
		// Scale and translate the path into the painted box.
		t := canvas.Translate(box[0], box[1])
		scl := canvas.Scale(box[2]/(maxX-minX), box[3]/(maxY-minY))
		m := t.Mul(scl)
		p.Transform(m)
		shapes[i] = canvas.MaskShape{Path: p, Rule: s.Rule}
	}
	return &clipRegion{shapes: shapes}
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
}

// measured is the region of a node against the canvas it is being painted on:
// the shapes carry the element's own transform already, and the matrix that
// maps the drawing onto the canvas is what is left to put on them.
func (r *clipRegion) measured(m canvas.Matrix) []canvas.MaskShape {
	out := make([]canvas.MaskShape, len(r.shapes))
	for i, s := range r.shapes {
		p := clonePath(s.Path)
		p.Transform(m)
		out[i] = canvas.MaskShape{Path: p, Rule: s.Rule}
	}
	return out
}
