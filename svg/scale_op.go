package svg

import "github.com/gabrielluizsf/antui/canvas"

// scaleOp is `scale(sx [sy])`: sy defaults to sx.
func scaleOp(args []string) (canvas.Matrix, bool) {
	if len(args) != 1 && len(args) != 2 {
		return canvas.Matrix{}, false
	}
	sx, err1 := parseNumber(args[0])
	sy, err2 := sx, err1
	if len(args) == 2 {
		sy, err2 = parseNumber(args[1])
	}
	if err1 != nil || err2 != nil {
		return canvas.Matrix{}, false
	}
	return canvas.Scale(sx, sy), true
}
