package css

// A Length is a size a template can draw with: fixed reference pixels, a
// percentage of the surrounding measure, a length in any CSS unit, or a
// keyword such as auto.
type Length struct {
	u     unit
	value float64
}

// Fixed builds a pixel length.
func Fixed(px float64) Length { return Length{u: unitPx, value: px} }

// Pct builds a percentage length.
func Pct(p float64) Length { return Length{u: unitPct, value: p} }

// Auto is the "auto" keyword, used by width, height and margins.
func Auto() Length { return Length{u: unitAuto} }

// Zero is a zero pixel length.
func Zero() Length { return Fixed(0) }

// IsPct reports whether the length is a percentage.
func (l Length) IsPct() bool { return l.u == unitPct }

// Auto reports whether the length is the auto keyword.
func (l Length) Auto() bool { return l.u == unitAuto }

// None reports whether the length represents the none keyword.
func (l Length) None() bool { return l.u == unitNone }

// Unit returns the length's unit name as CSS writes it, or "" for auto/none.
func (l Length) Unit() string { return unitName[l.u] }
