package css

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func parseLenPrefix(s string, u unit, suffix string) (Length, error) {
	body := strings.TrimSuffix(s, suffix)
	v, err := strconv.ParseFloat(strings.TrimSpace(body), 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return Length{}, fmt.Errorf("css: %q is not a length", s)
	}
	return Length{u: u, value: v}, nil
}
