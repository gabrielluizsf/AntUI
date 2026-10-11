package template

import "time"

// getTimeout is the timeout in force.
func (l *loading) getTimeout() time.Duration {
	if l == nil {
		return 0
	}
	return time.Duration(l.timeout.Load())
}
