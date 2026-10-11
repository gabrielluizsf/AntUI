package template

import "time"

// newLoading starts a screen at the default: on, counting from now, with
// the default timeout to close a program that never says Ready.
func newLoading() *loading {
	l := &loading{start: time.Now()}
	l.timeout.Store(int64(defaultLoadingTimeout))
	l.disabled.Store(!loadingEnabledByDefault)
	return l
}
