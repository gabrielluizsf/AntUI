package svg

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// scalesOf is [scaleOf] taken one axis at a time: the length of the two unit
// vectors after the transform. A pure scale(2, 8) comes out two and eight, where
// the area's square root of four is what a width read off both axes at once
// would give — see [paintStroke] for one that goes further and follows the
// path's own coordinates rather than either axis.
func scalesOf(m canvas.Matrix) scale2 {
	return scale2{x: math.Hypot(m.A, m.C), y: math.Hypot(m.B, m.D)}
}
