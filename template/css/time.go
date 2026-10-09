package css

import (
	"strings"
)

// Time is a duration in milliseconds, the unit the frame loop ticks in. CSS
// spells seconds and milliseconds; both collapse here so transition and
// animation durations compare the same way.
type Time struct {
	ms float64
}

// Sec builds a time from seconds.
func Sec(s float64) Time { return Time{ms: s * 1000} }

// MSec builds a time from milliseconds.
func MSec(ms float64) Time { return Time{ms: ms} }

// MS returns the duration in milliseconds.
func (t Time) MS() float64 { return t.ms }

// Sec returns the duration in seconds.
func (t Time) Sec() float64 { return t.ms / 1000 }

// ParseTime reads a CSS time: 300ms or 1.5s. A unitless number is invalid in
// CSS, so it reports failure.
func ParseTime(raw string) (Time, bool) {
	s := strings.TrimSpace(strings.ToLower(raw))
	if n, ok := numberSuffix(s, "ms"); ok {
		return MSec(n), true
	}
	if n, ok := numberSuffix(s, "s"); ok {
		return Sec(n), true
	}
	return Time{}, false
}
