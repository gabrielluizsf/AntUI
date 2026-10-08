package svg

import (
	"math"
	"slices"
	"strings"
	"sync"

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
// drawn once for every shape that asks the same of it rather than once per
// shape, and once per pixel of neither, which is what keeps a pattern no dearer
// than a gradient to paint with.
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
	// visible is `overflow`: whether the picture may run out of the tile. What
	// it does not say — which is what the spec's own stylesheet says a pattern
	// stands for — is that it may not, so the tile cuts what runs past its
	// edge unless `overflow="visible"` is written, and then the picture is
	// drawn into the tile from the copies of it that land there. See
	// [pattern.picture].
	visible bool
	// tiles is the picture this pattern already drew, kept for the shapes that
	// ask the same of it — the same tile at the same place and size, in the
	// same colour, over the same chain of patterns already being drawn, which
	// is the whole of what a picture is drawn from. Ten shapes under one
	// pattern cost one picture rather than ten, and drawn says how many were
	// drawn rather than found kept, which is what a test reads to know the
	// keeping holds — a picture reused and a picture drawn again come out the
	// same on the page. mu guards both because two goroutines may paint the
	// same drawing at once; the picture is drawn with mu off and only the
	// finding and the keeping of it are under mu, so a pattern inside a pattern
	// never waits on the one holding it — and where two draw the same picture
	// at once, whichever keeps it first stands and the other's is let go rather
	// than both being held.
	mu    sync.Mutex
	tiles []tilePicture
	drawn int
}

// tilePicture is one kept picture of a pattern: what it was drawn from, the
// chain of patterns being drawn when it was — kept because the same tile
// reached through two different chains may come out differently, a picture
// naming a pattern that one chain has already open is cut off in that one —
// and the picture itself, which is only ever read after it was drawn.
type tilePicture struct {
	key   tileKey
	chain []*pattern
	pic   *canvas.Canvas
}

// tileKey is everything a pattern's picture is drawn from apart from the
// chain: the tile where it stands and how big it lands, the colour carried
// into it, and the box of the shape — but only where the picture is written
// in fractions of that box, since a picture in the drawing's own coordinates
// never looks at it and shapes of different sizes share one picture.
type tileKey struct {
	box           [4]float64
	x, y, w, h, s float64
	current       canvas.Color
}

// maxTilePixels is how many pixels one tile may be drawn into. A tile the size
// of the drawing at the drawing's own scale is ordinary and worth every pixel
// of it; a tile written in fractions of a box on a shape drawn enormous is not
// worth running the machine out of memory for. What is worth it is what the
// drawing can reach: such a tile is drawn as the part of it the canvas can ask
// about, at the scale the pixels are sampled at, and comes out smaller and
// blurrier only where even that part is over the limit — the tiling is still
// the tiling, and a pattern is a texture rather than the thing itself.
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
	// Read after any href above, so that what this element says about its own
	// room is what stands rather than what the pattern it borrowed from said.
	p.visible = overflowIn(e)
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
func patternShade(pt *pattern, n *Node, m canvas.Matrix, width, height int, current canvas.Color, opacity float64, patterning []*pattern) (func(x, y int) canvas.Color, bool) {
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
	// that moves the pattern, since that transform moves the tile too. A tile
	// worth more than [maxTilePixels] comes out drawn as the part of it the
	// canvas can ask about at this same scale, and smaller only where even
	// that is over the limit.
	s := scaleOf(toCanvas.Mul(pt.transform))
	if !(s > 0) || math.IsInf(s, 0) {
		return nil, false
	}
	// Where the picture is drawn into: the whole tile, unless the whole tile
	// at this scale is dearer than [maxTilePixels] is worth — then the part of
	// the tile the canvas corners fall in, where even that fits, drawn at the
	// same scale the pixels are sampled at; and where it does not, the whole
	// tile at a scale that does, which is what every tile over the limit used
	// to come out as. A pattern that draws past its own edges (an
	// overflow=visible one) is tied to the tile's edges rather than to where
	// the canvas looks, so it always takes the last road.
	rx, ry, rw, rh := x, y, w, h
	if w*s*h*s > maxTilePixels {
		if !pt.visible {
			// The box the four corners of the canvas fall in, taken back
			// through the same two maps every point goes through: everything
			// outside it is never asked about. A side the box runs past the
			// tile on is the whole of the tile that way — a box straddling
			// the tile's edge wraps round onto both sides of it.
			u0, v0, u1, v1 := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
			for _, c := range [][2]float64{{0, 0}, {float64(width), 0}, {0, float64(height)}, {float64(width), float64(height)}} {
				ux, uy := inv.Map(c[0], c[1])
				bx, by := toTile.Map(ux, uy)
				u0, v0 = min(u0, bx-x), min(v0, by-y)
				u1, v1 = max(u1, bx-x), max(v1, by-y)
			}
			if u0 >= 0 && u1 <= w {
				rx, rw = x+u0, u1-u0
			}
			if v0 >= 0 && v1 <= h {
				ry, rh = y+v0, v1-v0
			}
		}
		if rw*s*rh*s > maxTilePixels {
			// Even the part the canvas reaches is over the limit, or the
			// pattern draws past its own edges: the whole tile smaller.
			s *= math.Sqrt(maxTilePixels / (w * s * h * s))
			rx, ry, rw, rh = x, y, w, h
		}
	}
	// The picture is cut to the part above, so a pixel takes its colour from
	// the tile at the same scale minus how far into the tile the part starts.
	xoff, yoff := (rx-x)*s, (ry-y)*s
	tile := pt.picture(box, rx, ry, rw, rh, s, current, patterning)
	if tile == nil {
		return nil, false
	}
	return func(px, py int) canvas.Color {
		ux, uy := inv.Map(float64(px)+0.5, float64(py)+0.5)
		bx, by := toTile.Map(ux, uy)
		ix := int(modUnit(bx-x, w)*s - xoff)
		iy := int(modUnit(by-y, h)*s - yoff)
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

// picture is one tile at the size it lands at on the page, drawn into a canvas
// of its own — and kept, so that the shapes asking the same of the pattern
// share it rather than every one of them drawing it again. What it is drawn
// from is the whole of the key it is kept under: the tile where it stands and
// how big it lands, the colour carried into it, the box of the shape where the
// picture is written in fractions of that box, and the chain of patterns that
// was already being drawn — the same tile reached through two chains may come
// out differently, a picture naming a pattern one chain has open being cut off
// in that one. It answers nil where there is nothing to draw it with, which is
// the same as a picture that keeps nothing, and a nil is kept the same as any
// other.
func (pt *pattern) picture(box [4]float64, x, y, w, h, s float64, current canvas.Color, patterning []*pattern) *canvas.Canvas {
	key := tileKey{x: x, y: y, w: w, h: h, s: s, current: current}
	if pt.contentBox {
		key.box = box
	}
	pt.mu.Lock()
	for _, t := range pt.tiles {
		if t.key == key && slices.Equal(t.chain, patterning) {
			pic := t.pic
			pt.mu.Unlock()
			return pic
		}
	}
	pt.drawn++
	pt.mu.Unlock()
	pic := pt.drawPicture(box, x, y, w, h, s, current, patterning)
	pt.keepTile(key, patterning, pic)
	return pic
}

// keepTile keeps a picture drawn under this key and chain for whatever asks
// the same of the pattern later — and where a picture drawn in parallel got
// there first, lets go of this one so that only one of the two is held. The
// chain is copied: the slice it comes in belongs to the caller, which goes on
// appending to it as the drawing goes deeper.
func (pt *pattern) keepTile(key tileKey, chain []*pattern, pic *canvas.Canvas) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	for _, t := range pt.tiles {
		if t.key == key && slices.Equal(t.chain, chain) {
			return
		}
	}
	pt.tiles = append(pt.tiles, tilePicture{key: key, chain: slices.Clone(chain), pic: pic})
}

// drawPicture is one tile drawn into a canvas of its own, at the size it lands
// at on the page: the content of the pattern put through the coordinates it
// was written in. What runs out of the tile is cut by the canvas being exactly
// the tile — which is what `overflow` on a pattern says unless it says
// `visible`, and then the picture is painted into the tile from every copy of
// it whose overflow reaches in, so a picture written wider than its tile lands
// in the tiles beside it instead of being cut off along the edge. It answers
// nil where there is nothing to draw it with, which is the same as a picture
// that keeps nothing.
func (pt *pattern) drawPicture(box [4]float64, x, y, w, h, s float64, current canvas.Color, patterning []*pattern) *canvas.Canvas {
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
	// the same as a gradient in objectBoundingBox units does. A copy of the
	// picture reaches into the tile from between the two transforms: the offset
	// that brings it in is counted in the tile's own units, so it has to be
	// laid after the box has had its say or the box scales it away.
	base := canvas.Scale(s, s).Mul(canvas.Translate(-x, -y))
	boxT := canvas.Identity()
	if pt.contentBox {
		boxT = canvas.Translate(box[0], box[1]).Mul(canvas.Scale(box[2], box[3]))
	}
	one := func(at canvas.Matrix) {
		paintNodes(layer, pt.content, at, current, nil, append(patterning, pt))
	}
	if pt.visible {
		// `overflow="visible"`: the copies of the picture at the offsets that
		// bring its overflow into the tile, painted left to right in rows from
		// the top — the order the spec says tiles are laid in, which is what
		// says which of two overlapping copies is on top. The picture's span
		// in the tile's own coordinates comes out of what it painted, through
		// the same box the picture itself went through.
		if b, ok := paintedBox(pt.content); ok {
			x0, y0, x1, y1 := b[0], b[1], b[0]+b[2], b[1]+b[3]
			if pt.contentBox {
				x0, x1 = box[0]+x0*box[2], box[0]+x1*box[2]
				y0, y1 = box[1]+y0*box[3], box[1]+y1*box[3]
			}
			if nx0, nx1, ok := copyRange(x0-x, x1-x, w); ok {
				if ny0, ny1, ok := copyRange(y0-y, y1-y, h); ok {
					for ny := ny0; ny <= ny1; ny++ {
						for nx := nx0; nx <= nx1; nx++ {
							one(base.Mul(canvas.Translate(float64(nx)*w, float64(ny)*h)).Mul(boxT))
						}
					}
					return layer
				}
			}
		}
		// The picture reaches nothing into the tile, or reaches so far past it
		// that drawing it over and over would be dearer than the corner of it
		// that would come out: one copy, cut to the tile, which is the default.
	}
	one(base.Mul(boxT))
	return layer
}

// copyRange is which copies of a picture written over [lo, hi) — in the
// coordinates of the tile, however far out those run — reach into the tile
// itself, [0, size): the copy at the offset n*size, for each n from the first
// to the last of them. ok answers false where no copy reaches in at all, and
// also where they reach so far out that there would be too many of them to
// draw one by one — at most five along one way, which is a picture spanning
// about five tiles, and anything beyond that takes the plain cut instead.
func copyRange(lo, hi, size float64) (first, last int, ok bool) {
	f := math.Floor(-hi/size) + 1
	l := math.Ceil((size-lo)/size) - 1
	if f > l || l-f >= 5 || math.Abs(f) > 1e6 || math.Abs(l) > 1e6 {
		return 0, 0, false
	}
	return int(f), int(l), true
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
