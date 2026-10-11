//go:build windows

package windows

import (
	"github.com/gabrielluizsf/antui/backend"
	"unsafe"
)

// dropped reads the files out of a WM_DROPFILES and lets the shell go.
func (d *Driver) dropped(hdrop uintptr) {
	if d.win == nil || hdrop == 0 {
		return
	}
	files := filesFromDrop(hdrop)

	// Where the pointer was when it let go, in window coordinates — which is
	// what everything else here counts in.
	var point struct{ X, Y int32 }
	procDragQueryPoint.Call(hdrop, uintptr(unsafe.Pointer(&point)))
	procDragFinish.Call(hdrop)

	if len(files) > 0 {
		d.win.Push(backend.Event{Type: backend.EventDropFiles, Files: files,
			X: int(point.X), Y: int(point.Y)})
	}
}
