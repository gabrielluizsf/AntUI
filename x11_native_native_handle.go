//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

func (n *x11Native) nativeHandle() (NativeHandle, bool) {
	window, display := n.d.Native()
	if window == 0 {
		return NativeHandle{}, false
	}
	return NativeHandle{System: NativeX11, Window: uintptr(window), Display: display}, true
}
