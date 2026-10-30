package antui

// NativeHandle is the platform's own handle to the window, for a renderer of
// its own — OpenGL, Vulkan, WebGPU — to draw into. ok is false when there is
// no window or the platform has no handle to give; the software canvas is
// the way to draw there.
//
// A renderer drawing through it wants [Window.SetExternal] too, or the next
// frame's canvas lands on top of what it drew.
func (win *Window) NativeHandle() (h NativeHandle, ok bool) {
	if win.native == nil {
		return NativeHandle{}, false
	}
	return win.native.nativeHandle()
}
