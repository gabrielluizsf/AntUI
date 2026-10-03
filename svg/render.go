package svg

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// supersample is how many pixels to either side of the one being asked for a
// drawing is painted at before being brought back down. Drawing an icon at
// twenty-four pixels from a drawing of curves and diagonals leaves the edges
// ragged at one pixel; drawn at forty-eight they are smooth, and the halving
// averages four of them into one, which is what the antialiasing of a fill
// already does and the same cost.
const supersample = 2

// Render paints the drawing into a canvas of the given size and answers it. The
// size is in pixels and is independent of the size the drawing says it is: an
// icon drawn at twenty-four and asked for at forty-eight comes out at
// forty-eight, with everything about it twice as big.
//
// A drawing that paints itself in `currentColor` is painted in black, which is
// what the keyword means where nothing else has said what the colour is. Use
// [Image.RenderWith] to name one.
//
// The result is remembered, so painting the same drawing at the same size again
// answers the canvas that was already made rather than drawing it a second time.
// That is what makes an icon that is on screen every frame cost nothing after
// the first one. Painting at a new size draws it again and replaces what was
// remembered, so the cache holds one size and no more.
func (img *Image) Render(width, height int) *canvas.Canvas {
	return img.RenderWith(width, height, canvas.RGB(0, 0, 0))
}

// RenderWith is [Image.Render] with the colour that a `currentColor` paint
// takes. It is what lets one drawing be written once and painted in as many
// colours as the interface it belongs to has states: the same tick in the
// normal, hovered and pressed colour of a checkbox, without being written three
// times.
//
// The colour is part of what is remembered, so asking for the same drawing at
// the same size in two colours draws it twice, and asking again in either of
// them costs nothing. A caller that wants an icon to hold a while and be
// repainted often is better off asking for one colour at a time.
func (img *Image) RenderWith(width, height int, current canvas.Color) *canvas.Canvas {
	if img == nil || img.Root == nil || width <= 0 || height <= 0 {
		return nil
	}
	// The one size and colour already drawn, and the canvas it was drawn on. A
	// call for the same one hands the same canvas back; a caller that draws on
	// it is drawing on the icon itself, which is the point.
	img.mu.Lock()
	if img.drawn != nil && img.drawnWidth == width && img.drawnHeight == height && img.drawnColour == current {
		cv := img.drawn
		img.mu.Unlock()
		return cv
	}
	img.mu.Unlock()

	cv := img.paint(width, height, current)

	img.mu.Lock()
	img.drawn, img.drawnWidth, img.drawnHeight, img.drawnColour = cv, width, height, current
	img.mu.Unlock()
	return cv
}

// paint is the drawing itself: a canvas twice the size asked for, the nodes
// walked in order with the transform that maps the drawing's own coordinates
// onto it, and the result brought back down to size.
func (img *Image) paint(width, height int, current canvas.Color) *canvas.Canvas {
	// A layer rather than a plain canvas, because an icon is a picture with soft
	// edges that gets laid over a widget, not a surface of its own. A plain
	// canvas composites over what is there and answers opaque, which on an empty
	// one means a quarter-covered edge is stored as a quarter-bright opaque
	// pixel: a dark rim all the way round the icon.
	big, err := canvas.NewLayer(width*supersample, height*supersample)
	if err != nil {
		return nil
	}
	// The viewBox maps onto the whole of the drawing, and the drawing maps onto
	// the canvas, so the two compose into the transform every node is drawn
	// with. A drawing with a viewBox that does not start at the origin is
	// shifted, which is what keeps a viewBox of "0 0 24 24" and one of
	// "-12 -12 24 24" both filling the icon the same way.
	m := fitTransform(img.ViewBox, float64(width*supersample), float64(height*supersample))
	paintNodes(big, img.Root, m, current, nil, nil)

	if supersample == 1 {
		return big
	}
	// Scaled averages the block each pixel covers weighted by alpha and writes
	// the result rather than compositing it, so what comes down keeps the
	// transparency the drawing had. Halving by Blit would blend instead and
	// flatten every edge to solid.
	return canvas.Scaled(big, width, height)
}

// paintNodes walks the tree, drawing each node and then its children inside it.
// The matrix is the one the node's own transform composed onto the one it
// inherits, so a child is drawn where its parent put it.
//
// masking is the masks whose pictures are being drawn on the way to this node,
// and patterning the patterns whose tiles are being drawn on the way to it. A
// mask may name one of them — itself, or two that name each other — and
// following that would draw the same picture while it is already being drawn,
// for ever; the mask is left off instead, which is what keeps a broken drawing
// a picture rather than a program that does not come back. A pattern is left
// off the same way, and the shape that asked for it takes the colour after the
// `url(...)` if it brought one.
func paintNodes(cv *canvas.Canvas, n *Node, m canvas.Matrix, current canvas.Color, masking []*maskDef, patterning []*pattern) {
	if n.Style.Hidden {
		return
	}
	if n.clip != nil || n.mask != nil || len(n.filters) > 0 {
		paintLayered(cv, n, m, current, masking, patterning)
		return
	}
	paintBody(cv, n, m, current, masking, patterning)
}

// paintBody is one node and everything inside it, drawn straight onto the
// canvas it is given: its shape first, then its writing, then the picture an
// `<image>` is, then its children in the order they were written.
func paintBody(cv *canvas.Canvas, n *Node, m canvas.Matrix, current canvas.Color, masking []*maskDef, patterning []*pattern) {
	if n.Path != nil && !n.Path.Empty() {
		paintShape(cv, n, m, current, masking, patterning)
	}
	if len(n.Runs) > 0 {
		paintText(cv, n, m, current)
	}
	if n.Pic != nil {
		paintPicture(cv, n, m)
	}
	for _, k := range n.Kids {
		paintNodes(cv, k, m, current, masking, patterning)
	}
}

// paintLayered draws a node into a picture of its own, puts that picture
// through the filter list the element asked for, cuts the result down to the
// shapes its clip-path named, takes what its mask says of what is left over,
// and lays it over what is already there. The picture has to be made before it
// is touched, whole: a group is one picture with one filter, one clip and one
// mask over it, and cutting each of its children as it was drawn would cut
// them apart from each other, so that two translucent shapes overlapping under
// the edge of the clip would meet somewhere the clip does not keep at all.
//
// The order the picture goes through them is the order the spec puts them in:
// the filter first, because it changes the picture itself and reaches a little
// past where the element stopped; then the clip, which says where of that
// picture is kept; then the mask, which says how much of what is left is
// there. A clip or a mask measured before the filter would cut away the reach
// of a blur that was still coming, and the filter after them would soften the
// cut they made instead of the picture they were cut from.
//
// The picture is the size of the whole canvas rather than of the node, because
// a clip may reach anywhere and the cut has to be of the picture where the clip
// says and nowhere else. What a `<clipPath>` keeps nothing of comes out of this
// the same way: a picture drawn, cut to nothing, and laid over what was there
// as nothing at all — and so does a mask whose reference the drawing could not
// follow, which keeps nothing rather than keeping everything.
func paintLayered(cv *canvas.Canvas, n *Node, m canvas.Matrix, current canvas.Color, masking []*maskDef, patterning []*pattern) {
	if n.mask != nil && n.mask.def.content == nil {
		// There is no picture to make: the mask keeps nothing, so nothing of
		// what it masks is drawn.
		return
	}
	layer, err := canvas.NewLayer(cv.Width, cv.Height)
	if err != nil {
		// There is no picture to make — a canvas of no size has nothing to
		// draw on and nothing to draw into either.
		return
	}
	paintBody(layer, n, m, current, masking, patterning)
	if len(n.filters) > 0 {
		// The scale the element is drawn at goes with the list because what
		// the group's transform drew is bigger or smaller than the writing
		// says: a blur written as 1 is 1 of the drawing's units wherever those
		// units landed on the canvas. The transform is composed the same way
		// the children were drawn with it.
		applyFilters(layer, n.filters, scaleOf(m.Mul(n.Style.Transform)))
	}
	if n.clip != nil {
		layer.MaskShapes(n.clip.measured(m)...)
	}
	if n.mask != nil {
		paintMasked(layer, n, m, current, masking, patterning)
	}
	cv.BlitOver(0, 0, layer)
}

// paintMasked takes what the element painted down to what its mask keeps of
// it: the mask's picture is drawn into a picture of its own, at the same size
// and over the same transform so that what it says lines up with what it is
// being said about, and then every pixel of the element is multiplied by how
// much of the mask is over it.
//
// The mask's picture is drawn under the element's own transform, because a
// mask — like the shapes a clip cuts with — is written in the coordinates of
// the element it is put on: a `transform` on that element moves the mask along
// with everything else it moves. Its region is measured in the drawing's own
// coordinates and taken to the canvas by the same transform the element is,
// for the same reason.
func paintMasked(layer *canvas.Canvas, n *Node, m canvas.Matrix, current canvas.Color, masking []*maskDef, patterning []*pattern) {
	md := n.mask
	for _, being := range masking {
		if being == md.def {
			// This mask is already being drawn further up, so following it
			// again would draw it while it is being drawn, for ever: a mask
			// that names itself, or two that name each other, is left off and
			// the element stands as it is rather than the drawing never
			// coming back.
			return
		}
	}
	ml, err := canvas.NewLayer(layer.Width, layer.Height)
	if err != nil {
		// With no picture of the mask there is nothing to say how much is
		// kept, and what no mask has anything to say about keeps nothing.
		layer.MaskShapes()
		return
	}
	if md.def.content != nil {
		paintNodes(ml, md.def.content, m.Mul(n.Style.Transform), current, append(masking, md.def), patterning)
	}
	if md.cuts {
		x0, y0, x1, y1 := regionArea(md.x, md.y, md.w, md.h, m)
		layer.MaskRect(x0, y0, x1-x0, y1-y0)
	}
	layer.MaskBy(ml, md.def.measure)
}

// regionArea is where a mask reaches on the canvas: the four corners of the
// region go through the transform that puts the drawing on it, and the box
// around where they land is what is cut to — exact for the transforms that
// keep a rectangle a rectangle, and never smaller than the truth for the ones
// that do not, which is what a region has to be.
func regionArea(x, y, w, h float64, m canvas.Matrix) (x0, y0, x1, y1 int) {
	ax, ay := m.Map(x, y)
	bx, by := m.Map(x+w, y)
	cx, cy := m.Map(x, y+h)
	dx, dy := m.Map(x+w, y+h)
	minX := math.Floor(min(ax, bx, cx, dx))
	minY := math.Floor(min(ay, by, cy, dy))
	maxX := math.Ceil(max(ax, bx, cx, dx))
	maxY := math.Ceil(max(ay, by, cy, dy))
	return int(minX), int(minY), int(maxX), int(maxY)
}

// paintShape draws one node's shape: its fill first and its stroke over it,
// which is the order SVG paints, so a stroke half its width over a fill covers
// the fill's edge as it should, and the marker over the pair of them last of
// all — an arrowhead on the line rather than under it.
func paintShape(cv *canvas.Canvas, n *Node, m canvas.Matrix, current canvas.Color, masking []*maskDef, patterning []*pattern) {
	st := n.Style
	if st.HasFill {
		switch {
		case st.fillGradient != nil:
			// A gradient is a colour at every point rather than one colour, so
			// the fill is laid one pixel at a time. The path is still transformed
			// first: the gradient moves with the shape it paints, and what is
			// asked about each pixel is where that pixel falls in the shape's own
			// coordinates rather than where it falls on the canvas.
			fill := clonePath(n.Path)
			fill.Transform(m)
			cv.FillPathFunc(fill, st.FillRule, gradientShade(st.fillGradient, n, m, current, st.Opacity*st.FillOpacity))
		case st.fillPattern != nil:
			// A pattern is a picture at every point rather than one colour, so
			// the fill is laid one pixel at a time the same way. Where the
			// pattern has nothing to paint with — it names itself, the tile has
			// no size to it, the transform has no way back — the colour after
			// the `url(...)` is what the shape falls back on, and nothing at
			// all where it brought none.
			fill := clonePath(n.Path)
			fill.Transform(m)
			if shade, ok := patternShade(st.fillPattern, n, m, current, st.Opacity*st.FillOpacity, patterning); ok {
				cv.FillPathFunc(fill, st.FillRule, shade)
			} else if c, ok := st.fillFallbackColour(current); ok {
				if c := fade(c, st.Opacity*st.FillOpacity); c.A() > 0 {
					cv.FillPath(fill, c, st.FillRule)
				}
			}
		default:
			c := st.Fill
			if st.FillCurrent {
				c = current
			}
			if c := fade(c, st.Opacity*st.FillOpacity); c.A() > 0 {
				// The path is drawn in its own coordinates and the transform is
				// put on a copy of it, so the node keeps the shape the drawing
				// wrote while the picture comes out where the transform put it.
				// The stroke takes its width from the same transform rather
				// than from the geometry — see [Style.strokeWidth].
				fill := clonePath(n.Path)
				fill.Transform(m)
				cv.FillPath(fill, c, st.FillRule)
			}
		}
	}
	if st.HasStroke && st.Width > 0 {
		paintStroke(cv, n, m, current, patterning)
	}
	// The marker comes after the fill and the stroke, which is the order the
	// spec puts them in. It is drawn where the shape has no stroke at all, and
	// where it encloses no fill either: a marker is a picture put on a vertex,
	// and the line beside it is not what puts it there — the vertex is.
	if len(n.markers) > 0 {
		paintMarkers(cv, n, m, current, masking, patterning)
	}
}

// paintStroke widens the outline of one node's shape and paints what comes of
// it. It is the whole of the stroke rather than a part of the shape, so that
// the marker can be drawn after all of it rather than after only some.
func paintStroke(cv *canvas.Canvas, n *Node, m canvas.Matrix, current canvas.Color, patterning []*pattern) {
	st := n.Style
	// A dashed stroke is cut into pieces along the line before it is widened,
	// so that each dash is capped at both ends and a dash never starts halfway
	// through a corner. The cut was made where the shape was written and the
	// transform of the element has carried it along with the shape ever since,
	// so the outline taken here is already in the runs the pattern asked for,
	// and taking it to the canvas carries the cut the rest of the way — the
	// viewBox and the size asked for included.
	outline := n.Path
	if n.dashed != nil {
		outline = n.dashed
	}
	stroke := clonePath(outline)
	stroke.Transform(m)
	style := canvas.StrokeStyle{
		Width:      st.strokeWidth(m),
		Cap:        st.Cap,
		Join:       st.Join,
		MiterLimit: st.MiterLimit,
	}
	if st.strokeGradient != nil {
		// The stroke is widened into the shape it covers, and that shape is
		// filled with the gradient, so the colour runs across the stroke the same
		// way it runs across the fill rather than being one colour of it.
		outline := canvas.StrokeOutline(stroke, style)
		cv.FillPathFunc(outline, canvas.FillNonZero, gradientShade(st.strokeGradient, n, m, current, st.Opacity*st.StrokeOpacity))
		return
	}
	if st.strokePattern != nil {
		// The stroke is widened into the shape it covers, and that shape is
		// painted with the pattern, so the picture runs across the stroke the
		// same way it runs across the fill — with the colour after the
		// `url(...)` where the pattern has nothing to paint with.
		outline := canvas.StrokeOutline(stroke, style)
		if shade, ok := patternShade(st.strokePattern, n, m, current, st.Opacity*st.StrokeOpacity, patterning); ok {
			cv.FillPathFunc(outline, canvas.FillNonZero, shade)
		} else if c, ok := st.strokeFallbackColour(current); ok {
			if c := fade(c, st.Opacity*st.StrokeOpacity); c.A() > 0 {
				cv.FillPath(outline, c, canvas.FillNonZero)
			}
		}
		return
	}
	c := st.Stroke
	if st.StrokeCurrent {
		c = current
	}
	c = fade(c, st.Opacity*st.StrokeOpacity)
	if c.A() == 0 {
		return
	}
	cv.StrokePath(stroke, c, style)
}

// gradientShade is the colour one pixel of a gradient-painted shape takes: the
// middle of the pixel is taken back to the coordinates the shape was written in,
// and the gradient is asked what colour it is there.
//
// The way back is through the transform that moves the shape, which is already
// baked into the path and not the one the gradient knows about — a gradient in
// fractions of the shape's own box follows the shape through it, and one in the
// drawing's own units follows it too, because both are written against where the
// shape is before it is moved. The opacity comes in here rather than into the
// coverage, so a gradient at half strength is the same colours half as strong.
func gradientShade(g *gradient, n *Node, m canvas.Matrix, current canvas.Color, opacity float64) func(x, y int) canvas.Color {
	s := newSampler(g, n.box[0], n.box[1], n.box[2], n.box[3])
	inv, ok := m.Mul(n.Style.Transform).Inverse()
	if !ok {
		// A transform with no way back has no way to say which point of the
		// gradient a pixel is over, so the whole shape takes one colour of it —
		// the one at the corner the shape starts from — rather than none at all.
		c := fade(s.colorAt(n.box[0], n.box[1], current), opacity)
		return func(x, y int) canvas.Color { return c }
	}
	return func(x, y int) canvas.Color {
		gx, gy := inv.Map(float64(x)+0.5, float64(y)+0.5)
		return fade(s.colorAt(gx, gy, current), opacity)
	}
}

// clonePath is a copy of a path, because the one on the node is shared by every
// frame that paints this drawing and must not be transformed in place.
func clonePath(p *canvas.Path) *canvas.Path {
	pts, closed := p.Points()
	out := canvas.NewPath()
	for i, sub := range pts {
		if len(sub) == 0 {
			continue
		}
		out.MoveTo(sub[0].X, sub[0].Y)
		for _, q := range sub[1:] {
			out.LineTo(q.X, q.Y)
		}
		if closed[i] {
			out.Close()
		}
	}
	return out
}

// fade is a colour at an opacity, as a multiply of the two the node was given.
func fade(c canvas.Color, opacity float64) canvas.Color {
	if opacity >= 1 {
		return c
	}
	a := int(float64(c.A())*clampFloat(opacity, 0, 1) + 0.5)
	return canvas.RGBA(c.R(), c.G(), c.B(), uint8(a))
}

// fillFallbackColour and strokeFallbackColour are the colour after a `url(...)`
// where the reference is not one the drawing can follow — never there, or there
// and with nothing to paint with — which is the second choice the file gave, if
// it gave one. Writing it `currentColor` keeps the keyword: the colour of the
// widget is only known when it is painted.
func (st Style) fillFallbackColour(current canvas.Color) (canvas.Color, bool) {
	if !st.hasFillFallback {
		return 0, false
	}
	if st.fillFallbackCurrent {
		return current, true
	}
	return st.fillFallback, true
}

func (st Style) strokeFallbackColour(current canvas.Color) (canvas.Color, bool) {
	if !st.hasStrokeFallback {
		return 0, false
	}
	if st.strokeFallbackCurrent {
		return current, true
	}
	return st.strokeFallback, true
}

// scaleOf is how much a transform makes things bigger, taken as the square root
// of how much area it covers. A drawing scaled by two has its stroke twice as
// wide, and one drawn through a matrix that only rotates is left alone, which
// is the average of the two numbers on the diagonal.
func scaleOf(m canvas.Matrix) float64 {
	det := math.Abs(m.A*m.D - m.B*m.C)
	return math.Sqrt(det)
}

// strokeWidth is how wide the stroke comes out on the canvas: the width the
// drawing wrote, taken through everything that scaled the shape — the viewBox
// onto the canvas and the transforms the shape went through, which are the same
// ones that moved the geometry, so that a shape scaled twice is a stroke twice
// as wide. That is the default SVG has.
//
// `vector-effect="non-scaling-stroke"` asks for the other one: the width it
// says, in the pixels the drawing comes out at, with nothing scaled at all —
// neither the viewBox nor any transform, which is what keeps a hairline a
// hairline when the drawing is zoomed. The canvas is painted supersample times
// the size asked for and brought back down at the end, so the width is in those
// bigger pixels and comes down with the rest of the picture.
func (st Style) strokeWidth(m canvas.Matrix) float64 {
	if st.NonScalingStroke {
		return st.Width * supersample
	}
	// The transform only counts where it was put into the path, under the same
	// condition the build used to put it there: a style with no transform of
	// its own scales nothing, and a zero matrix would scale the width to none.
	if st.HasTransform && st.Transform != (canvas.Matrix{}) {
		m = m.Mul(st.Transform)
	}
	return st.Width * scaleOf(m)
}

// fitTransform is the transform that maps a viewBox onto a rectangle of the
// given size, keeping the two in the proportion the drawing was made in. A
// drawing whose viewBox is not the same shape as the canvas is fitted inside it
// and centred, the way a picture is, rather than stretched to fill it. The size
// is any two numbers in the coordinates the drawing is being put into: pixels
// when it is painted, the units of the drawing when a `<symbol>` is fitted into
// the box a `<use>` gave it.
func fitTransform(view [4]float64, width, height float64) canvas.Matrix {
	vw, vh := view[2], view[3]
	if vw <= 0 || vh <= 0 {
		return canvas.Scale(width, height)
	}
	scale := math.Min(width/vw, height/vh)
	sx, sy := scale, scale
	// Keep the drawing's own proportions inside the canvas.
	if vw*scale > width {
		sx = width / vw
	}
	if vh*scale > height {
		sy = height / vh
	}
	// The drawing is centred in what is left over, which is nothing at all when
	// the canvas and the viewBox have the same shape, as an icon's usually do.
	ox := (width - vw*sx) / 2
	oy := (height - vh*sy) / 2
	return canvas.Translate(ox-view[0]*sx, oy-view[1]*sy).Mul(canvas.Scale(sx, sy))
}
