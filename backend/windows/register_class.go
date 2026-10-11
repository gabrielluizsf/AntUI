//go:build windows

package windows

import (
	"fmt"
	"syscall"
	"unsafe"
)

// registerClass registers the window class once for the whole process. Two
// windows in one program share it, which is what the class is for.
func registerClass(instance uintptr) error {
	classOnce.Do(func() {
		cursor, _, _ := procLoadCursorW.Call(0, 32512) // IDC_ARROW
		wc := wndClassExW{
			Size:      uint32(unsafe.Sizeof(wndClassExW{})),
			Style:     csHRedraw | csVRedraw | csOwnDC,
			WndProc:   wndProcCB,
			Instance:  instance,
			Cursor:    cursor,
			ClassName: utf16Ptr(className),
		}
		atom, _, err := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
		if atom == 0 {
			if errno, ok := err.(syscall.Errno); !ok || uintptr(errno) != errClassAlreadyExists {
				classErr = fmt.Errorf("antui: could not register the window class: %w", err)
			}
		}
	})
	return classErr
}
