package css

import (
	"fmt"
	"strings"
)

// parseLength reads "12", "12px", "50%", "1.5rem", "3vw", "2in" or "auto".
func parseLength(raw string) (Length, error) {
	s := strings.TrimSpace(strings.ToLower(raw))
	if s == "" {
		return Length{}, fmt.Errorf("css: empty length")
	}
	if s == "auto" {
		return Auto(), nil
	}
	if s == "none" {
		return Length{u: unitNone}, nil
	}
	if s == "inherit" || s == "initial" || s == "unset" {
		return Length{u: unitNone}, nil
	}
	if strings.HasSuffix(s, "%") {
		return parseLenPrefix(s, unitPct, "%")
	}
	for _, suf := range unitSuffixes {
		if strings.HasSuffix(s, suf.name) {
			return parseLenPrefix(s, suf.u, suf.name)
		}
	}
	return parseLenPrefix(s, unitPx, "")
}
