package template

import "time"

// setTimeout takes a new timeout, which is both how long the app may take
// and the span the bar fills over. A duration of zero or less would divide
// the bar by nothing and make the quit rule meaningless, so it is left out
// — an argument this package cannot use keeps the one already there, the
// way an unreadable attribute keeps the value it found.
func (l *loading) setTimeout(d time.Duration) {
	if l != nil && d > 0 {
		l.timeout.Store(int64(d))
	}
}
