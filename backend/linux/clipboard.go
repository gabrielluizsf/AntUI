//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// Clipboard is what the X selections hold: the files they name, and their
// text.
func (x *Driver) Clipboard() (string, []string) {
	if !x.alive {
		return "", nil
	}
	x.dropAtoms()

	// Files first: a file manager puts both on the clipboard, and the files
	// are the more useful answer when there are any.
	if data := x.askSelection(x.atomClipboard, x.atomUriList, 0); len(data) > 0 {
		if paths := parseURIList(data); len(paths) > 0 {
			return string(data), paths
		}
	}
	if x.atomUTF8String != 0 {
		if data := x.askSelection(x.atomClipboard, x.atomUTF8String, 0); len(data) > 0 {
			return string(data), parseURIList(data)
		}
	}
	data := x.askSelection(x.atomClipboard, atomString, 0)
	return string(data), parseURIList(data)
}
