package css

import (
	"strconv"
	"strings"
)

// parseFontWeight reads a font-weight: a number rounded to the nearest
// hundred, or normal/bold/bolder/lighter. The relative keywords cannot see the
// inherited weight in a single declaration, so bolder resolves to bold and
// lighter to 300, which is the common approximation.
func parseFontWeight(raw string) (uint16, bool) {
	s := strings.TrimSpace(raw)
	switch s {
	case "normal":
		return FontWeightNormal, true
	case "bold", "bolder":
		return FontWeightBold, true
	case "lighter":
		return 300, true
	}
	v, err := strconv.Atoi(s)
	if err != nil || v < 1 || v > 1000 {
		return 0, false
	}
	return uint16((v + 50) / 100 * 100), true
}
