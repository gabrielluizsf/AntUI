//go:build windows

package windows

// notifyIconDataW is NOTIFYICONDATAW in its Vista-and-later shape, the one
// with the balloon's own icon at the end. The Go field order and natural
// alignment match the C layout on both 32- and 64-bit Windows.
type notifyIconDataW struct {
	cbSize           uint32
	hWnd             uintptr
	uID              uint32
	uFlags           uint32
	uCallbackMessage uint32
	hIcon            uintptr
	szTip            [128]uint16
	dwState          uint32
	dwStateMask      uint32
	szInfo           [256]uint16
	uVersion         uint32 // a union with uTimeout, ignored since Vista
	szInfoTitle      [64]uint16
	dwInfoFlags      uint32
	guidItem         [16]byte
	hBalloonIcon     uintptr
}

// Shell_NotifyIcon's messages and flags.
const (
	nimAdd    = 0x0
	nimModify = 0x1
	nimDelete = 0x2

	nifIcon = 0x02
	nifTip  = 0x04
	nifInfo = 0x10

	niifInfo = 0x1

	// notifyIconID is which of the window's notification icons this is. A
	// window only ever has the one.
	notifyIconID = 1

	idiApplication = 32512
)
