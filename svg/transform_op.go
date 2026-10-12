package svg

import "github.com/gabrielluizsf/antui/canvas"

// transformOp reads one named operation of a `transform` and answers the
// matrix it multiplies in, and whether its arguments were readable.
func transformOp(name string, args []string) (canvas.Matrix, bool) {
	switch name {
	case "matrix":
		return matrixOp(args)
	case "translate":
		return translateOp(args)
	case "scale":
		return scaleOp(args)
	case "rotate":
		return rotateOp(args)
	case "skewx":
		return skewOp(args, true)
	case "skewy":
		return skewOp(args, false)
	}
	return canvas.Matrix{}, false
}
