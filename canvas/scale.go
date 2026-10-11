package canvas

// Scale multiplies sizes by sx along x and sy along y. A negative factor
// mirrors the axis.
func Scale(sx, sy float64) Matrix { return Matrix{A: sx, D: sy} }
