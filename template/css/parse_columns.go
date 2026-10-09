package css

import (
	"strconv"
	"strings"
)

// parseColumns reads the columns shorthand, a count and a width in either
// order, which is how a stylesheet asks for "three columns, each at least this
// wide" in a single declaration. A lone number is a count, a lone length a
// width, and auto on its own changes neither.
func parseColumns(raw string, ctx Units) (count int, width Length, ok bool) {
	width = Auto()
	words := splitWords(raw)
	if len(words) == 0 {
		return 0, width, false
	}
	for _, wrd := range words {
		if strings.ToLower(wrd) == "auto" {
			continue
		}
		if n, err := strconv.Atoi(wrd); err == nil {
			if n < 0 {
				return 0, width, false
			}
			count = n
			continue
		}
		l, err := parseLengthAt(wrd, ctx)
		if err != nil || l.Auto() || l.None() || l.value < 0 {
			return 0, width, false
		}
		width = l
	}
	return count, width, true
}
