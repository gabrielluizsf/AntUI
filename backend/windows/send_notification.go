//go:build windows

package windows

import (
	"os"
	"path/filepath"
	"unsafe"
)

// SendNotification shows a notification balloon from the window's icon in
// the notification area — on Windows 10 and later the system draws it as a
// toast and keeps it in the notification centre. It reports whether the
// shell took it.
//
// The balloon has to come from a notification-area icon, so the first
// notification adds one and later ones reuse it; [Driver.Close] takes it
// away again, which is what keeps a closed program from leaving a ghost
// icon behind.
func (d *Driver) SendNotification(title, body string) bool {
	if d.hwnd == 0 {
		return false
	}
	data := d.notifyIconData()
	data.uFlags = nifIcon | nifTip | nifInfo
	data.dwInfoFlags = niifInfo
	copyUTF16(data.szTip[:], filepath.Base(os.Args[0]))
	copyUTF16(data.szInfoTitle[:], title)
	copyUTF16(data.szInfo[:], body)

	message := uintptr(nimModify)
	if !d.notifyIcon {
		message = nimAdd
	}
	ok, _, _ := procShellNotifyIconW.Call(message, uintptr(unsafe.Pointer(&data)))
	if ok == 0 {
		return false
	}
	d.notifyIcon = true
	return true
}

// notifyIconData is the part of a NOTIFYICONDATAW every message about the
// window's icon shares: which icon it is, and the picture it shows — the
// window's own small icon when it has been given one.
func (d *Driver) notifyIconData() notifyIconDataW {
	data := notifyIconDataW{hWnd: d.hwnd, uID: notifyIconID, hIcon: d.iconSmall}
	data.cbSize = uint32(unsafe.Sizeof(data))
	if data.hIcon == 0 {
		data.hIcon, _, _ = procLoadIconW.Call(0, idiApplication)
	}
	return data
}
