//go:build android

package app

func awaitPermission(code int32) chan PermissionResult {
	ch := make(chan PermissionResult, 1)
	waitMu.Lock()
	permWait[code] = ch
	waitMu.Unlock()
	return ch
}
