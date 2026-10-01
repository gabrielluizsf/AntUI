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
	paintNodes(big, img.Root, m, current)

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
func paintNodes(cv *canvas.Canvas, n *Node, m canvas.Matrix, current canvas.Color) {
	if n.Style.Hidden {
		return
	}
	if n.clip != nil {
		paintClipped(cv, n, m, current)
		return
	}
	paintBody(cv, n, m, current)
}

// paintBody is one node and everything inside it, drawn straight onto the
// canvas it is given: its shape first, then its writing, then its children in
// the order they were written.
func paintBody(cv *canvas.Canvas, n *Node, m canvas.Matrix, current canvas.Color) {
	if n.Path != nil && !n.Path.Empty() {
		paintShape(cv, n, m, current)
	}
	if len(n.Runs) > 0 {
		paintText(cv, n, m, current)
	}
	for _, k := range n.Kids {
		paintNodes(cv, k, m, current)
	}
}

// paintClipped draws a node into a picture of its own, cuts that picture down
// to the shapes its clip-path named, and lays the result over what is already
// there. The picture has to be made before it is cut, whole: a group is one
// picture with one clip over it, and cutting each of its children as it was
// drawn would cut them apart from each other, so that two translucent shapes
// overlapping under the edge of the clip would meet somewhere the clip does
// not keep at all.
//
// The picture is the size of the whole canvas rather than of the node, because
// a clip may reach anywhere and the cut has to be of the picture where the clip
// says and nowhere else. What a `<clipPath>` keeps nothing of comes out of this
// the same way: a picture drawn, cut to nothing, and laid over what was there
// as nothing at all.
func paintClipped(cv *canvas.Canvas, n *Node, m canvas.Matrix, current canvas.Color) {
	layer, err := canvas.NewLayer(cv.Width, cv.Height)
	if err != nil {
		// There is no picture to make — a canvas of no size has nothing to
		// draw on and nothing to draw into either.
		return
	}
	paintBody(layer, n, m, current)
	layer.MaskShapes(n.clip.measured(m)...)
	cv.BlitOver(0, 0, layer)
}

// paintShape draws one node's shape: its fill first and its stroke over it,
// which is the order SVG paints, so a stroke half its width over a fill covers
// the fill's edge as it should.
func paintShape(cv *canvas.Canvas, n *Node, m canvas.Matrix, current canvas.Color) {
	st := n.Style
	if st.HasFill {
		if st.fillGradient != nil {
			// A gradient is a colour at every point rather than one colour, so
			// the fill is laid one pixel at a time. The path is still transformed
			// first: the gradient moves with the shape it paints, and what is
			// asked about each pixel is where that pixel falls in the shape's own
			// coordinates rather than where it falls on the canvas.
			fill := clonePath(n.Path)
			fill.Transform(m)
			cv.FillPathFunc(fill, st.FillRule, gradientShade(st.fillGradient, n, m, current, st.Opacity*st.FillOpacity))
		} else {
			c := st.Fill
			if st.FillCurrent {
				c = current
			}
			if c := fade(c, st.Opacity*st.FillOpacity); c.A() > 0 {
				// The path is drawn in its own coordinates and the transform is
				// put on a copy of it, so the fill follows the drawing while the
				// stroke keeps its own width rather than being scaled with it.
				fill := clonePath(n.Path)
				fill.Transform(m)
				cv.FillPath(fill, c, st.FillRule)
			}
		}
	}
	if !st.HasStroke || st.Width <= 0 {
		return
	}
	// A dashed stroke is cut into pieces along the line before it is widened,
	// so that each dash is capped at both ends and a dash never starts halfway
	// through a corner.
	stroke := clonePath(n.Path)
	stroke.Transform(m)
	if st.Dash != nil {
		stroke = stroke.Dashed(*st.Dash)
	}
	style := canvas.StrokeStyle{
		Width:      st.Width * scaleOf(m),
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

// scaleOf is how much a transform makes things bigger, taken as the square root
// of how much area it covers. A drawing scaled by two has its stroke twice as
// wide, and one drawn through a matrix that only rotates is left alone, which
// is the average of the two numbers on the diagonal.
func scaleOf(m canvas.Matrix) float64 {
	det := math.Abs(m.A*m.D - m.B*m.C)
	return math.Sqrt(det)
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
