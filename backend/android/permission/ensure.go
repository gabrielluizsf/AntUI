//go:build android

package permission

// Ensure asks only for what is not already held, and reports whether
// everything asked for is held afterwards. It is what most callers want:
// asking again for something already granted shows no dialog, but it does
// pause the app for a moment, and doing it every launch is visible.
//
// **Not from the frame loop**, for [Request]'s reason — and the trap here is
// worse than Request's, because Ensure only blocks when something is
// actually missing. That is the *first* launch and no other. An app that
// calls this from its frame loop works perfectly every time its author runs
// it, and hangs on the one launch that matters: the user's first.
//
// Run it on a goroutine and let the loop keep drawing:
//
//	go func() { granted <- ensure() }()
func Ensure(names ...string) (bool, error) {
	var missing []string
	for _, n := range names {
		held, err := Held(n)
		if err != nil {
			return false, err
		}
		if !held {
			missing = append(missing, n)
		}
	}
	if len(missing) == 0 {
		return true, nil
	}
	r, err := Request(missing...)
	if err != nil {
		return false, err
	}
	return r.All(), nil
}
