//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// SetClipboard puts text on the X selections, taking the clipboard over and
// keeping the text to hand out.
func (x *Driver) SetClipboard(text string) bool {
	if !x.alive || x.window == 0 {
		return false
	}
	x.dropAtoms()
	if x.atomClipboard == 0 {
		return false
	}
	if x.atomTargets == 0 {
		x.atomTargets = x.internAtom("TARGETS")
		x.atomText = x.internAtom("TEXT")
	}

	x.copied = text
	// SetSelectionOwner: window, selection, time. CurrentTime is nought,
	// which the protocol allows and every toolkit uses for a copy that
	// happens now.
	x.send(xSetSelectionOwner, 0, le32(x.window, x.atomClipboard, 0))
	return x.alive
}
