package css

import (
	"math"
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
