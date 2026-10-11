package svg

import "github.com/gabrielluizsf/antui/canvas"

// clipShapesOf is what a built `<clipPath>` adds to the clip: the outline of
// every shape inside it, read with the rule it asked for. Everything that
// draws nothing adds nothing — a group is its children, a definition is
// nothing until something points at it, and what is hidden is nowhere — and
// the one thing that has an outline but no shape to cut with says so, because
// a drawing that asked to cut with writing is a drawing that is not going to
// get what it asked for.
//
// A `clip-path` written on the way down to a shape cuts that shape: each one
// found on the way is carried the rest of the way and put on the shape, so
// what comes out is a shape that keeps where it keeps and where every clip
// over it keeps too. The list of them is copied as it is carried down, so
// that two shapes side by side do not append to the same list behind each
// other's back and end up cutting each other's clips away.
func (img *Image) clipShapesOf(n *Node) []canvas.MaskShape {
	var shapes []canvas.MaskShape
	var take func(n *Node, cut [][]canvas.MaskShape)
	take = func(n *Node, cut [][]canvas.MaskShape) {
		if n.clip != nil {
			// The clip is carried even when it holds nothing: a reference the
			// drawing could not follow keeps nothing, so the shape it was
			// written on falls out of the clipPath, the same as a clip of
			// nothing anywhere else in the drawing.
			cut = append(append([][]canvas.MaskShape(nil), cut...), append([]canvas.MaskShape(nil), n.clip.shapes...))
		}
		switch {
		case n.Path != nil && !n.Path.Empty():
			shapes = append(shapes, canvas.MaskShape{Path: n.Path, Rule: n.Style.ClipRule, And: cut})
		case len(n.Runs) > 0:
			img.warnings.warn(n.Name,
				"writing inside a <clipPath> has no outline to cut with here, so it is left out of the clip")
		case n.Pic != nil:
			img.warnings.warn(n.Name,
				"a picture inside a <clipPath> has no outline to cut with here, so it is left out of the clip")
		default:
			for _, k := range n.Kids {
				take(k, cut)
			}
		}
	}
	take(n, nil)
	return shapes
}
