package svg

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// scaleOf is how much a transform makes things bigger, taken as the square root
// of how much area it covers. A drawing scaled by two has its stroke twice as
// wide, and one drawn through a matrix that only rotates is left alone, which
// is the average of the two numbers on the diagonal.
func scaleOf(m canvas.Matrix) float64 {
	det := math.Abs(m.A*m.D - m.B*m.C)
	return math.Sqrt(det)
}
