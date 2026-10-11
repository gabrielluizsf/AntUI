//go:build windows

package windows

import (
	"github.com/gabrielluizsf/antui/backend"
)

// handleChar turns a UTF-16 code unit into a text event, joining the two
// halves of a surrogate pair into the one character they stand for.
func (d *Driver) handleChar(unit uint16) uintptr {
	var r rune
	switch {
	case unit >= 0xD800 && unit <= 0xDBFF: // the high half
		d.highSurrogate = unit
		return 0
	case unit >= 0xDC00 && unit <= 0xDFFF: // the low half
		if d.highSurrogate != 0 {
			r = 0x10000 + (rune(d.highSurrogate)-0xD800)<<10 + (rune(unit) - 0xDC00)
			d.highSurrogate = 0
		}
	default:
		r = rune(unit)
	}

	// Below 32 is a control character, and 127 is delete: neither is text.
	if r >= 32 && r != 127 {
		mods := winMods()
		if mods&backend.ModControl == 0 {
			d.win.Push(backend.Event{Type: backend.EventText, Rune: r, Text: string(r), Mods: mods})
		}
	}
	return 0
}
