//go:build windows

package windows

// acceptDrops tells the shell this window will take a drag. Without it the
// pointer shows the no-entry sign and WM_DROPFILES never arrives.
func (d *Driver) acceptDrops() {
	if d.hwnd != 0 {
		procDragAcceptFiles.Call(d.hwnd, 1)
	}
}
