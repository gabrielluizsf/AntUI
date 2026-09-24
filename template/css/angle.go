package css

import (
	"math"
	"strconv"
	"strings"
)

// Angle is a rotation in radians, the unit a transform matrix needs. Every CSS
// spelling — degrees, radians, gradians, turns — collapses into one float, so
// rotate(30deg) and rotate(0.5236rad) are the same value.
type Angle struct {
	rad float64
}

// Rad builds an angle from radians.
func Rad(rad float64) Angle { return Angle{rad: rad} }

// Deg builds an angle from degrees (360 to a full circle).
func Deg(deg float64) Angle { return Angle{rad: deg * math.Pi / 180} }

// Grad builds an angle from gradians (400 to a full circle).
func Grad(grad float64) Angle { return Angle{rad: grad * math.Pi / 200} }

// Turn builds an angle from turns (1 turn is a full circle).
func Turn(turns float64) Angle { return Angle{rad: turns * 2 * math.Pi} }

// Rad returns the angle in radians.
func (a Angle) Rad() float64 { return a.rad }

// Deg returns the angle in degrees.
func (a Angle) Deg() float64 { return a.rad * 180 / math.Pi }

// ParseAngle reads a CSS angle: 30deg, 1.57rad, 200grad, 0.5turn. The zero
// value is a valid rotation, so the returned bool reports success separately.
func ParseAngle(raw string) (Angle, bool) {
	s := strings.TrimSpace(strings.ToLower(raw))
	vals := []struct {
		suf string
		to  func(float64) Angle
	}{
		{"deg", Deg},
		{"grad", Grad},
		{"rad", Rad},
		{"turn", Turn},
	}
	for _, v := range vals {
		if n, ok := numberSuffix(s, v.suf); ok {
			return v.to(n), true
		}
	}
	return Angle{}, false
}

// Time is a duration in milliseconds, the unit the frame loop ticks in. CSS
// spells seconds and milliseconds; both collapse here so transition and
// animation durations compare the same way.
type Time struct {
	ms float64
}

// Sec builds a time from seconds.
func Sec(s float64) Time { return Time{ms: s * 1000} }

// MSec builds a time from milliseconds.
func MSec(ms float64) Time { return Time{ms: ms} }

// MS returns the duration in milliseconds.
func (t Time) MS() float64 { return t.ms }

// Sec returns the duration in seconds.
func (t Time) Sec() float64 { return t.ms / 1000 }

// ParseTime reads a CSS time: 300ms or 1.5s. A unitless number is invalid in
// CSS, so it reports failure.
func ParseTime(raw string) (Time, bool) {
	s := strings.TrimSpace(strings.ToLower(raw))
	if n, ok := numberSuffix(s, "ms"); ok {
		return MSec(n), true
	}
	if n, ok := numberSuffix(s, "s"); ok {
		return Sec(n), true
	}
	return Time{}, false
}

// numberSuffix parses the float before a unit suffix. It requires the whole
// string to be number+suffix, so a "turn" does not leak into "deg" and "s"
// does not leak out of "ms" (ms is checked first by callers that care).
func numberSuffix(s, suffix string) (float64, bool) {
	if !strings.HasSuffix(s, suffix) {
		return 0, false
	}
	body := strings.TrimSpace(strings.TrimSuffix(s, suffix))
	if body == "" {
		return 0, false
	}
	v, err := strconv.ParseFloat(body, 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, false
	}
	return v, true
}
