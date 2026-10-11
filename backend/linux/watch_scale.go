//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// watchScale asks the server to say when a property of the root window is
// written, which is how a desktop announces the display scale it has just
// changed: Xft.dpi lives in the resource database on the root window, and
// rewriting that database is the write. The selection is kept per
// connection, so nothing else on the display is told anything it did not
// ask for, and only the database's own atom drops the cached number — every
// other property the window manager writes on the root passes through.
//
// A desktop that changed the scale without writing the database would go
// unnoticed, but the value is read from there, so a change to it has to be
// said from there too: this is the whole of the signal X11 carries.
func (x *Driver) watchScale() {
	if !x.alive || x.root == 0 {
		return
	}
	x.send(xChangeWindowAttributes, 0, le32(x.root, cwEventMask, xPropertyChangeMask))
}
