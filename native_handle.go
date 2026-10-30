package antui

// NativeHandle is what another renderer — OpenGL, Vulkan, WebGPU — needs to
// draw straight into a window. Which fields mean anything depends on System.
type NativeHandle struct {
	System NativeSystem

	// Window is the window itself: the XID on X11, the HWND on Win32.
	Window uintptr

	// Display is the X11 display the window lives on, as DISPLAY spells it.
	// The window was made on antui's own connection, so a renderer opens a
	// connection of its own to this display (XOpenDisplay, xcb_connect) and
	// draws into the XID through it.
	Display string

	// Instance is the module the window class was registered from, the
	// HINSTANCE a Win32 surface is created with.
	Instance uintptr
}

// NativeSystem names the window system behind a [NativeHandle].
type NativeSystem uint8

const (
	NativeNone  NativeSystem = iota // no window, or none the platform hands out
	NativeX11                       // Window is an XID, Display names the server
	NativeWin32                     // Window is an HWND, Instance its HINSTANCE
)
