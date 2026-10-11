//go:build windows

package windows

func defWindowProc(hwnd, message, wparam, lparam uintptr) uintptr {
	ret, _, _ := procDefWindowProcW.Call(hwnd, message, wparam, lparam)
	return ret
}
