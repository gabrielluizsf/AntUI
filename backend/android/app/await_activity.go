//go:build android

package app

func awaitActivity(code int32) chan ActivityResult {
	ch := make(chan ActivityResult, 1)
	waitMu.Lock()
	actWait[code] = ch
	waitMu.Unlock()
	return ch
}
