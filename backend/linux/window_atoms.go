//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// internWindowAtoms looks up the properties the window core sets on every
// window, once, so the frame loop never pays a round trip for them.
func (x *Driver) internWindowAtoms() {
	x.atomWMProtocols = x.internAtom("WM_PROTOCOLS")
	x.atomWMDeleteWindow = x.internAtom("WM_DELETE_WINDOW")
	x.atomNetWMName = x.internAtom("_NET_WM_NAME")
	x.atomUTF8String = x.internAtom("UTF8_STRING")
	x.atomNetWMState = x.internAtom("_NET_WM_STATE")
	x.atomNetWMStateFullscreen = x.internAtom("_NET_WM_STATE_FULLSCREEN")
	x.atomNetWMOpacity = x.internAtom("_NET_WM_WINDOW_OPACITY")
}
