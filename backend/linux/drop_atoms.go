//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// dropAtoms are looked up once, the first time they are wanted: a window that
// nobody ever drags anything onto should not pay for six round trips.
func (x *Driver) dropAtoms() {
	if x.atomXdndAware != 0 {
		return
	}
	x.atomXdndAware = x.internAtom("XdndAware")
	x.atomXdndEnter = x.internAtom("XdndEnter")
	x.atomXdndPosition = x.internAtom("XdndPosition")
	x.atomXdndStatus = x.internAtom("XdndStatus")
	x.atomXdndLeave = x.internAtom("XdndLeave")
	x.atomXdndDrop = x.internAtom("XdndDrop")
	x.atomXdndFinished = x.internAtom("XdndFinished")
	x.atomXdndSelection = x.internAtom("XdndSelection")
	x.atomXdndActionCopy = x.internAtom("XdndActionCopy")
	x.atomUriList = x.internAtom("text/uri-list")
	x.atomClipboard = x.internAtom("CLIPBOARD")
	x.atomDropProperty = x.internAtom("ANTUI_DROP")
	x.atomTargets = x.internAtom("TARGETS")
	x.atomText = x.internAtom("TEXT")
}
