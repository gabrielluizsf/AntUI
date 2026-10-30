//go:build windows

package windows

// SetExternal says whether another renderer owns the window, which keeps
// WM_PAINT from blitting the canvas over that renderer's frame. The class
// has no background brush, so there is nothing else to switch off. It
// reports whether there is a window to hand over.
func (d *Driver) SetExternal(on bool) bool {
	d.external = on
	return d.hwnd != 0
}
