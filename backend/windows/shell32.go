//go:build windows

package windows

import "syscall"

var (
	shell32 = syscall.NewLazyDLL("shell32.dll")

	procDragAcceptFiles = shell32.NewProc("DragAcceptFiles")
	procDragQueryFileW  = shell32.NewProc("DragQueryFileW")
	procDragQueryPoint  = shell32.NewProc("DragQueryPoint")
	procDragFinish      = shell32.NewProc("DragFinish")

	procShellNotifyIconW = shell32.NewProc("Shell_NotifyIconW")
	procLoadIconW        = user32.NewProc("LoadIconW")

	procOpenClipboard    = user32.NewProc("OpenClipboard")
	procCloseClipboard   = user32.NewProc("CloseClipboard")
	procGetClipboardData = user32.NewProc("GetClipboardData")
	procGlobalLock       = kernel.NewProc("GlobalLock")
	procGlobalUnlock     = kernel.NewProc("GlobalUnlock")
)
