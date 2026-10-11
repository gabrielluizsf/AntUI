//go:build windows

package windows

// awayFrom is how far a size is from the one that is wanted.
func awayFrom(size, want int) int {
	if size < want {
		return want - size
	}
	return size - want
}
