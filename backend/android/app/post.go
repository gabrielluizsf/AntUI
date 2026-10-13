//go:build android

package app

import "unsafe"

// post hands an event to the app's goroutine. When wait is set it does not
// return until the app has acknowledged it, which is what makes it safe for
// Android to take the window away on the next line.
func post(ev Event, ptr unsafe.Pointer, wait bool) {
	a := current()
	if a == nil {
		return
	}
	m := message{ev: ev, ptr: ptr}
	if wait {
		m.ack = make(chan struct{})
	}
	select {
	case a.msgs <- m:
	case <-a.dead:
		return
	}
	if m.ack == nil {
		return
	}
	select {
	case <-m.ack:
	case <-a.dead:
	}
}
