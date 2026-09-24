package css

import (
	"math"
	"strconv"
	"strings"
)

// The kinds a [TransformFunc] can be. rotate() and skew() carry their angles
// in the Ax/Ay fields; a matrix() fills M; scale() lands in Sx/Sy and
// translate() keeps Lengths so percentages stay percentages until the box is
// known.
const (
	TransformTranslate uint8 = iota
	TransformScale
	TransformRotate
	TransformSkew
	TransformMatrix
)

// TransformFunc is one function of a transform list. The zero value is a
// translate with no offset, and [Style.Transform] being nil means "none".
type TransformFunc struct {
	Kind uint8
	// translate() offsets; a percentage is of the border box.
	Dx, Dy Length
	// scale() factors.
	Sx, Sy float64
	// rotate() uses Ax; skew() uses both.
	Ax, Ay Angle
	// matrix() is the six numbers of the CSS/SVG matrix(a, b, c, d, e, f).
	M [6]float64
}

// parseTransform reads a space-separated transform list: functions such as
// translate(10px 20%), rotate(45deg), scale(2), skew(20deg) or matrix(...).
// "none" clears the list, and any function the engine does not know drops the
// whole property exactly as an unknown value would.
func parseTransform(raw string, ctx Units) ([]TransformFunc, bool) {
	s := strings.TrimSpace(raw)
	if s == "none" {
		return nil, true
	}
	var out []TransformFunc
	for _, tk := range splitTokens(s) {
		open := strings.IndexByte(tk, '(')
		if open < 0 || !strings.HasSuffix(tk, ")") {
			return nil, false
		}
		name := strings.ToLower(strings.TrimSpace(tk[:open]))
		inner := strings.TrimSpace(tk[open+1 : len(tk)-1])
		f, ok := parseTransformFunc(name, inner, ctx)
		if !ok {
			return nil, false
		}
		out = append(out, f)
	}
	if len(out) == 0 {
		return nil, false
	}
	return out, true
}

// parseTransformFunc reads one function's arguments. translateX/Y and skewX/Y
// fold into their two-axis siblings, with the missing axis at its zero.
func parseTransformFunc(name, inner string, ctx Units) (TransformFunc, bool) {
	args, ok := transformArgs(inner)
	if !ok {
		return TransformFunc{}, false
	}
	switch name {
	case "translate":
		if len(args) != 1 && len(args) != 2 {
			return TransformFunc{}, false
		}
		dx, ok := transformLength(args[0], ctx)
		if !ok {
			return TransformFunc{}, false
		}
		dy := Zero()
		if len(args) == 2 {
			if dy, ok = transformLength(args[1], ctx); !ok {
				return TransformFunc{}, false
			}
		}
		return TransformFunc{Kind: TransformTranslate, Dx: dx, Dy: dy}, true
	case "translatex":
		if len(args) != 1 {
			return TransformFunc{}, false
		}
		dx, ok := transformLength(args[0], ctx)
		if !ok {
			return TransformFunc{}, false
		}
		return TransformFunc{Kind: TransformTranslate, Dx: dx}, true
	case "translatey":
		if len(args) != 1 {
			return TransformFunc{}, false
		}
		dy, ok := transformLength(args[0], ctx)
		if !ok {
			return TransformFunc{}, false
		}
		return TransformFunc{Kind: TransformTranslate, Dy: dy}, true
	case "scale":
		if len(args) != 1 && len(args) != 2 {
			return TransformFunc{}, false
		}
		sx, ok := transformFactor(args[0])
		if !ok {
			return TransformFunc{}, false
		}
		sy := sx
		if len(args) == 2 {
			if sy, ok = transformFactor(args[1]); !ok {
				return TransformFunc{}, false
			}
		}
		return TransformFunc{Kind: TransformScale, Sx: sx, Sy: sy}, true
	case "scalex":
		if len(args) != 1 {
			return TransformFunc{}, false
		}
		sx, ok := transformFactor(args[0])
		if !ok {
			return TransformFunc{}, false
		}
		return TransformFunc{Kind: TransformScale, Sx: sx, Sy: 1}, true
	case "scaley":
		if len(args) != 1 {
			return TransformFunc{}, false
		}
		sy, ok := transformFactor(args[0])
		if !ok {
			return TransformFunc{}, false
		}
		return TransformFunc{Kind: TransformScale, Sx: 1, Sy: sy}, true
	case "rotate":
		if len(args) != 1 {
			return TransformFunc{}, false
		}
		a, ok := ParseAngle(args[0])
		if !ok {
			return TransformFunc{}, false
		}
		return TransformFunc{Kind: TransformRotate, Ax: a}, true
	case "skew":
		if len(args) != 1 && len(args) != 2 {
			return TransformFunc{}, false
		}
		ax, ok := ParseAngle(args[0])
		if !ok {
			return TransformFunc{}, false
		}
		ay, _ := ParseAngle(args[1])
		return TransformFunc{Kind: TransformSkew, Ax: ax, Ay: ay}, true
	case "skewx":
		if len(args) != 1 {
			return TransformFunc{}, false
		}
		ax, ok := ParseAngle(args[0])
		if !ok {
			return TransformFunc{}, false
		}
		return TransformFunc{Kind: TransformSkew, Ax: ax}, true
	case "skewy":
		if len(args) != 1 {
			return TransformFunc{}, false
		}
		ay, ok := ParseAngle(args[0])
		if !ok {
			return TransformFunc{}, false
		}
		return TransformFunc{Kind: TransformSkew, Ay: ay}, true
	case "matrix":
		if len(args) != 6 {
			return TransformFunc{}, false
		}
		var m [6]float64
		for i, a := range args {
			v, ok := transformNumber(a)
			if !ok {
				return TransformFunc{}, false
			}
			m[i] = v
		}
		return TransformFunc{Kind: TransformMatrix, M: m}, true
	}
	return TransformFunc{}, false
}

// transformArgs splits a function's arguments on commas and whitespace and
// rejects leftovers that are not real arguments.
func transformArgs(inner string) ([]string, bool) {
	if strings.HasPrefix(inner, ",") || strings.HasSuffix(inner, ",") || strings.Contains(inner, ",,") {
		return nil, false
	}
	var out []string
	for _, f := range strings.FieldsFunc(inner, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t' || r == '\n' || r == '\r'
	}) {
		if strings.ContainsAny(f, "()") {
			return nil, false
		}
		out = append(out, f)
	}
	return out, true
}

// transformLength is a length a transform accepts: a calc(), a percentage or
// a fixed unit surface, but never the auto/none keywords.
func transformLength(s string, ctx Units) (Length, bool) {
	l, err := parseLengthAt(s, ctx)
	if err != nil || l.Auto() || l.None() {
		return Length{}, false
	}
	return l, true
}

// transformNumber is one of matrix()'s six numbers.
func transformNumber(s string) (float64, bool) {
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, false
	}
	return v, true
}

// transformFactor is a scale argument: a unitless number or a percentage,
// which divides by a hundred the way scale(150%) means scale(1.5).
func transformFactor(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "%") {
		v, ok := transformNumber(strings.TrimSuffix(s, "%"))
		return v / 100, ok
	}
	return transformNumber(s)
}

// TransformOrigin puts the pivot of a transform in the border box. The two
// lengths are the x and y of the pivot, a percentage of the box; the initial
// value is "50% 50%", the centre.
type TransformOrigin = [2]Length

// TransformOriginPlain is the CSS spelling the property stores under.
const TransformOriginPlain = "transform-origin"

// InitialTransformOrigin is the CSS initial value: the centre of the box.
var InitialTransformOrigin = [2]Length{Pct(50), Pct(50)}

// parseTransformOrigin reads 1–3 values: a vertical keyword alone sets x to
// the centre; two values order left/right/centre for x and top/centre/bottom
// for y in either spelling; a third, the z, is validated and dropped.
func parseTransformOrigin(raw string, ctx Units) ([2]Length, bool) {
	parts := strings.Fields(strings.ToLower(raw))
	if len(parts) == 0 || len(parts) > 3 {
		return [2]Length{}, false
	}
	if len(parts) == 1 {
		if y, ok := originKeywordY(parts[0]); ok {
			return [2]Length{Pct(50), y}, true
		}
		x, ok := originKeywordX(parts[0], ctx)
		if !ok {
			return [2]Length{}, false
		}
		return [2]Length{x, Pct(50)}, true
	}
	xTok, yTok := parts[0], parts[1]
	if parts[0] == "top" || parts[0] == "bottom" {
		xTok, yTok = parts[1], parts[0]
	}
	x, ok := originKeywordX(xTok, ctx)
	if !ok {
		return [2]Length{}, false
	}
	y, ok := originKeywordYWithX(yTok, ctx)
	if !ok {
		return [2]Length{}, false
	}
	if len(parts) == 3 {
		if _, err := parseLengthAt(parts[2], ctx); err != nil {
			if _, ok := transformNumber(parts[2]); !ok {
				return [2]Length{}, false
			}
		}
	}
	return [2]Length{x, y}, true
}

// originKeywordX reads the x of an origin: a horizontal keyword or a length.
func originKeywordX(s string, ctx Units) (Length, bool) {
	switch s {
	case "left":
		return Pct(0), true
	case "center":
		return Pct(50), true
	case "right":
		return Pct(100), true
	}
	return originAxisLength(s, ctx)
}

// originKeywordY reads the y of an origin: a vertical keyword or a length.
func originKeywordY(s string) (Length, bool) {
	switch s {
	case "top":
		return Pct(0), true
	case "center":
		return Pct(50), true
	case "bottom":
		return Pct(100), true
	}
	return Length{}, false
}

// originKeywordYWithX reads the y of an origin in two-value form, where the
// horizontal keywords are not allowed and a length is.
func originKeywordYWithX(s string, ctx Units) (Length, bool) {
	switch s {
	case "top":
		return Pct(0), true
	case "bottom":
		return Pct(100), true
	case "center":
		return Pct(50), true
	}
	return originAxisLength(s, ctx)
}

// originAxisLength is the shared length/percentage tail of an origin axis.
func originAxisLength(s string, ctx Units) (Length, bool) {
	if strings.HasSuffix(s, "%") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
		if err != nil {
			return Length{}, false
		}
		return Pct(v), true
	}
	if s == "left" || s == "right" || s == "top" || s == "bottom" {
		return Length{}, false
	}
	l, err := parseLengthAt(s, ctx)
	if err != nil || l.Auto() || l.None() {
		return Length{}, false
	}
	return l, true
}
