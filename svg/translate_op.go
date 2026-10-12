package svg

import "github.com/gabrielluizsf/antui/canvas"

// translateOp is `translate(tx [ty])`: ty defaults to zero.
func translateOp(args []string) (canvas.Matrix, bool) {
	if len(args) != 1 && len(args) != 2 {
		return canvas.Matrix{}, false
	}
	tx, _ := parseNumber(args[0])
	ty := 0.0
	if len(args) == 2 {
		ty, _ = parseNumber(args[1])
	}
	return canvas.Translate(tx, ty), true
}
