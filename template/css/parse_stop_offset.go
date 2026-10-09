package css

import (
	"strings"
)

import (
	"strconv"
)

import (
	"fmt"
)

// parseStopOffset reads one stop position: a percentage or a bare number,
// both in fractions, or a length, which is left to the percentage scale since
// the box is not known yet. Angles and the auto/none keywords do not belong.
func parseStopOffset(s string, ctx Units) (float64, error) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "%") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
		return v / 100, err
	}
	l, err := parseLengthAt(s, ctx)
	if err != nil {
		return 0, err
	}
	if l.IsPct() {
		return l.value / 100, nil
	}
	return 0, fmt.Errorf("css: unsupported gradient stop position %q", s)
}
