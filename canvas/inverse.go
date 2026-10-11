package canvas

import "math"

// Inverse returns the transform that undoes m. A degenerate matrix — a zero
// determinant, as in scale(0) or a collapsed skew — has no inverse, and the
// bool reports it so a caller can paint nothing instead of dividing by zero.
func (m Matrix) Inverse() (Matrix, bool) {
	det := m.A*m.D - m.B*m.C
	if math.Abs(det) < 1e-12 {
		return Matrix{}, false
	}
	return Matrix{
		A: m.D / det,
		B: -m.B / det,
		C: -m.C / det,
		D: m.A / det,
		E: (m.C*m.F - m.D*m.E) / det,
		F: (m.B*m.E - m.A*m.F) / det,
	}, true
}
