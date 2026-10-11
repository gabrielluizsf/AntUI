package svg

import "github.com/gabrielluizsf/antui/canvas"

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
