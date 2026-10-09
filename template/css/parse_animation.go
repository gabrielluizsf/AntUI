package css

import (
	"strconv"
	"strings"
)

// parseAnimation reads one comma-free animation. Names are kept as the
// author wrote them; the keywords are case-insensitive.
func parseAnimation(part string) (Animation, bool) {
	part = strings.TrimSpace(part)
	if part == "" {
		return Animation{}, false
	}
	a := Animation{Iterations: 1, Direction: AnimNormal, Fill: FillNone}
	var sawDur, sawDelay bool
	for _, tk := range splitTokens(part) {
		t := strings.ToLower(tk)
		if applyAnimKeyword(&a, t) {
			continue
		}
		if tim, ok := parseTiming(t); ok {
			a.Timing = tim
			continue
		}
		if tm, ok := ParseTime(t); ok {
			if !sawDur {
				a.Duration = tm
				sawDur = true
			} else if !sawDelay {
				a.Delay = tm
				sawDelay = true
			} else {
				return Animation{}, false
			}
			continue
		}
		if v, err := strconv.ParseFloat(t, 64); err == nil {
			a.Iterations = v
			continue
		}
		if strings.ContainsRune(t, '(') {
			return Animation{}, false
		}
		if a.Name != "" {
			return Animation{}, false
		}
		a.Name = trimQuotes(tk)
	}
	return a, a.Name != ""
}
