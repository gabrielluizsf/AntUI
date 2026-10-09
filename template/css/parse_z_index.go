package css

import "strconv"

// parseZIndex accepts an integer (possibly negative) z-index. auto is the
// initial value and resets the layer to 0.
func parseZIndex(raw string) (int, bool) {
	if raw == "auto" {
		return 0, true
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, false
	}
	return v, true
}
