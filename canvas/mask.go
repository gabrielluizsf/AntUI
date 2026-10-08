package canvas

import "math"

// MaskRect zeros everything outside a rectangle, keeping only the pixels
// inside it. It is how a background painted into a layer is clipped to part of
// the box — the padding or content box — before the layer is composited.
// It honours the clip on every canvas format.
func (cv *Canvas) MaskRect(x, y, w, h int) {
	for py := cv.Clip.Y; py < cv.Clip.Y+cv.Clip.Height; py++ {
		for px := cv.Clip.X; px < cv.Clip.X+cv.Clip.Width; px++ {
			if px < x || px >= x+w || py < y || py >= y+h {
				cv.Put(px, py, 0)
			}
		}
	}
}

// MaskRoundRect zeros everything outside a rounded rectangle, the clip a
// background layer gets so the corner transparency of the box shows through
// to what lies behind. A radius of zero makes it exactly [Canvas.MaskRect].
func (cv *Canvas) MaskRoundRect(x, y, w, h, rx, ry int) {
	if w <= 0 || h <= 0 {
		return
	}
	rx = min(rx, w/2)
	ry = min(ry, h/2)
	if rx <= 0 || ry <= 0 {
		cv.MaskRect(x, y, w, h)
		return
	}
	for py := cv.Clip.Y; py < cv.Clip.Y+cv.Clip.Height; py++ {
		// The straight run between the corner bands has no curve: only the
		// columns outside the rectangle need clearing. The top and bottom
		// bands (and the rows beyond the rectangle) run the ellipse test, with
		// the area outside the rectangle cleared before it.
		if py >= y+ry && py < y+h-ry {
			for px := cv.Clip.X; px < cv.Clip.X+cv.Clip.Width; px++ {
				if px < x || px >= x+w {
					cv.Put(px, py, 0)
				}
			}
			continue
		}
		for px := cv.Clip.X; px < cv.Clip.X+cv.Clip.Width; px++ {
			if px < x || px >= x+w || py < y || py >= y+h {
				cv.Put(px, py, 0)
				continue
			}
			if !inRoundRect(px, py, x, y, w, h, rx, ry) {
				cv.Put(px, py, 0)
			}
		}
	}
}

// A MaskShape is one outline a mask is cut from: the path, and the rule that
// says which of its insides count — the same rule a fill takes, because a mask
// is a fill measured rather than painted.
type MaskShape struct {
	Path *Path
	Rule FillRule
	// And is the cuts the shape itself carries: what it keeps has to be under
	// the shape and under every one of these too, each entry being the union
	// of one clip of shapes — the clip written inside a `<clipPath>` cutting
	// the shape it was written on. What any entry covers is inside that entry,
	// and the shape is kept only where the shape and every entry agree. An
	// entry with no shapes in it keeps nothing anywhere, so a shape cut by a
	// reference the drawing could not follow keeps nothing, the same as an
	// empty `<clipPath>` keeps nothing. Nil is a shape with no cuts at all.
	And [][]MaskShape
}

// MaskShapes keeps only what lies under one of the shapes and empties the rest
// of the clip region, which is a clip: the shapes are the inside of it and
// everything else is drawn nowhere. The edge of every shape is soft — a pixel
// it covers by half keeps half of what it had — which is the antialiasing a
// rectangle cannot have and the reason this is not [Canvas.MaskRect] with a
// path in it.
//
// The shapes are a union and not a run of cuts one after the other: two of
// them overlapping leave their overlap standing, which is what a clip written
// as several shapes means. The union of no shapes at all is nothing, so no
// shapes at all empties the clip region — that is the clip an empty
// `<clipPath>` is, where nothing of what it cuts is drawn. A shape that
// carries cuts of its own (see [MaskShape.And]) has to be under them as well:
// what its own outline covers and what they cover is where it keeps.
//
// What is cut is the alpha alone: the colours stay as they are and only how
// much of each there is changes, so this is a cut to make on a layer, where
// emptiness is a colour a pixel can have. The narrow formats keep no alpha —
// what the mask empties there becomes the colour of nothing, and the soft edge
// along the way is lost, the same as it is in [Canvas.MaskRect].
func (cv *Canvas) MaskShapes(shapes ...MaskShape) {
	cx0, cy0 := cv.Clip.X, cv.Clip.Y
	cx1, cy1 := cx0+cv.Clip.Width, cy0+cv.Clip.Height
	if cx0 >= cx1 || cy0 >= cy1 {
		return
	}
	// The band is the box the shapes can reach, cut down to the clip region:
	// what is beside or beyond it is under no shape at all, and only inside it
	// is there anything to measure.
	bx0, by0, bx1, by1 := cx1, cy1, cx0, cy0
	for _, s := range shapes {
		if s.Path == nil {
			continue
		}
		minX, minY, maxX, maxY, ok := s.Path.Bounds()
		if !ok {
			continue
		}
		x0 := clampInt(int(math.Floor(minX)), cx0, cx1)
		y0 := clampInt(int(math.Floor(minY)), cy0, cy1)
		x1 := clampInt(int(math.Ceil(maxX)), cx0, cx1)
		y1 := clampInt(int(math.Ceil(maxY)), cy0, cy1)
		if x0 >= x1 || y0 >= y1 {
			continue
		}
		bx0, by0 = min(bx0, x0), min(by0, y0)
		bx1, by1 = max(bx1, x1), max(by1, y1)
	}
	if bx0 >= bx1 || by0 >= by1 {
		// Nothing reaches the clip region, so nothing in it is under a shape.
		for y := cy0; y < cy1; y++ {
			cv.cutRow(y, cx0, cx1, nil)
		}
		return
	}
	// Every shape measured over the whole band at once, along with the clips
	// each one carries and the clips those carry, all the way down. The runs
	// of one list come back in row order, so each is walked from where it
	// left off rather than being searched for again on every row.
	type measured struct {
		at    int
		spans []span
		and   [][]measured
	}
	var measure func(s MaskShape) measured
	measure = func(s MaskShape) measured {
		w := measured{spans: fillCoverage(s.Path, bx0, by0, bx1, by1, s.Rule)}
		for _, e := range s.And {
			entry := make([]measured, 0, len(e))
			for _, es := range e {
				if es.Path == nil {
					continue
				}
				entry = append(entry, measure(es))
			}
			w.and = append(w.and, entry)
		}
		return w
	}
	all := make([]measured, 0, len(shapes))
	for _, s := range shapes {
		if s.Path == nil {
			continue
		}
		all = append(all, measure(s))
	}
	row := make([]float64, bx1-bx0)
	// A row's worth of working room, kept per depth: a shape, the clip it is
	// cut by and the clip that one is cut by are never measured into the same
	// row at once, and the depth they need is as many as the clips cut into
	// clips along the way — which is never more than the picture is nested.
	var bufs [][]float64
	bufAt := func(d int) []float64 {
		for len(bufs) <= d {
			bufs = append(bufs, make([]float64, bx1-bx0))
		}
		return bufs[d]
	}
	// cov is how much of one row of the band one measured shape covers: what
	// its own outline covers, cut down to what every clip it carries covers
	// too, one of them at a time. A pixel outside one of them multiplies out
	// to nothing, which is the cut it is. It hands back the buffer it worked
	// in at depth d, for the caller to take what it needs of before the next
	// shape at that depth works in it again.
	var cov func(w *measured, d, y int) []float64
	cov = func(w *measured, d, y int) []float64 {
		dst := bufAt(d)
		clear(dst)
		for w.at < len(w.spans) && w.spans[w.at].y == y {
			s := w.spans[w.at]
			w.at++
			for x := s.x0; x < s.x1; x++ {
				if s.a > dst[x-bx0] {
					dst[x-bx0] = s.a
				}
			}
		}
		for j := range w.and {
			over := bufAt(d + 1)
			clear(over)
			for k := range w.and[j] {
				m := cov(&w.and[j][k], d+2, y)
				for x := range over {
					if m[x] > over[x] {
						over[x] = m[x]
					}
				}
			}
			for x := range dst {
				dst[x] *= over[x]
			}
		}
		return dst
	}
	for y := cy0; y < cy1; y++ {
		if y < by0 || y >= by1 {
			cv.cutRow(y, cx0, cx1, nil)
			continue
		}
		clear(row)
		for i := range all {
			m := cov(&all[i], 0, y)
			for x := range row {
				// The most of a pixel any one shape covers is how much of
				// it the union covers: a second shape over the same pixel
				// adds nothing to what is already inside the clip.
				if m[x] > row[x] {
					row[x] = m[x]
				}
			}
		}
		// The two ends of the row are past the band and so past every shape;
		// the middle of it keeps what the shapes measured.
		cv.cutRow(y, cx0, bx0, nil)
		cv.cutRow(y, bx0, bx1, row)
		cv.cutRow(y, bx1, cx1, nil)
	}
}

// A MaskMode is how the picture a mask is made of says how much of what it
// covers is kept.
type MaskMode uint8

// The two ways a mask may measure. Luminance is what an SVG `<mask>` means
// unless it says otherwise: what is bright keeps and what is dark takes away.
// Alpha keeps by how much of the mask is there at all, so an opaque black
// shape keeps everything under it and only what is clear takes away — the
// difference between the two is a black mask, which under one is a hole and
// under the other is no mask at all.
const (
	MaskLuminance MaskMode = iota
	MaskAlpha
)

// MaskBy takes the alpha of every pixel of cv down to what the picture in m
// keeps of it, which is a multiply and not a cut: where m keeps half, what cv
// holds comes out half, and where m keeps nothing nothing of it is left. That
// is what makes a mask different from [Canvas.MaskShapes], which decides
// whether a pixel is there at all — a mask says how much of it is, and a
// picture painted white over half a shape leaves that half whole and takes the
// other half away rather than putting a hard edge between them.
//
// How much is kept is the mode's business: how bright the pixel of m is, by
// the same coefficients [Canvas.FilterRegion] grays a colour with, and how much
// of that pixel is there on top of it — a clear pixel keeps nothing however
// bright the colour behind its zero alpha happens to be. In MaskAlpha only the
// second half counts.
//
// A pixel of cv that m has no say about — one beside a picture smaller than cv
// — is under no mask at all and so keeps nothing, the same as the union of no
// shapes empties the clip in [Canvas.MaskShapes]. Only the alpha moves: the
// colours are as they were and only how much of each there is changes, which
// is what makes it a call to make on a layer. The narrow formats have no
// per-pixel alpha to multiply and are left alone.
func (cv *Canvas) MaskBy(m *Canvas, mode MaskMode) {
	if cv.Pixels == nil || m == nil {
		return
	}
	for y := cv.Clip.Y; y < cv.Clip.Y+cv.Clip.Height; y++ {
		for x := cv.Clip.X; x < cv.Clip.X+cv.Clip.Width; x++ {
			keep := maskKeep(m.At(x, y), mode)
			c := cv.At(x, y)
			a := int(float64(c.A())*keep + 0.5)
			if a == int(c.A()) {
				continue
			}
			cv.Put(x, y, Fade(c, a))
		}
	}
}

// maskKeep is how much of one pixel of cv the pixel of the mask over it keeps.
//
// The brightness is taken on the sRGB values the canvas holds rather than on
// the linear ones the spec's own colour space would ask for, which is the same
// choice grayscale() makes and the same numbers it uses: a mask and a graying
// filter that disagreed about what a colour's brightness was would disagree
// about how much of a picture a grey mask keeps.
func maskKeep(c Color, mode MaskMode) float64 {
	alpha := float64(c.A()) / 255
	if mode == MaskAlpha {
		return alpha
	}
	bright := (0.2126*float64(c.R()) + 0.7152*float64(c.G()) + 0.0722*float64(c.B())) / 255
	return bright * alpha
}

// cutRow takes the alpha of a run of pixels down to what the mask keeps of it:
// keep[x-x0] is how much of the pixel at x that is, and a nil keep empties the
// whole run. Only the alpha moves, so a pixel that already holds what it would
// hold is left alone rather than written back unchanged.
func (cv *Canvas) cutRow(y, x0, x1 int, keep []float64) {
	for x := x0; x < x1; x++ {
		c := cv.At(x, y)
		a := 0
		if keep != nil {
			a = int(float64(c.A())*keep[x-x0] + 0.5)
		}
		if a == int(c.A()) {
			continue
		}
		cv.Put(x, y, Fade(c, a))
	}
}
