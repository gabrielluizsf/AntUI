package template

// boxStart is the first tile position on an axis so the tiling lands on the
// requested edge, wrapping negative positions around the tile cycle.
func boxStart(p, dim int) int {
	if dim <= 0 {
		return 0
	}
	return p - ((p%dim)+dim)%dim
}
