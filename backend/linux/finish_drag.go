//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// finishDrag tells the source we are done with the drop, which is what
// releases it to carry on.
func (x *Driver) finishDrag() {
	if x.dragFrom == 0 {
		return
	}
	x.sendDragMessage(x.atomXdndFinished, [5]uint32{
		x.window, 1, x.atomXdndActionCopy, 0, 0,
	})
	x.dragFrom = 0
}
