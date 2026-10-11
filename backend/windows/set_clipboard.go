//go:build windows

package windows

import (
	"syscall"
	"unsafe"
)

// SetClipboard puts text on the clipboard as UTF-16, which is what every
// Windows program reads.
func (d *Driver) SetClipboard(text string) bool {
	if ok, _, _ := procOpenClipboard.Call(d.hwnd); ok == 0 {
		return false
	}
	defer procCloseClipboard.Call()
	procEmptyClipboard.Call()

	chars, err := syscall.UTF16FromString(text)
	if err != nil {
		return false
	}
	size := uintptr(len(chars) * 2)
	handle, _, _ := procGlobalAlloc.Call(gmemMoveable, size)
	if handle == 0 {
		return false
	}
	at, _, _ := procGlobalLock.Call(handle)
	if at == 0 {
		procGlobalFree.Call(handle)
		return false
	}
	copy(unsafe.Slice((*uint16)(unsafe.Pointer(at)), len(chars)), chars)
	procGlobalUnlock.Call(handle)

	// Once it is handed over the clipboard owns the block; freeing it here
	// would be freeing somebody else's memory.
	if ok, _, _ := procSetClipboardData.Call(cfUnicodeText, handle); ok == 0 {
		procGlobalFree.Call(handle)
		return false
	}
	return true
}
