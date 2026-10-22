//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "encoding/binary"

// gcForeground and gcBackground are bits in a CreateGC request's value-mask,
// picking out which of the two colours in its value-list is which.
const (
	gcForeground = 1 << 2
	gcBackground = 1 << 3
)

// createGCFor makes a graphics context bound to drawable and returns its id
// — [Driver.createGC]'s window-bound GC again, but for a pixmap whose depth
// does not match the window's, which a GC has to agree with to draw into it
// at all. Foreground is set to 1 and background to 0, the reverse of a
// fresh GC's own defaults: [Driver.putBitmap] sends a cursor's shape as
// literal 1-and-0 pixels, through a GC whose foreground and background
// decide what those bits actually become, and the default GC has them the
// other way round — every 1 would come out 0 and the cursor would be its
// own negative.
func (x *Driver) createGCFor(drawable uint32) uint32 {
	gc := x.newID()
	body := make([]byte, 20)
	binary.LittleEndian.PutUint32(body[0:], gc)
	binary.LittleEndian.PutUint32(body[4:], drawable)
	binary.LittleEndian.PutUint32(body[8:], gcForeground|gcBackground)
	binary.LittleEndian.PutUint32(body[12:], 1) // foreground
	binary.LittleEndian.PutUint32(body[16:], 0) // background
	x.send(xCreateGC, 0, body)
	return gc
}
