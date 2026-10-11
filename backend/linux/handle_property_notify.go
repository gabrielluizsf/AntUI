//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "encoding/binary"

// handlePropertyNotify marks the display scale stale when the resource
// database on the root window is rewritten, which is the only root property
// here that is not the window's own business.
func (x *Driver) handlePropertyNotify(packet []byte) {
	// The desktop wrote a property on the root window. The resource database
	// being rewritten is a display scale the user just changed, and the
	// frame that follows reads it again. Every other property on the root —
	// the window manager writes several — falls through untouched.
	if binary.LittleEndian.Uint32(packet[4:]) == x.root &&
		binary.LittleEndian.Uint32(packet[8:]) == xResourceManager {
		x.scaleDirty = true
	}
}
