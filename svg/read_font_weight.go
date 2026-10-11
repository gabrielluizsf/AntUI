package svg

import (
	"strconv"
	"strings"
)

// readFontWeight reads how heavy the writing is, out of the weight it already
// has — because `bolder` and `lighter` are not weights of their own, they are
// the step CSS takes from the one in hand. The steps are not evenly spaced: a
// family is not cut evenly, and a weight past either end of the scale is the
// end of it. Anything that is not a weight at all leaves the one in hand alone,
// since a weight this cannot read is not one to draw at.
func readFontWeight(was uint16, raw string, warn func(string, ...any)) uint16 {
	if was == 0 {
		was = FontWeightNormal
	}
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "normal":
		return FontWeightNormal
	case "bold":
		return FontWeightBold
	case "bolder":
		return bolderWeight[was/100]
	case "lighter":
		return lighterWeight[was/100]
	}
	if v, err := strconv.Atoi(strings.TrimSpace(raw)); err == nil && v >= 100 && v <= 900 {
		// A weight on the scale is a weight on the hundreds the words name,
		// so one written between two of them is the nearer of the two.
		return uint16(min((v+50)/100, 9) * 100)
	}
	if strings.TrimSpace(raw) != "" {
		warn("the font weight %q is not a weight, so the writing keeps the weight it had", raw)
	}
	return was
}
