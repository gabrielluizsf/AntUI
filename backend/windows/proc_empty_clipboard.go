//go:build windows

package windows

var (
	procEmptyClipboard   = user32.NewProc("EmptyClipboard")
	procSetClipboardData = user32.NewProc("SetClipboardData")
	procGlobalAlloc      = kernel.NewProc("GlobalAlloc")
	procGlobalFree       = kernel.NewProc("GlobalFree")
)
