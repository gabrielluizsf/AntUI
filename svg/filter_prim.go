package svg

import "github.com/gabrielluizsf/antui/canvas"

// filterPrim is one primitive of a `<filter>`: what it takes, what it leaves
// its result under, and the handful of numbers its own kind needs. The inputs
// are the names as written; what they name is resolved while the picture is
// being made, where the pictures are to hand.
type filterPrim struct {
	kind            primKind
	in, in2, result string
	stdDev          float64      // feGaussianBlur
	stdDevY         float64      // feGaussianBlur's second number, if it gave one
	dx, dy          float64      // feOffset, in the drawing's own units
	flood           canvas.Color // feFlood, its colour and opacity apart
	floodOpacity    float64      // feFlood, 1 when none was written
	floodCurrent    bool         // feFlood's colour was currentColor
	matrix          [20]float64  // feColorMatrix
	comp            compositeOp  // feComposite
	blend           blendMode    // feBlend
	k               [4]float64   // feComposite's arithmetic
}
