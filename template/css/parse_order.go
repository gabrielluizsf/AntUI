package css

import (
	"strconv"
)

// parseOrder reads an integer order, auto resets it to zero.
func parseOrder(raw string) (int, bool) {
	if raw == "auto" {
		return 0, true
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, false
	}
	return v, true
}
