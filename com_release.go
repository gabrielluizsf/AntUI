//go:build windows

package antui

// release is IUnknown::Release, which every interface has third.
func comRelease(this uintptr) {
	if this != 0 {
		comCall(this, 2)
	}
}
