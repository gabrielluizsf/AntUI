//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"errors"
	"syscall"

	"github.com/gabrielluizsf/antui/backend"
)

// open connects to the server, creates the window and maps it.
func (x *Driver) open(win backend.Face, title string, width, height int) error {
	x.win = win

	conn, displayNumber, err := dialDisplay()
	if err != nil {
		return err
	}
	x.conn = conn
	// Raw access to the socket, for reading without waiting. Without it the
	// backend still works, a millisecond at a time.
	if sc, ok := conn.(interface {
		SyscallConn() (syscall.RawConn, error)
	}); ok {
		x.raw, _ = sc.SyscallConn()
	}
	if err := x.handshake(displayNumber); err != nil {
		conn.Close()
		x.alive = false
		return err
	}

	x.internWindowAtoms()
	// Say the window will take a drag. Without this property no source
	// offers it anything, so a window that does not ask never finds out that
	// anybody was dragging.
	x.announceDrops()
	x.loadKeymap()

	x.createWindow(title, width, height)
	x.createGC()
	x.send(xMapWindow, 0, le32(x.window))

	// The display's scale is read now, so the first frame finds the number
	// here instead of paying a round trip for it, and watched for after:
	// the desktop rewrites the resource database when the user changes the
	// scale, and that write is this end's way of hearing about it.
	x.startScale()

	if !x.alive {
		return errors.New("antui: the connection to the X server dropped while opening the window")
	}
	return nil
}
