//go:build windows

package windows

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// SetIcon gives the window a big icon and a small one. Windows asks for both:
// the big one is the task bar and Alt-Tab, the small one the title bar.
func (d *Driver) SetIcon(images []*canvas.Canvas) bool {
	if d.hwnd == 0 || len(images) == 0 {
		return false
	}
	// Largest first is the order they arrive in; the small icon wants
	// something near 16 and the big one something near 32.
	big := iconFor(images, 32)
	small := iconFor(images, 16)
	if big == 0 && small == 0 {
		return false
	}

	if big != 0 {
		procSendMessageW.Call(d.hwnd, wmSetIcon, iconBig, big)
		if d.iconBig != 0 {
			procDestroyIcon.Call(d.iconBig)
		}
		d.iconBig = big
	}
	if small != 0 {
		procSendMessageW.Call(d.hwnd, wmSetIcon, iconSmall, small)
		if d.iconSmall != 0 {
			procDestroyIcon.Call(d.iconSmall)
		}
		d.iconSmall = small
	}
	return true
}
