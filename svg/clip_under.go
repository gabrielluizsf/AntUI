package svg

import "github.com/gabrielluizsf/antui/canvas"

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
