//go:build windows

package windows

// ContentScale is pixels per point, or 0 when the system does not say.
//
// GetDpiForWindow is the per-monitor answer and is right on a machine with
// two displays at different scales, which is now ordinary. It arrived in
// Windows 10 1607, so a machine older than that falls back to the one number
// the whole desktop shares.
func (d *Driver) ContentScale() float64 {
	if d.hwnd != 0 && procGetDpiForWindow.Find() == nil {
		if dpi, _, _ := procGetDpiForWindow.Call(d.hwnd); dpi > 0 {
			return float64(dpi) / defaultWindowsDPI
		}
	}
	hdc, _, _ := procGetDC.Call(0)
	if hdc == 0 {
		return 0
	}
	defer procReleaseDC.Call(0, hdc)
	dpi, _, _ := procGetDeviceCaps.Call(hdc, logPixelsX)
	if dpi == 0 {
		return 0
	}
	return float64(dpi) / defaultWindowsDPI
}
