//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "github.com/gabrielluizsf/antui/backend"

// pushKeyText pushes the text a key stands for, composing a waiting dead
// accent if there is one.
func (x *Driver) pushKeyText(win backend.Face, sym uint32, mods backend.Mod) {
	var r rune
	switch {
	case sym >= 0x01000000: // the Unicode keysym block
		r = rune(sym & 0x00FFFFFF)
	case sym >= 0x20 && sym <= 0xFF && sym != 0x7F:
		r = rune(sym)
	}

	if x.deadKey != 0 && r != 0 {
		if composed := compose(x.deadKey, r); composed != 0 {
			r = composed
		}
		x.deadKey = 0
	}

	// Ctrl-C is a command, not the letter C. Alt is left alone because some
	// layouts put real characters behind it.
	if r != 0 && mods&backend.ModControl == 0 {
		win.Push(backend.Event{
			Type: backend.EventText,
			Rune: r,
			Text: string(r),
			Mods: mods,
		})
	}
}
