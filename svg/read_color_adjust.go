package svg

import "github.com/gabrielluizsf/antui/canvas"

// readColorAdjust reads the functions that take one amount: grayscale, sepia,
// invert, brightness and contrast.
func readColorAdjust(name, args string) (filterOp, bool) {
	v, ok := filterAmount(args)
	if !ok {
		return filterOp{}, false
	}
	kind := canvas.FilterGrayscale
	switch name {
	case "sepia":
		kind = canvas.FilterSepia
	case "invert":
		kind = canvas.FilterInvert
	case "brightness":
		kind = canvas.FilterBrightness
	case "contrast":
		kind = canvas.FilterContrast
	}
	return filterOp{kind: kind, amount: v}, true
}
