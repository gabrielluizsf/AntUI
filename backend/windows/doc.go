//go:build windows

// Package windows drives windows on Win32 (user32 + gdi32), driven through syscall.
//
// The window core of AntUI lives elsewhere; this package is the platform
// half. It receives the window to drive as a [backend.Face], which is the
// half of Window a platform is allowed to touch, and drives the display with
// the exported calls below.
//
// The Windows backend talks to user32 and gdi32 through syscall, so there is
// no cgo here and no C toolchain needed to build. The frame is handed to GDI
// as a top-down 32-bit DIB, which is exactly the layout a Canvas already
// holds — no conversion per frame.
package windows
