//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid
#include <android/native_window.h>
*/
import "C"

import (
	"fmt"
	"syscall"
	"unsafe"
)

// Format is how the pixels in a window's buffer are laid out. The names are
// Android's, and Android names them in **byte order**, not in the order the
// bits sit in a word — so RGBA8888 is R,G,B,A in memory, which on a
// little-endian machine reads back as 0xAABBGGRR. A canvas.Color is
// 0xAARRGGBB, so the two differ by a red/blue swap and every blit pays for
// it. This is worth knowing before debugging a picture that comes out blue.
type Format int32

// The formats a window buffer can be asked for.
const (
	RGBA8888 Format = 1
	RGBX8888 Format = 2 // the same bytes, alpha ignored
	RGB888   Format = 3 // three bytes, no padding; rare and slow
	RGB565   Format = 4
)

// String names the pixel format the way Android does.
func (f Format) String() string {
	switch f {
	case RGBA8888:
		return "RGBA_8888"
	case RGBX8888:
		return "RGBX_8888"
	case RGB888:
		return "RGB_888"
	case RGB565:
		return "RGB_565"
	}
	return fmt.Sprintf("format(%d)", int32(f))
}

// Bytes is how much one pixel takes.
func (f Format) Bytes() int {
	switch f {
	case RGB565:
		return 2
	case RGB888:
		return 3
	}
	return 4
}

// Window is an ANativeWindow: the surface the system gives an app to draw
// into. It is only ever created by Android and handed to a callback; this
// package never makes one.
//
// A Window is valid only between the system saying it exists and the system
// saying it does not. Drawing into one after that is a crash, not an error,
// which is why [Window.Acquire] exists.
type Window struct {
	ptr *C.ANativeWindow
}

// Wrap takes the pointer a native-activity callback was given. It does not
// acquire: the caller owns whatever lifetime the system promised.
func Wrap(p unsafe.Pointer) *Window {
	if p == nil {
		return nil
	}
	return &Window{ptr: (*C.ANativeWindow)(p)}
}

// Pointer is the raw handle, for passing between packages that each have
// their own idea of what a C type is.
func (w *Window) Pointer() unsafe.Pointer {
	if w == nil {
		return nil
	}
	return unsafe.Pointer(w.ptr)
}

// Acquire adds a reference, so the window outlives the callback that
// announced it. Every Acquire needs a Release.
func (w *Window) Acquire() { C.ANativeWindow_acquire(w.ptr) }

// Release drops a reference taken by Acquire.
func (w *Window) Release() { C.ANativeWindow_release(w.ptr) }

// Width and Height are the size in pixels of the buffers this window hands
// out — which is the size asked for in SetGeometry when one was asked for,
// and the physical size of the surface otherwise.
func (w *Window) Width() int { return int(C.ANativeWindow_getWidth(w.ptr)) }

// Height is the surface in pixels, top to bottom.
func (w *Window) Height() int { return int(C.ANativeWindow_getHeight(w.ptr)) }

// Format is the layout of the pixels in the buffer.
func (w *Window) Format() Format { return Format(C.ANativeWindow_getFormat(w.ptr)) }

// SetGeometry asks for buffers of a given size and format. A width and
// height of zero means "whatever the surface is", which is what an app that
// draws at native resolution wants; giving a smaller size is how a game
// renders at half resolution and lets the compositor scale it up for free.
func (w *Window) SetGeometry(width, height int, f Format) error {
	rc := C.ANativeWindow_setBuffersGeometry(w.ptr, C.int32_t(width), C.int32_t(height), C.int32_t(f))
	return wrapErr("ANativeWindow_setBuffersGeometry", int32(rc))
}

// Rect is a rectangle in window pixels, half-open on the right and bottom.
type Rect struct{ Left, Top, Right, Bottom int32 }

// Buffer is a locked window buffer: memory owned by the compositor that the
// app may write to until it posts.
type Buffer struct {
	Width, Height int
	// Stride is the distance between rows **in pixels, not bytes**, and it
	// is not the width — the compositor rounds it up, often to 16 or 32.
	// Ignoring it shears every frame diagonally.
	Stride int
	Format Format
	Bits   unsafe.Pointer
}

// Pixels32 is the buffer as 32-bit words, Stride*Height of them. It is the
// compositor's memory, not Go's, and stops being valid the moment the window
// is unlocked.
func (b Buffer) Pixels32() []uint32 {
	if b.Bits == nil || b.Format.Bytes() != 4 {
		return nil
	}
	return unsafe.Slice((*uint32)(b.Bits), b.Stride*b.Height)
}

// Pixels16 is the buffer as 16-bit words, for RGB565.
func (b Buffer) Pixels16() []uint16 {
	if b.Bits == nil || b.Format != RGB565 {
		return nil
	}
	return unsafe.Slice((*uint16)(b.Bits), b.Stride*b.Height)
}

// Row is one row of the buffer as bytes, whatever the format.
func (b Buffer) Row(y int) []byte {
	if b.Bits == nil || y < 0 || y >= b.Height {
		return nil
	}
	px := b.Format.Bytes()
	start := y * b.Stride * px
	return unsafe.Slice((*byte)(b.Bits), b.Stride*b.Height*px)[start : start+b.Width*px]
}

// Lock takes the next buffer to draw into. When dirty is not nil it names
// the region about to be written, and Lock rewrites it with the region the
// compositor actually requires — which is usually larger, because it has to
// include whatever the previous frames in the swap chain left behind. Code
// that trusts its own rectangle instead of the one that comes back leaves
// stale pixels on screen every second or third frame.
func (w *Window) Lock(dirty *Rect) (Buffer, error) {
	var buf C.ANativeWindow_Buffer
	var rect *C.ARect
	if dirty != nil {
		rect = (*C.ARect)(unsafe.Pointer(dirty))
	}
	if rc := C.ANativeWindow_lock(w.ptr, &buf, rect); rc != 0 {
		return Buffer{}, wrapErr("ANativeWindow_lock", int32(rc))
	}
	return Buffer{
		Width:  int(buf.width),
		Height: int(buf.height),
		Stride: int(buf.stride),
		Format: Format(buf.format),
		Bits:   buf.bits,
	}, nil
}

// UnlockAndPost hands the buffer back and puts it on screen.
func (w *Window) UnlockAndPost() error {
	return wrapErr("ANativeWindow_unlockAndPost", int32(C.ANativeWindow_unlockAndPost(w.ptr)))
}

// wrapErr turns the NDK's convention — zero for success, a negative errno
// for anything else — into an error that says which call failed.
func wrapErr(call string, rc int32) error {
	if rc == 0 {
		return nil
	}
	if rc < 0 {
		return fmt.Errorf("ndk: %s: %w", call, syscall.Errno(-rc))
	}
	return fmt.Errorf("ndk: %s failed (%d)", call, rc)
}
