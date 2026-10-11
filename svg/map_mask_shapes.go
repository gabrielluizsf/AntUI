package svg

import "github.com/gabrielluizsf/antui/canvas"

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
