package svg

import (
	"math"
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// A `<pattern>` is a paint rather than a colour, the same as a gradient:
// `fill="url(#p)"` names one elsewhere in the drawing, and the shape is painted
// by asking the pattern what it has at each point rather than being painted one
// flat colour. What it has is a picture, written once, laid over the shape
// again and again from one edge of it to the other — which is what makes a
// chequered fill or a field of dots out of one small drawing.
//
// The picture is drawn into a canvas the size of one tile and kept there. Every
// pixel of the shape asks that canvas which part of the tile it is over, and
// the answer wraps round by the size of the tile, so the same picture comes
// back as many times as the shape is big enough to want it. The tile itself is
// drawn once per shape rather than once per pixel, which is what keeps a pattern
// no dearer than a gradient to paint with.
//
// Three things say where the tile and its picture are, and each is read the way
// a gradient's geometry is: `patternUnits` says whether the tile is written in
// fractions of the box of the shape it paints or in the drawing's own
// coordinates — the box being the default — and `patternContentUnits` says the
// same for what is inside the tile, whose default is the other way round, the
// drawing's own. `patternTransform` moves the tile and everything in it
// without moving the shape it paints.
type pattern struct {
	// boxUnits is whether the tile below is written in fractions of the box of
	// the shape the pattern paints, which is what `patternUnits` chooses
	// between and what it says when it says nothing.
	boxUnits bool
	// contentBox is the same for the picture inside the tile, which is what
	// `patternContentUnits` chooses between — and there the drawing's own
	// coordinates are the default, not the fractions.
	contentBox bool
	// x, y, w, h is the tile, in whichever units boxUnits says: where it
	// starts and how big it is, before the pattern's own transform moves it.
	x, y, w, h float64
	// transform is `patternTransform`: it moves the tile and the picture in it,
	// the same as `gradientTransform` moves a gradient, and the shape the
	// pattern paints stays where it is. It is always a matrix that can be put
	// back, so an empty one is the identity rather than nothing.
	transform canvas.Matrix
	// content is the picture inside the tile, built where it was written — the
	// transform of the pattern and of anything inside it counts, because that
	// is part of how they were written — and drawn into the tile over and over
	// as the shape is painted. It is not part of the picture where it stands,
	// which is what forceKids says when it is built.
	content *Node
}

// maxTilePixels is how many pixels one tile may be drawn into. A tile the size
// of the drawing at the drawing's own scale is ordinary and worth every pixel
// of it; a tile written in fractions of a box on a shape drawn enormous is not
// worth running the machine out of memory for, and comes out smaller and
// blurrier instead — the tiling is still the tiling, and a pattern is a texture
// rather than the thing itself.
const maxTilePixels = 2_000_000

// readPatterns takes every `<pattern>` that named itself, under the id a `fill`
// or a `stroke` may point at, and builds the picture each one holds. They are
// named before any of them is built, the way the gradients and the masks are:
// a pattern may be written after the element that names it, and what is inside
// one may name anything else the drawing has — a gradient, a clip, a mask, a
// `<use>`, or another pattern written further down the file.
func (img *Image) readPatterns(root *element) {
	type declared struct {
		id string
		e  *element
	}
	var found []declared
	seen := map[string]bool{}
	var walk func(e *element)
	walk = func(e *element) {
		if e.Name == "pattern" {
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

	img.patterns = make(map[string]*pattern, len(found))
	// Every pattern is named with an empty picture in it first, and only then
	// is any of them read and built: a pattern's picture may point at another
	// by an id that has not been reached yet, and a reference to a pattern that
	// is not there is a warning rather than a picture.
	for _, d := range found {
		img.patterns[d.id] = &pattern{transform: canvas.Identity()}
	}
	for _, d := range found {
		img.readPattern(d.id, d.e)
	}
	for _, d := range found {
		img.patterns[d.id].content = img.build(d.e, baseStyle(), true)
	}
}

// readPattern reads where one `<pattern>`'s tile is and how it is moved —
// everything about it except the picture it holds, which is built afterwards by
// [Image.readPatterns]. What it does not say is the SVG default, which is the
// box of the shape it paints for the tile and the drawing's own coordinates for
// what is inside it.
func (img *Image) readPattern(id string, e *element) {
	warn := func(format string, args ...any) { img.warn(e, format, args...) }
	p := img.patterns[id]
	p.boxUnits = !strings.EqualFold(strings.TrimSpace(e.attr("patternUnits")), "userspaceonuse")
	p.contentBox = strings.EqualFold(strings.TrimSpace(e.attr("patternContentUnits")), "objectboundingbox")

	// The tile before anything says where it is: nothing at all of the drawing,
	// or the whole of it in the drawing's own coordinates. A percentage is of
	// the box in the first case and of the drawing in the second, which is the
	// same two worlds a gradient's numbers live in.
	span := [2]float64{img.ViewBox[2], img.ViewBox[3]}
	p.x, p.y = 0, 0
	p.w, p.h = 1, 1
	if !p.boxUnits {
		p.w, p.h = span[0], span[1]
	}
	for _, f := range []struct {
		name string
		dst  *float64
		axis int
	}{
		{"x", &p.x, 0},
		{"y", &p.y, 1},
		{"width", &p.w, 0},
		{"height", &p.h, 1},
	} {
		if !e.hasAttr(f.name) {
			continue
		}
		raw := e.attr(f.name)
		unit := 1.0
		if !p.boxUnits {
			unit = span[f.axis]
		}
		v, ok := gradientLength(raw, unit)
		if !ok {
			warn("the pattern %s %q is not a length, so the tile keeps %g", f.name, raw, *f.dst)
			continue
		}
		*f.dst = v
	}
	if v := strings.TrimSpace(e.attr("patternTransform")); v != "" {
		if m, ok := parseTransform(v); ok {
			p.transform = m
		} else {
			warn("the pattern transform %q is not one this package can read, so the tile stays where it is", v)
		}
	}
	for _, name := range []string{"href", "xlink:href"} {
		if ref, ok := hashRef(e.attr(name)); ok && ref != "" {
			// Follow the reference to another pattern, the same way gradients
			// follow href references. The referenced pattern must be named with an
			// id and be within this drawing.
			id := strings.TrimPrefix(ref, "#")
			if refPattern, ok := img.patterns[id]; ok {
				// Use the referenced pattern's content and transform,
				// inheriting its boxUnits and contentBox settings.
				p.content = refPattern.content
				p.transform = refPattern.transform
				p.boxUnits = refPattern.boxUnits
				p.contentBox = refPattern.contentBox
				warn("the pattern takes its content from #%s", id)
			} else {
				warn("the pattern references #%s, which this drawing does not have, so its content is empty", id)
			}
			break
		}
	}
	if e.hasAttr("viewBox") {
		vbStr := strings.TrimSpace(e.attr("viewBox"))
		if vb, ok := parseViewBox(vbStr); ok {
			// viewBox establishes the pattern's coordinate system, overriding
			// x, y, width, height and patternUnits/patternContentUnits.
			p.x, p.y, p.w, p.h = vb[0], vb[1], vb[2], vb[3]
			p.boxUnits = true
			// patternContentUnits is also ignored per spec when viewBox is given,
			// but we keep the current behavior for backward compatibility.
			if e.hasAttr("patternContentUnits") {
				p.contentBox = strings.EqualFold(
					strings.TrimSpace(e.attr("patternContentUnits")), "objectboundingbox")
			}
			warn("the pattern has a viewBox, which sets the tile coordinates; " +
				"x, y, width and height attributes are ignored")
		} else {
			warn("the pattern has an invalid viewBox %q, so its picture is drawn where it was written", vbStr)
		}
	} else if e.hasAttr("viewBox") {
		// viewBox present but invalid — fall through to warning below
	}
	if p.w <= 0 || p.h <= 0 {
		warn("the pattern tile is %g by %g, which tiles nothing, so the shape is painted without it", p.w, p.h)
	}
}

// patternShade is the colour one pixel of a pattern-painted shape takes: the
// tile is drawn into a picture of its own, and the pixel asks that picture
// which part of the tile it is over — and the same part is asked for every
// pixel all the way across the shape, wrapped round by the size of the tile.
//
// The pixel is taken back to the coordinates the shape was written in, the same
// as for a gradient, and then backwards through `patternTransform` — because
// that transform moves the tile rather than the shape, so where the tile now
// stands is what has to be found. What comes out is a place inside the tile,
// worked out before the transform moved it, which is where the picture was
// drawn. A transform with no way back, a tile of no size and a pattern that is
// already being drawn further out all answer the same way: there is no shade
// here, and the shape takes the colour after the `url(...)` if it brought one
// with it and nothing at all if it did not.
func patternShade(pt *pattern, n *Node, m canvas.Matrix, current canvas.Color, opacity float64, patterning []*pattern) (func(x, y int) canvas.Color, bool) {
	for _, being := range patterning {
		if being == pt {
			// This pattern is already being drawn further up: following it
			// again would draw it while it is being drawn, for ever — a
			// pattern that names itself, or two that name each other, is left
			// off so the drawing comes out as a picture rather than not at all.
			return nil, false
		}
	}
	box := n.box
	x, y, w, h := pt.x, pt.y, pt.w, pt.h
	if pt.boxUnits {
		if box[2] <= 0 || box[3] <= 0 {
			return nil, false
		}
		x, y = box[0]+x*box[2], box[1]+y*box[3]
		w, h = w*box[2], h*box[3]
	}
	if pt.contentBox && (box[2] <= 0 || box[3] <= 0) {
		return nil, false
	}
	if w <= 0 || h <= 0 {
		return nil, false
	}
	toCanvas := m.Mul(n.Style.Transform)
	inv, ok := toCanvas.Inverse()
	if !ok {
		return nil, false
	}
	toTile, ok := pt.transform.Inverse()
	if !ok {
		return nil, false
	}
	// How much of the canvas one step of the tile is worth, which is the size
	// the tile lands at: the drawing's own scale, put through the transform
	// that moves the pattern, since that transform moves the tile too. What is
	// more pixels than a tile is worth is a tile drawn smaller rather than a
	// machine run out of room — see [maxTilePixels].
	s := scaleOf(toCanvas.Mul(pt.transform))
	if !(s > 0) || math.IsInf(s, 0) {
		return nil, false
	}
	if area := w * s * h * s; area > maxTilePixels {
		s *= math.Sqrt(maxTilePixels / area)
	}
	tile := pt.picture(box, x, y, w, h, s, current, patterning)
	if tile == nil {
		return nil, false
	}
	return func(px, py int) canvas.Color {
		ux, uy := inv.Map(float64(px)+0.5, float64(py)+0.5)
		bx, by := toTile.Map(ux, uy)
		ix := int(modUnit(bx-x, w) * s)
		iy := int(modUnit(by-y, h) * s)
		if ix < 0 {
			ix = 0
		} else if ix >= tile.Width {
			ix = tile.Width - 1
		}
		if iy < 0 {
			iy = 0
		} else if iy >= tile.Height {
			iy = tile.Height - 1
		}
		return fade(tile.At(ix, iy), opacity)
	}, true
}

// picture is one tile drawn into a canvas of its own, at the size it lands at
// on the page: the content of the pattern put through the coordinates it was
// written in, and clipped to the tile by the canvas being exactly the tile —
// which is what `overflow` on a pattern means without having to be read. It
// answers nil where there is nothing to draw it with, which is the same as a
// picture that keeps nothing.
func (pt *pattern) picture(box [4]float64, x, y, w, h, s float64, current canvas.Color, patterning []*pattern) *canvas.Canvas {
	px, py := int(math.Ceil(w*s)), int(math.Ceil(h*s))
	if px < 1 {
		px = 1
	}
	if py < 1 {
		py = 1
	}
	layer, err := canvas.NewLayer(px, py)
	if err != nil {
		return nil
	}
	// The picture's own coordinates onto the tile: a point of the picture goes
	// to where it stands in the tile and from there to a pixel of it. What is
	// written in fractions of the box of the shape goes through that box first,
	// the same as a gradient in objectBoundingBox units does.
	m := canvas.Scale(s, s).Mul(canvas.Translate(-x, -y))
	if pt.contentBox {
		m = m.Mul(canvas.Translate(box[0], box[1])).Mul(canvas.Scale(box[2], box[3]))
	}
	paintNodes(layer, pt.content, m, current, nil, append(patterning, pt))
	return layer
}

// modUnit is how far along a tile a point is, wrapped into [0, size): the
// subtraction that makes a pattern repeat, done so that a point before the
// first tile lands in the last one rather than nowhere. A tile is never no
// size here — a tile that was has already said so and paints nothing.
func modUnit(a, size float64) float64 {
	v := math.Mod(a, size)
	if v < 0 {
		v += size
	}
	return v
}
