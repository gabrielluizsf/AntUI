package css

import (
	"strings"
)

// parseTransition reads one comma-free transition: a property name (which may
// be "all" or "none"), two times with the duration first, an easing, in any
// order. A missing property means every property, and a missing timing means
// the default silence of [Linear].
func parseTransition(part string) (Transition, bool) {
	part = strings.TrimSpace(part)
	if part == "" {
		return Transition{}, false
	}
	var tr Transition
	var sawDur, sawDelay bool
	for _, tk := range splitTokens(part) {
		if tim, ok := parseTiming(tk); ok {
			tr.Timing = tim
			continue
		}
		if tm, ok := ParseTime(tk); ok {
			if !sawDur {
				tr.Duration = tm
				sawDur = true
			} else if !sawDelay {
				tr.Delay = tm
				sawDelay = true
			} else {
				return Transition{}, false
			}
			continue
		}
		if tr.Prop == "" {
			if tk == "none" {
				tr.Prop = "none"
			} else if !strings.ContainsRune(tk, '(') {
				tr.Prop = strings.ToLower(tk)
			}
		}
	}
	if tr.Prop == "" {
		tr.Prop = "all"
	}
	return tr, true
}
