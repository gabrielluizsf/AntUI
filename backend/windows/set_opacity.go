//go:build windows

package windows

// SetOpacity fades the whole window so what sits behind it shows through;
// 255 is fully opaque. The window is kept opaque to GDI and made layered
// instead: the compositor paints it at a uniform alpha, so the frame itself
// never loses a bit, only its overall brightness changes.
func (d *Driver) SetOpacity(alpha uint8) bool {
	if d.hwnd == 0 {
		return false
	}
	style, _, _ := procGetWindowLongPtr.Call(d.hwnd, gwlExStyle)
	procSetWindowLongPtr.Call(d.hwnd, gwlExStyle, style|wsExLayered)
	// The bAlpha of SetLayeredWindowAttributes is a byte already, so the
	// value travels as-is. A full fade is asked for by putting opacity back
	// at 255 the same way it was lowered.
	procSetLayeredWindowAttr.Call(d.hwnd, 0, uintptr(alpha), lwaAlpha)
	return true
}
