package svg

import "github.com/gabrielluizsf/antui/canvas"

// baseStyle is where the style of everything starts: the fill a drawing that
// says nothing about colour paints with, a stroke the width the spec gives it,
// the opacities that let anything be seen through, and the transform nothing
// has moved from. A `<clipPath>` is built from this too rather than from the
// style where it stands, because its shapes are written in the coordinates of
// whoever points at it and take only what was written on the way down into
// them.
func baseStyle() Style {
	return Style{
		Fill:        canvas.RGB(0, 0, 0),
		HasFill:     true,
		FillRule:    canvas.FillNonZero,
		Width:       1,
		Cap:         canvas.CapButt,
		Join:        canvas.JoinMiter,
		MiterLimit:  canvas.DefaultMiterLimit,
		Opacity:     1,
		FillOpacity: 1, StrokeOpacity: 1,
		Transform: canvas.Identity(), HasTransform: true,
	}
}
