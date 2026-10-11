//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"
)

// answerSelection hands over what was copied. It reports whether the event
// was one of the selection ones.
func (x *Driver) answerSelection(packet []byte) bool {
	switch packet[0] & 0x7F {
	case xSelectionClear:
		// Somebody else copied something: what this window holds is no longer
		// the clipboard, and keeping it would be keeping a lie.
		x.copied = ""
		return true

	case xSelectionRequest:
		// The fields, in the order the protocol puts them: the time, who is
		// asking, which selection, what they want it as, and where to put it.
		when := binary.LittleEndian.Uint32(packet[4:])
		requestor := binary.LittleEndian.Uint32(packet[12:])
		target := binary.LittleEndian.Uint32(packet[20:])
		property := binary.LittleEndian.Uint32(packet[24:])
		if property == 0 {
			// An old-style requestor: the target is where the answer goes.
			property = target
		}

		switch {
		case target == x.atomTargets:
			// What this window can give: the list itself, and the three ways
			// of asking for text that are still in use.
			kinds := le32(x.atomTargets, x.atomUTF8String, atomString, x.atomText)
			x.changePropertyOn(requestor, property, atomAtom, 32, kinds, 4)

		case target == x.atomUTF8String || target == atomString || target == x.atomText:
			x.changePropertyOn(requestor, property, target, 8,
				[]byte(x.copied), uint32(len(x.copied)))

		default:
			property = 0 // nothing this window can turn it into
		}
		x.sendSelectionNotify(requestor, property, target, when)
		return true
	}
	return false
}
