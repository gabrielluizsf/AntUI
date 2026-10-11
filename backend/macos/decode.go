//go:build darwin

package macos

/*
#cgo CFLAGS: -x objective-c -Wno-deprecated-declarations
#cgo LDFLAGS: -framework Cocoa -framework CoreGraphics -framework CoreVideo
#include <stdlib.h>
#include "backend_darwin.h"
*/
import "C"
import "github.com/gabrielluizsf/antui/backend"

// decode turns one Cocoa-side event into the package's own.
func (d *Driver) decode(win backend.Face, ev *C.antui_d_event) backend.Event {
	out := backend.Event{
		Type:    backend.EventType(ev._type),
		Key:     backend.Key(ev.key),
		Button:  backend.MouseButton(ev.button),
		Mods:    backend.Mod(ev.mods),
		Repeat:  ev.repeat != 0,
		X:       int(ev.x),
		Y:       int(ev.y),
		Wheel:   int(ev.wheel),
		Width:   int(ev.width),
		Height:  int(ev.height),
		Focused: ev.focused != 0,
	}
	if out.Type == backend.EventText {
		out.Rune = rune(ev.codepoint)
		out.Text = string(out.Rune)
	}
	// A drop carries how many files there are and not the names; the names
	// are read off the window, which is holding them.
	if out.Type == backend.EventDropFiles {
		out.Files = d.droppedFiles(d.handle, out.Wheel)
		out.Wheel = 0
	}
	// The Cocoa side has already grown its own buffer; the canvas has to
	// follow before anything is drawn into it.
	if out.Type == backend.EventResize {
		win.ResizeCanvas(out.Width, out.Height)
	}
	return out
}
