package template

import "time"

// progress is how far along the bar is: the time gone over the time allowed,
// which puts 100% exactly where the timeout lands. It is read from the
// clock rather than ticked, so a frame that took long to draw still reports
// the truth and the bar never lies about how much time is left.
func (l *loading) progress() float32 {
	if l == nil {
		return 0
	}
	return progressAt(time.Since(l.start), time.Duration(l.timeout.Load()))
}
