//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "github.com/gabrielluizsf/antui/backend"

// handleEvent decodes one 32-byte event packet and folds it into the window.
func (x *Driver) handleEvent(packet []byte) {
	win := x.win
	if win == nil {
		return
	}
	// The top bit marks an event the server sent on someone's behalf rather
	// than one it generated; either way the body reads the same.
	switch packet[0] & 0x7F {
	case xError:
		// An error outside a request we are waiting on is informational: the
		// window manager rejecting a hint is not a reason to stop drawing.

	case xKeyPress, xKeyRelease:
		x.handleKey(packet)

	case xButtonPress, xButtonRelease:
		x.handleButton(packet)

	case xMotionNotify:
		x.handleMotion(packet)

	case xExpose, xMapNotify:
		// The window has just appeared, which means everything in it has to
		// be drawn. Some servers send an Expose for that and some do not —
		// a bare X server with no window manager is the case that does not —
		// so this says it either way rather than leaving a program that has
		// only just opened waiting to be told to draw. Pushing the expose is
		// the whole of it: that is what marks the frame as a full redraw.
		win.Push(backend.Event{Type: backend.EventExpose})

	case xConfigureNotify:
		x.handleConfigure(packet)

	case xSelectionClear, xSelectionRequest:
		x.answerSelection(packet)

	case xClientMessage:
		x.handleClientMessage(packet)

	case xFocusIn, xFocusOut:
		win.Push(backend.Event{Type: backend.EventFocus, Focused: packet[0]&0x7F == xFocusIn})

	case xMappingNotify:
		// The layout changed under us; the old table would now type the
		// wrong letters.
		x.loadKeymap()

	case xPropertyNotify:
		x.handlePropertyNotify(packet)
	}
}
