//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// sendDragStatus tells the source whether we will take what it is dragging.
func (x *Driver) sendDragStatus(accept bool) {
	if x.dragFrom == 0 {
		return
	}
	accepted := uint32(0)
	action := uint32(0)
	if accept {
		accepted = 1
		action = x.atomXdndActionCopy
	}
	x.sendDragMessage(x.atomXdndStatus, [5]uint32{
		x.window, accepted, 0, 0, action,
	})
}
