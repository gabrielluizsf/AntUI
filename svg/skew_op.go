package svg

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// skewOp is `skewX(deg)` or `skewY(deg)`, which skew along one axis only.
func skewOp(args []string, x bool) (canvas.Matrix, bool) {
	if len(args) != 1 {
		return canvas.Matrix{}, false
	}
	deg, err := parseAngle(args[0])
	if err != nil {
		return canvas.Matrix{}, false
	}
	if x {
		return canvas.Skew(deg*math.Pi/180, 0), true
	}
	return canvas.Skew(0, deg*math.Pi/180), true
}
