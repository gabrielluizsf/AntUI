package svg

import "github.com/gabrielluizsf/antui/canvas"

// fitTransform is the transform that maps a viewBox onto a rectangle of the
// given size, keeping the two in the proportion the drawing was made in — the
// centred meeting a drawing that says nothing about preserveAspectRatio gets,
// rather than the fit and alignment the attribute could ask for — see
// [fitViewBox]. The size is any two numbers in the coordinates the drawing is
// being put into: pixels when it is painted, the units of the drawing when a
// `<symbol>` is fitted into the box a `<use>` gave it.
func fitTransform(view [4]float64, width, height float64) canvas.Matrix {
	return fitViewBox(view, &Picture{
		Box:   [4]float64{0, 0, width, height},
		Fit:   PictureFitMeet,
		Align: AlignXMidYMid,
	})
}
