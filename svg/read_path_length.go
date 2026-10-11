package svg

import (
	"strconv"
	"strings"
)

// readPathLength reads the total length an element says its path has, in the
// drawing's own units. A drawing writes one so that a stroke can be laid out
// along a path in numbers of its own — ten dashes across a line, whatever the
// line comes out at — rather than in the ones the geometry measures, which are
// only known once the path is built. It is a number and not a length: no unit
// goes on it, and one that is not a positive number is nothing to measure a
// path against, so it is said and left out rather than taken as zero.
func readPathLength(raw string, warn func(string, ...any)) float64 {
	s := strings.TrimSpace(raw)
	if s == "" {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || v <= 0 {
		warn("the path length %q is not a positive number, so the dash pattern is not scaled to it", raw)
		return 0
	}
	return v
}
