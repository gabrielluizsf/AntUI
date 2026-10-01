package svg

import (
	"math"
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
		name := strings.ToLower(strings.TrimSpace(s[:open]))
		args := splitArgs(s[open+1 : open+close])
		switch name {
		case "matrix":
			// The whole transform named at once: a, b, c, d, e and f, which are
			// the six numbers every other form is worked out to.
			if len(args) == 6 {
				var v [6]float64
				read := true
				for i, a := range args {
					n, err := parseNumber(a)
					if err != nil {
						read = false
						break
					}
					v[i] = n
				}
				if read {
					m = m.Mul(canvas.Matrix{A: v[0], B: v[1], C: v[2], D: v[3], E: v[4], F: v[5]})
					found = true
				}
			}
		case "translate":
			if len(args) == 1 || len(args) == 2 {
				tx, _ := parseNumber(args[0])
				ty := 0.0
				if len(args) == 2 {
					ty, _ = parseNumber(args[1])
				}
				m = m.Mul(canvas.Translate(tx, ty))
				found = true
			}
		case "scale":
			if len(args) == 1 || len(args) == 2 {
				sx, err1 := parseNumber(args[0])
				sy, err2 := sx, err1
				if len(args) == 2 {
					sy, err2 = parseNumber(args[1])
				}
				if err1 == nil && err2 == nil {
					m = m.Mul(canvas.Scale(sx, sy))
					found = true
				}
			}
		case "rotate":
			if len(args) == 1 || len(args) == 3 {
				deg, err := parseAngle(args[0])
				if err != nil {
					break
				}
				r := canvas.Rotate(deg * math.Pi / 180)
				ok := true
				if len(args) == 3 {
					// Rotating about a point: move there, spin, move back.
					cx, err1 := parseNumber(args[1])
					cy, err2 := parseNumber(args[2])
					if err1 == nil && err2 == nil {
						r = canvas.Translate(cx, cy).Mul(r).Mul(canvas.Translate(-cx, -cy))
					} else {
						ok = false
					}
				}
				if ok {
					m = m.Mul(r)
					found = true
				}
			}
		case "skewx":
			if len(args) == 1 {
				if deg, err := parseAngle(args[0]); err == nil {
					m = m.Mul(canvas.Skew(deg*math.Pi/180, 0))
					found = true
				}
			}
		case "skewy":
			if len(args) == 1 {
				if deg, err := parseAngle(args[0]); err == nil {
					m = m.Mul(canvas.Skew(0, deg*math.Pi/180))
					found = true
				}
			}
		}
		s = strings.TrimSpace(s[open+close+1:])
		// A comma or whitespace may sit between two operations.
		if s != "" && (s[0] == ',' || s[0] == ' ') {
			s = s[1:]
		}
	}
	return m, found
}

// splitArgs breaks the inside of a transform function into its numbers, which
// may be separated by commas, by spaces, or by both.
func splitArgs(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t' || r == '\n' || r == '\r'
	})
}
