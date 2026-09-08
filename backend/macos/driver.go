//go:build darwin

// Package macos drives windows on Cocoa (AppKit) through the Objective-C runtime.
//
// The window core of AntUI lives elsewhere; this package is the platform
// half. It receives the window to drive as a [backend.Face], which is the
// half of Window a platform is allowed to touch, and drives the display with
// the exported calls below.
//
// macOS is the one platform that needs a C toolchain: AppKit is
// Objective-C, and there is no socket or DLL to talk to it through. The
// Cocoa side lives in backend_darwin.m; driver.go is the translation between
// its event struct and the package's own.
package macos

/*
#cgo CFLAGS: -x objective-c -Wno-deprecated-declarations
#cgo LDFLAGS: -framework Cocoa -framework CoreGraphics -framework CoreVideo
#include <stdlib.h>
#include "backend_darwin.h"
*/
import "C"

import (
	"errors"
	"runtime"
	"unsafe"

	"github.com/gabrielluizsf/antui/backend"
	"github.com/gabrielluizsf/antui/canvas"
)

// Driver is one window on macOS. It owns the platform connections the way
// the window core owns the input state.
type Driver struct {
	handle *C.antui_d_window
	win    backend.Face
	events [64]C.antui_d_event
}

// Open creates the window that Options describes.
//
// AppKit insists on the main thread. Locking pins this goroutine to the
// thread it is already on, so a program that opens its window from main gets
// the main thread and everything works; one that opens it from another
// goroutine at least keeps every AppKit call on one thread.
func Open(win backend.Face, opts backend.Options) (*Driver, error) {
	if win == nil {
		return nil, errors.New("antui: macos backend with no window to drive")
	}
	runtime.LockOSThread()

	d := &Driver{win: win}
	cTitle := C.CString(opts.Title)
	defer C.free(unsafe.Pointer(cTitle))

	d.handle = C.antui_d_open(cTitle, C.int(opts.Width), C.int(opts.Height))
	if d.handle == nil {
		runtime.UnlockOSThread()
		return nil, errors.New("antui: could not create the window; on macOS this has to happen on the main thread")
	}
	return d, nil
}

// Close destroys the window and lets go of everything the platform held.
func (d *Driver) Close() {
	if d.handle == nil {
		return
	}
	C.antui_d_close(d.handle)
	d.handle = nil
	runtime.UnlockOSThread()
}

// Pump handles whatever happened this frame, folding it into the window
// through win.Push.
func (d *Driver) Pump(win backend.Face) {
	if d.handle == nil {
		return
	}
	// Loop until the queue comes back short, so a burst bigger than the
	// buffer is drained rather than left to arrive a frame late.
	for {
		n := int(C.antui_d_pump(d.handle, &d.events[0], C.int(len(d.events))))
		for i := range n {
			win.Push(d.decode(win, &d.events[i]))
		}
		if n < len(d.events) {
			return
		}
	}
}

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

// Present copies the window's dirty rectangle to the platform surface.
func (d *Driver) Present(win backend.Face, dirty canvas.Area) {
	if d.handle == nil || dirty.Width <= 0 || dirty.Height <= 0 {
		return
	}
	cv := win.Canvas()
	if cv.Pixels == nil {
		return
	}
	// The pixels are copied on the other side, not retained, so handing a Go
	// pointer across for the length of this call is within the cgo rules.
	C.antui_d_present(d.handle,
		(*C.uint32_t)(unsafe.Pointer(&cv.Pixels[0])),
		C.int(cv.Width), C.int(cv.Height),
		C.int(dirty.X), C.int(dirty.Y), C.int(dirty.Width), C.int(dirty.Height))
}

// SetTitle names the window in the system's decoration.
func (d *Driver) SetTitle(title string) {
	if d.handle == nil {
		return
	}
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	C.antui_d_set_title(d.handle, cTitle)
}

// SetFullscreen asks the system for a full-screen window, reporting whether
// the request could be made at all.
func (d *Driver) SetFullscreen(on bool) bool {
	if d.handle == nil {
		return false
	}
	flag := C.int(0)
	if on {
		flag = 1
	}
	return C.antui_d_set_fullscreen(d.handle, flag) != 0
}

// DisplaySize is the size of the display the window is on.
func (d *Driver) DisplaySize() (w, h int, ok bool) {
	var width, height C.int
	if C.antui_d_display_size(d.handle, &width, &height) == 0 {
		return 0, 0, false
	}
	return int(width), int(height), true
}

// DisplayRefresh is how often the display repaints, in hertz.
func (d *Driver) DisplayRefresh() int {
	return int(C.antui_d_display_refresh(d.handle))
}

// SetLimits passes the size constraints on to the platform.
func (d *Driver) SetLimits(limits backend.Limits) {
	if d.handle == nil {
		return
	}
	fixed, noMaximize := C.int(0), C.int(0)
	if limits.Fixed {
		fixed = 1
	}
	if limits.NoMaximize {
		noMaximize = 1
	}
	C.antui_d_set_limits(d.handle,
		C.int(limits.MinWidth), C.int(limits.MinHeight),
		C.int(limits.MaxWidth), C.int(limits.MaxHeight),
		C.double(limits.Aspect), fixed, noMaximize)
}

// SetSize asks the platform for a new drawable size, reporting whether the
// request could be made.
func (d *Driver) SetSize(win backend.Face, width, height int) bool {
	if d.handle == nil {
		return false
	}
	return C.antui_d_set_size(d.handle, C.int(width), C.int(height)) != 0
}

// ContentScale is pixels per point. On macOS it is 1, and that is the true
// answer rather than a missing one.
//
// The view reports its bounds in points and the framebuffer is that size, so
// a window asked for at 1280 by 720 is 1280 by 720 pixels of Canvas and
// appears at 1280 by 720 points on every Mac — the postage-stamp problem the
// other two systems have simply does not arise here, because AppKit does the
// scaling on the way to the glass.
//
// What it costs is sharpness: on a retina display the frame is resampled up
// rather than drawn at the panel's own resolution. Drawing at
// convertRectToBacking size would fix that and double the pixels pushed
// every frame, which is a decision for a machine that can measure it.
func (d *Driver) ContentScale() float64 { return 1 }

// Clipboard is what the system clipboard holds: its text, and the files it
// names.
func (d *Driver) Clipboard() (text string, files []string) {
	if held := C.antui_d_clipboard_text(); held != nil {
		text = C.GoString(held)
		C.free(unsafe.Pointer(held))
	}

	for i := range int(C.antui_d_clipboard_count()) {
		if path := C.antui_d_clipboard_path(C.int(i)); path != nil {
			files = append(files, C.GoString(path))
		}
	}
	return text, files
}

// SetIcon gives the window an icon where the system shows one. The largest
// picture is handed to AppKit, which is the one the Dock wants: macOS has no
// small icon of its own to set.
func (d *Driver) SetIcon(images []*canvas.Canvas) bool {
	if d.handle == nil || len(images) == 0 {
		return false
	}
	best := images[0]
	for _, image := range images {
		if image.Width > best.Width {
			best = image
		}
	}
	pixels := make([]uint32, best.Width*best.Height)
	for y := range best.Height {
		for x := range best.Width {
			pixels[y*best.Width+x] = uint32(best.At(x, y))
		}
	}
	C.antui_d_set_icon(d.handle, (*C.uint32_t)(unsafe.Pointer(&pixels[0])),
		C.int(best.Width), C.int(best.Height))
	return true
}

// SetClipboard puts text on the system clipboard, reporting whether the
// system took it.
func (d *Driver) SetClipboard(text string) bool {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.antui_d_set_clipboard(cText)
	return true
}

// droppedFiles reads the paths of the drop that has just been reported. The
// window holds them until the next drop, so this is safe to do while the
// event is being decoded.
func (d *Driver) droppedFiles(handle *C.antui_d_window, count int) []string {
	if handle == nil || count <= 0 {
		return nil
	}
	held := int(C.antui_d_drop_count(handle))
	if held < count {
		count = held
	}
	files := make([]string, 0, count)
	for i := range count {
		if path := C.antui_d_drop_path(handle, C.int(i)); path != nil {
			files = append(files, C.GoString(path))
		}
	}
	return files
}