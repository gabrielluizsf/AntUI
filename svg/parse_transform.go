package svg

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// parseTransform reads a `transform` attribute, which is a list of the
// operations translate, scale, rotate, skewX and skewY, and answers what they
// come to as one matrix. They are multiplied in the order they are written, so
// that a translate after a scale is a move in the scaled space and not in the
// original one, which is what SVG says and what everybody expects.
func parseTransform(s string) (canvas.Matrix, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return canvas.Matrix{}, false
	}
	m := canvas.Identity()
	found := false
	// Split on the function names, keeping each one with its arguments.
	for len(s) > 0 {
		open := strings.IndexByte(s, '(')
		if open < 0 {
			break
		}
		close := strings.IndexByte(s[open:], ')')
		if close < 0 {
			break
		}
		name := keyword(s[:open])
		args := splitArgs(s[open+1 : open+close])
		if op, ok := transformOp(name, args); ok {
			m = m.Mul(op)
			found = true
		}
		s = strings.TrimSpace(s[open+close+1:])
		// A comma or whitespace may sit between two operations.
		if s != "" && (s[0] == ',' || s[0] == ' ') {
			s = s[1:]
		}
	}
	return m, found
}
