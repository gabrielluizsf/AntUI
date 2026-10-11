package canvas

// Translate shifts by tx to the right and ty down, the canvas's y growing down.
func Translate(tx, ty float64) Matrix { return Matrix{A: 1, D: 1, E: tx, F: ty} }
