//go:build windows

package windows

import (
	"fmt"

	"github.com/gabrielluizsf/antui/backend"
	"runtime"
)

// Open creates the window that Options describes.
//
// Windows delivers messages to the thread that created the window, so the
// frame loop has to stay on it. Callers must drive the window from the
// goroutine that opened it.
func Open(win backend.Face, opts backend.Options) (*Driver, error) {
	if win == nil {
		return nil, fmt.Errorf("antui: windows backend with no window to drive")
	}
	runtime.LockOSThread()

	d := &Driver{
		win:            win,
		windowedWidth:  opts.Width,
		windowedHeight: opts.Height,
		themeStale:     true, // nothing has read the theme yet
	}
	d.instance, _, _ = procGetModuleHandleW.Call(0)

	// Keeps the window sharp on scaled displays, with no manifest needed.
	if procSetProcessDPIAware.Find() == nil {
		procSetProcessDPIAware.Call()
	}

	if err := registerClass(d.instance); err != nil {
		runtime.UnlockOSThread()
		return nil, err
	}

	hwnd, err := createWindow(d, opts)
	if err != nil {
		runtime.UnlockOSThread()
		return nil, err
	}
	d.hwnd = hwnd

	windowsMu.Lock()
	windows[hwnd] = d
	windowsMu.Unlock()

	// Say the window takes a drag. Without it the pointer shows the no-entry
	// sign over the window and the message never comes.
	d.acceptDrops()

	procShowWindow.Call(hwnd, swShowNormal)
	procUpdateWindow.Call(hwnd)
	procSetForegroundWindow.Call(hwnd)
	return d, nil
}
