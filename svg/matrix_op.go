package svg

import "github.com/gabrielluizsf/antui/canvas"

// matrixOp is the `matrix(a b c d e f)` form: all six numbers at once, which
// are the six every other form is worked out to.
func matrixOp(args []string) (canvas.Matrix, bool) {
	if len(args) != 6 {
		return canvas.Matrix{}, false
	}
	var v [6]float64
	for i, a := range args {
		n, err := parseNumber(a)
		if err != nil {
			return canvas.Matrix{}, false
		}
		v[i] = n
	}
	return canvas.Matrix{A: v[0], B: v[1], C: v[2], D: v[3], E: v[4], F: v[5]}, true
}
