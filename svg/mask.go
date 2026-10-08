package svg

import (
	"strconv"
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// A `<mask>` is a picture that says how much of another element is there:
// where the mask is bright the element stands, where it is dark the element is
// nowhere, and every shade between the two leaves that much of it. Where a
// `<clipPath>` is a cut with an edge — a pixel is either kept or it is not —
// a mask is a multiply, and that is what lets a drawing fade an element away
// or cut it with a soft edge rather than only with a hard one.
//
// Like a clip, a mask is written once — usually at the end of the drawing,
// beside the gradients and the clip paths — and any number of elements point
// at it with `mask="url(#id)"`. Like a clip, it may be written after the
// element that names it, so every `<mask>` in the drawing is read before any
// of it is built. Unlike a clip, where the mask is written says nothing about
// where it reaches either: the picture is kept in the coordinates it was
// written in, and it is the element that names it that moves it, because a
// mask is laid over what is being masked.

// maskDef is one `<mask>` of a drawing: the picture it holds, how that picture
// says how much it keeps, and where it reaches.
type maskDef struct {
	// content is the mask's own tree, built in the coordinates it was written
	// in — the transform on the mask and on anything inside it counts, since
	// that is part of how they were written — and it is the element that names
	// the mask that moves it afterwards. It is nil only for [maskNowhere], the
	// mask a reference that the drawing cannot follow stands for.
	content *Node
	// measure is how the content says how much of what it covers is kept: by
	// how bright it is, which is what a mask means where nothing has said
	// otherwise, or by how much of it is there, which is what `mask-type:
	// alpha` asks for.
	measure canvas.MaskMode
	// region is where the mask reaches, and boxUnits inside it says which of
	// the two ways `maskUnits` has of writing that it was written in. See
	// [readMaskRegion] and [maskUnder].
	region maskRegion
	// contentBoxUnits is whether the content was written in fractions of the
	// box of the element the mask is put on rather than in the drawing's own
	// coordinates, which is the second of the two things `maskContentUnits`
	// may say. It is read as a whole or not at all — see [maskNamed].
	contentBoxUnits bool
}

// maskNowhere is the mask a reference the drawing cannot follow stands for: a
// picture with nothing in it, which keeps nothing, the same as a clip that
// names a clipPath that is not there. It is shared because it holds nothing
// that belongs to any one drawing.
var maskNowhere = &maskDef{}

// maskRegion is where a `<mask>` reaches: a rectangle, either in fractions of
// the box of the element it is put on or in the drawing's own coordinates,
// which is what `maskUnits` chooses between. What is written against neither
// is the SVG default of the box grown by a tenth all round, filled in when the
// mask is read. What the mask says beyond its region is beyond the region
// rather than under it, so it keeps nothing there — the region is where the
// mask is allowed to say anything at all.
type maskRegion struct {
	// boxUnits is whether x, y, w and h are fractions of the box of the
	// element the mask is put on rather than lengths of the drawing's own.
	boxUnits   bool
	x, y, w, h float64
}

// masked is the mask one node is painted under. It hangs off the node rather
// than off the style, because a mask is not inherited — a group is masked once
// as a whole, when its picture is made, rather than masked again at every
// child as it is drawn — which is the same as a clip. See the `mask` field of
// [Node].
type masked struct {
	// def is the mask that was named, which is [maskNowhere] where the drawing
	// could not follow the reference: that keeps nothing, rather than keeping
	// everything, so that a drawing which asked for a mask it does not have
	// shows the element nowhere rather than showing it whole.
	def *maskDef
	// x, y, w, h is where the mask reaches, in the drawing's own coordinates,
	// and cuts says the region narrows it at all: a region measured against a
	// box that cannot be measured leaves the mask reaching everywhere. See
	// [maskUnder].
	x, y, w, h float64
	cuts       bool
}

// readMasks takes every `<mask>` that named itself, under the id a `mask` may
// point at, and builds the picture each one holds. They are named before any
// of them is built, the way the clip paths are, because a mask may be written
// after the element that names it and one mask may be named from inside
// another — neither can be followed before every id is known, and a mask whose
// picture points at one written further down the file would otherwise find
// nothing there at all.
func (img *Image) readMasks(root *element) {
	type declared struct {
		id string
		e  *element
	}
	var found []declared
	seen := map[string]bool{}
	var walk func(e *element)
	walk = func(e *element) {
		if e.Name == "mask" {
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

	img.masks = make(map[string]*maskDef, len(found))
	for _, d := range found {
		img.masks[d.id] = &maskDef{}
	}
	// What a mask says about itself comes before any picture is built, because
	// the region is measured the moment an element that names the mask is
	// built — and one mask may name another, so a picture built before every
	// region is known would be measured against a region that was not filled
	// in yet, which is a zero rectangle reaching nothing, and would cut away
	// everything it was put on.
	for _, d := range found {
		img.readMask(d.id, d.e)
	}
	for _, d := range found {
		img.readMaskContent(d.id, d.e)
	}
}

// readMask reads one `<mask>`'s own sayings — the way it measures and where it
// reaches — into the place its id already put, before any picture is built.
// See [Image.readMasks] for why it is a separate step.
func (img *Image) readMask(id string, e *element) {
	warn := func(format string, args ...any) { img.warn(e, format, args...) }
	def := img.masks[id]
	def.measure = maskTypeIn(e, warn)
	units := strings.ToLower(strings.TrimSpace(e.attr("maskContentUnits")))
	def.contentBoxUnits = units != "" && units != "userspaceonuse"
	def.region = img.readMaskRegion(e, warn)
}

// readMaskContent builds the picture one `<mask>` holds, which is what forceKids
// says: what is inside a definition is not part of the picture where it stands,
// but here it is not being painted either, it is being read as the measure of
// something else.
func (img *Image) readMaskContent(id string, e *element) {
	img.masks[id].content = img.build(e, baseStyle(), true)
}

// readMaskRegion reads where a `<mask>` reaches: `maskUnits` says which of the
// two coordinates the four numbers are written in, and `x`, `y`, `width` and
// `height` are the rectangle itself. What is not written is the SVG default of
// the box grown by a tenth all round, which is the same numbers either way —
// a fraction of the box, or of the drawing when the region is written in the
// drawing's own coordinates.
func (img *Image) readMaskRegion(e *element, warn func(string, ...any)) maskRegion {
	units := strings.ToLower(strings.TrimSpace(e.attr("maskUnits")))
	r := maskRegion{boxUnits: units != "userspaceonuse", x: -0.1, y: -0.1, w: 1.2, h: 1.2}
	vb := img.ViewBox
	if !r.boxUnits {
		r.x, r.y, r.w, r.h = -0.1*vb[2], -0.1*vb[3], 1.2*vb[2], 1.2*vb[3]
	}
	for _, f := range []struct {
		name string
		dst  *float64
		span float64
	}{
		{"x", &r.x, vb[2]},
		{"y", &r.y, vb[3]},
		{"width", &r.w, vb[2]},
		{"height", &r.h, vb[3]},
	} {
		if !e.hasAttr(f.name) {
			continue
		}
		v, ok := maskRegionValue(e.attr(f.name), r.boxUnits, f.span)
		if !ok {
			warn("the mask region %s %q is not a number this package can read, so it is taken as the whole of it", f.name, e.attr(f.name))
			continue
		}
		*f.dst = v
	}
	return r
}

// maskRegionValue is one number of a mask's region: a percentage, which is a
// fraction of the box of the element when the region is written in fractions
// of it and a fraction of the drawing when it is not; or a number, which in
// the first case is a fraction on its own and in the second is a length of the
// drawing's own. A length with units is not a fraction of anything, so it is
// only read where the region is written in the drawing's coordinates.
func maskRegionValue(raw string, boxUnits bool, span float64) (float64, bool) {
	s := strings.TrimSpace(raw)
	if p, ok := strings.CutSuffix(s, "%"); ok {
		v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
		if err != nil {
			return 0, false
		}
		if boxUnits {
			return v / 100, true
		}
		return v / 100 * span, true
	}
	if boxUnits {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, false
		}
		return v, true
	}
	return parseLength(s)
}

// maskTypeIn is how one `<mask>` measures how much it keeps: `luminance`,
// which is what a mask means where nothing has said otherwise, or `alpha`,
// which keeps by how much of the mask is there rather than by how bright it
// is. It is read from the attribute and from a declaration in the `style`,
// the way clip-rule is, and it is not inherited — a mask is measured by what
// it says about itself and not by the place it stands in.
//
// A value that is neither is said out loud and read as the default, because
// the two measures disagree about the whole of a black mask and an author who
// asked for one of them should not be given the other in silence.
func maskTypeIn(e *element, warn func(string, ...any)) canvas.MaskMode {
	raw, ok := e.attr("mask-type"), e.hasAttr("mask-type")
	if !ok {
		if s := e.attr("style"); s != "" {
			raw, ok = parseDeclarations(s)["mask-type"]
		}
	}
	if !ok {
		return canvas.MaskLuminance
	}
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "alpha":
		return canvas.MaskAlpha
	case "luminance", "":
		return canvas.MaskLuminance
	default:
		warn("the mask-type %q is not luminance or alpha, so the mask measures how bright it is", raw)
		return canvas.MaskLuminance
	}
}

// maskNamed is the mask an element named, followed before anything is built —
// where a clip is followed — so that what cannot be followed says so once,
// while the element is being read rather than while it is being painted.
//
// Three things come out of it: nothing at all, where no mask was named or the
// element asks for `mask="none"`; the mask itself; and [maskNowhere], where
// the drawing does not have the id that was named, which keeps nothing. A
// mask whose content is written in fractions of a box this package does not
// measure is left off altogether with one warning saying so, which draws the
// element whole: an element drawn whole is one an author can see, and a
// warning is what tells them it was not masked.
func (img *Image) maskNamed(ref string, warn func(string, ...any)) *maskDef {
	if ref == "" {
		return nil
	}
	def := img.masks[ref]
	switch {
	case def == nil:
		warn("the mask names #%s, which the drawing does not have, so nothing of what it masks is drawn", ref)
		return maskNowhere
	case def.contentBoxUnits:
		warn("the mask is written in objectBoundingBox units, which are not read here, so the element is drawn without it")
		return nil
	}
	return def
}

// maskUnder is what a node carries for the mask it was found to be under: the
// region the mask reaches, worked out against the box of everything the node
// paints — which is only known once everything inside it is there, and is why
// the reference is named before the node is built and this is done after.
//
// The region of a mask written in the drawing's own coordinates is moved by
// the element's transform, the same as the shapes of a clip are: what is
// written against this element is written in its coordinates, and a transform
// on it moves the whole of it. The region of one written in fractions of the
// box needs no moving, because the box it is a fraction of is the box the
// element paints at, transform and all.
func maskUnder(def *maskDef, st Style, n *Node) *masked {
	if def == nil {
		return nil
	}
	m := &masked{def: def}
	r := def.region
	switch {
	case !r.boxUnits:
		m.x, m.y, m.w, m.h = r.x, r.y, r.w, r.h
		if st.HasTransform && st.Transform != (canvas.Matrix{}) {
			m.x, m.y, m.w, m.h = movedBox(m.x, m.y, m.w, m.h, st.Transform)
		}
		m.cuts = true
	default:
		box, ok := paintedBox(n)
		if !ok {
			// A region measured against a box that cannot be measured is not
			// one to cut with: a mask reaching everywhere is one an author can
			// see the whole of, and one reaching nowhere is a drawing that
			// vanished with nothing to say about it. Writing is measured the
			// same as a shape is now, so what is left with no box of its own
			// is a node with nothing in it that paints at all.
			return m
		}
		m.x, m.y = box[0]+r.x*box[2], box[1]+r.y*box[3]
		m.w, m.h = r.w*box[2], r.h*box[3]
		m.cuts = true
	}
	return m
}

// paintedBox is the box a node paints over, in the drawing's own coordinates:
// every shape it draws and every shape inside it draws, each where it was put
// — the shapes carry their transforms already, so this is the box as it lands
// on the page rather than as it was written. A picture has no transform put
// into it — the transform meets it while it is painted — so its box is taken
// through the element's own here for the same result. Writing has no box of
// its own until a font is asked, which is what [walkRuns] and [writingBox] do
// where this walks, and the letters land through the element's transform the
// same way a picture's box does. A node with nothing in it that paints has no
// box, and says so rather than answering an empty one.
func paintedBox(n *Node) ([4]float64, bool) {
	var minMax [4]float64
	found := false
	add := func(x0, y0, x1, y1 float64) {
		if !found {
			minMax = [4]float64{x0, y0, x1, y1}
			found = true
			return
		}
		minMax[0] = min(minMax[0], x0)
		minMax[1] = min(minMax[1], y0)
		minMax[2] = max(minMax[2], x1)
		minMax[3] = max(minMax[3], y1)
	}
	var walk func(n *Node)
	walk = func(n *Node) {
		if n.Path != nil && !n.Path.Empty() {
			if x0, y0, x1, y1, ok := n.Path.Bounds(); ok {
				add(x0, y0, x1, y1)
			}
		}
		if n.Pic != nil {
			x, y, w, h := n.Pic.Box[0], n.Pic.Box[1], n.Pic.Box[2], n.Pic.Box[3]
			if n.Style.HasTransform && n.Style.Transform != (canvas.Matrix{}) {
				x, y, w, h = movedBox(x, y, w, h, n.Style.Transform)
			}
			add(x, y, x+w, y+h)
		}
		if len(n.Runs) > 0 {
			pieces := walkRuns(n.Runs, canvas.Identity())
			// Every run hidden and every run empty leaves no pieces, and a box
			// of nothing is not one to measure a region against: the same no
			// box at all as a node with no writing in it.
			if len(pieces) > 0 {
				wb := writingBox(pieces)
				x0, y0, x1, y1 := wb[0], wb[1], wb[2], wb[3]
				if x1 > x0 || y1 > y0 {
					if n.Style.HasTransform && n.Style.Transform != (canvas.Matrix{}) {
						x, y, w, h := movedBox(x0, y0, x1-x0, y1-y0, n.Style.Transform)
						x0, y0, x1, y1 = x, y, x+w, y+h
					}
					add(x0, y0, x1, y1)
				}
			}
		}
		for _, k := range n.Kids {
			walk(k)
		}
	}
	walk(n)
	if !found {
		return [4]float64{}, false
	}
	return [4]float64{minMax[0], minMax[1], minMax[2] - minMax[0], minMax[3] - minMax[1]}, true
}

// movedBox is a rectangle taken through a transform: the four corners go and
// the box around where they land is what comes back, which is exact for the
// transforms that keep a rectangle a rectangle and never smaller than the
// truth for the ones that turn it into a diamond.
func movedBox(x, y, w, h float64, m canvas.Matrix) (float64, float64, float64, float64) {
	x0, y0 := m.Map(x, y)
	x1, y1 := m.Map(x+w, y)
	x2, y2 := m.Map(x, y+h)
	x3, y3 := m.Map(x+w, y+h)
	minX, maxX := min(x0, x1, x2, x3), max(x0, x1, x2, x3)
	minY, maxY := min(y0, y1, y2, y3), max(y0, y1, y2, y3)
	return minX, minY, maxX - minX, maxY - minY
}
