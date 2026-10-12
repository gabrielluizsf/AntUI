package svg

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// rotateOp is `rotate(deg [cx cy])`: with a centre it moves there, spins and
// moves back.
func rotateOp(args []string) (canvas.Matrix, bool) {
	if len(args) != 1 && len(args) != 3 {
		return canvas.Matrix{}, false
	}
	deg, err := parseAngle(args[0])
	if err != nil {
		return canvas.Matrix{}, false
	}
	r := canvas.Rotate(deg * math.Pi / 180)
	if len(args) == 3 {
		cx, err1 := parseNumber(args[1])
		cy, err2 := parseNumber(args[2])
		if err1 != nil || err2 != nil {
			return canvas.Matrix{}, false
		}
		r = canvas.Translate(cx, cy).Mul(r).Mul(canvas.Translate(-cx, -cy))
	}
	return r, true
}
