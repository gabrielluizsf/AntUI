package template

// active says whether the screen is up: switched on and not finished. A nil
// screen — a template built as a bare struct rather than through [New] — is
// never active, which keeps every caller below free of a nil check.
func (l *loading) active() bool {
	return l != nil && !l.disabled.Load() && !l.ready.Load()
}
