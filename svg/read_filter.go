package svg

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// readFilter reads one function of the list, such as `blur(2px)` or
// `grayscale(50%)`. What is not one of the functions this package can apply,
// does not open and close its arguments, or gives a number that is not one
// answers false, and the caller says so rather than guessing what was meant.
func readFilter(tok string) (filterOp, bool) {
	open := strings.IndexByte(tok, '(')
	if open < 0 || !strings.HasSuffix(tok, ")") {
		return filterOp{}, false
	}
	name := strings.ToLower(strings.TrimSpace(tok[:open]))
	args := strings.TrimSpace(tok[open+1 : len(tok)-1])
	switch name {
	case "grayscale", "sepia", "invert", "brightness", "contrast":
		return readColorAdjust(name, args)
	case "hue-rotate":
		deg, err := parseAngle(args)
		if err != nil {
			return filterOp{}, false
		}
		return filterOp{kind: canvas.FilterHueRotate, amount: deg}, true
	case "blur":
		if args == "" {
			// A blur of nothing is no blur, which is what the function with no
			// argument says in CSS, and it is read rather than refused so that
			// a list containing it is not warning about something harmless.
			return filterOp{kind: canvas.FilterBlur, amount: 0}, true
		}
		v, ok := parseLength(args)
		if !ok || v < 0 {
			return filterOp{}, false
		}
		return filterOp{kind: canvas.FilterBlur, amount: v}, true
	case "drop-shadow":
		sh, ok := readDropShadow(args)
		if !ok {
			return filterOp{}, false
		}
		return filterOp{shadow: sh}, true
	}
	return filterOp{}, false
}
