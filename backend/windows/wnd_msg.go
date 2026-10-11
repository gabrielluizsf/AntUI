//go:build windows

package windows

// wndMsg is one window message and the three values Windows delivers with it,
// bundled so a handler takes one argument rather than four.
type wndMsg struct {
	hwnd    uintptr
	message uintptr
	wparam  uintptr
	lparam  uintptr
}
