//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"

	"github.com/gabrielluizsf/antui/backend"
)

// X keysyms this backend recognises by name.
const (
	symBackspace = 0xFF08
	symTab       = 0xFF09
	symReturn    = 0xFF0D
	symKPEnter   = 0xFF8D
	symEscape    = 0xFF1B
	symHome      = 0xFF50
	symLeft      = 0xFF51
	symUp        = 0xFF52
	symRight     = 0xFF53
	symDown      = 0xFF54
	symPageUp    = 0xFF55
	symPageDown  = 0xFF56
	symEnd       = 0xFF57
	symInsert    = 0xFF63
	symDelete    = 0xFFFF
	symISOLevel3 = 0xFF7E // the AltGr some layouts send instead
	symF1        = 0xFFBE
	symF12       = 0xFFC9
	symKP0       = 0xFFB0
	symKP9       = 0xFFB9

	symShiftL = 0xFFE1
	symShiftR = 0xFFE2
	symCtrlL  = 0xFFE3
	symCtrlR  = 0xFFE4
	symCaps   = 0xFFE5
	symAltL   = 0xFFE9
	symAltR   = 0xFFEA
	symSuperL = 0xFFEB
	symSuperR = 0xFFEC

	symDeadLo = 0xFE50 // the dead-accent block
	symDeadHi = 0xFE5F
)

// The state mask the server sends with a key or button event.
const (
	stateShift  = 0x1
	stateLock   = 0x2 // caps lock
	stateCtrl   = 0x4
	stateAlt    = 0x8
	stateSuper  = 0x40
	stateAltGr  = 0x80
	stateGroup2 = 0x2000
)

// keysym resolves a keycode to the symbol it means with these modifiers held.
func (x *Driver) keysym(keycode byte, state uint16) uint32 {
	if x.keysyms == nil || keycode < x.minKeycode || keycode > x.maxKeycode {
		return 0
	}
	perCode := int(x.keysymsPerCode)
	base := (int(keycode) - int(x.minKeycode)) * perCode
	if base < 0 || base >= len(x.keysyms) {
		return 0
	}

	index := 0
	if state&stateAltGr != 0 && perCode > 2 {
		index = 2
	}
	if state&stateShift != 0 {
		index++
	}
	index = min(index, perCode-1)

	sym := x.keysyms[base+index]
	// An empty slot in the shifted column means the key has only one symbol.
	if sym == 0 && index > 0 {
		sym = x.keysyms[base]
	}

	// Caps lock only affects letters, and it inverts shift rather than
	// forcing upper case — which is why holding shift with caps on types
	// lower case, as it does everywhere else.
	switch {
	case state&stateLock != 0 && state&stateShift == 0 && sym >= 'a' && sym <= 'z':
		sym -= 32
	case state&stateLock != 0 && state&stateShift != 0 && sym >= 'A' && sym <= 'Z':
		sym += 32
	}
	return sym
}

// keysymToKey maps a symbol to the physical key it stands for. Printable
// keys report their uppercase ASCII code, so the caller can compare against
// KeyA regardless of the layout's case.
func keysymToKey(sym uint32) backend.Key {
	switch {
	case sym >= 'a' && sym <= 'z':
		return backend.Key(sym - 32)
	case sym >= 'A' && sym <= 'Z', sym >= '0' && sym <= '9':
		return backend.Key(sym)
	}

	switch sym {
	case ' ':
		return backend.KeySpace
	case symBackspace:
		return backend.KeyBackspace
	case symTab:
		return backend.KeyTab
	case symReturn, symKPEnter:
		return backend.KeyEnter
	case symEscape:
		return backend.KeyEscape
	case symHome:
		return backend.KeyHome
	case symLeft:
		return backend.KeyLeft
	case symUp:
		return backend.KeyUp
	case symRight:
		return backend.KeyRight
	case symDown:
		return backend.KeyDown
	case symPageUp:
		return backend.KeyPageUp
	case symPageDown:
		return backend.KeyPageDown
	case symEnd:
		return backend.KeyEnd
	case symInsert:
		return backend.KeyInsert
	case symDelete:
		return backend.KeyDelete
	case symShiftL:
		return backend.KeyLeftShift
	case symShiftR:
		return backend.KeyRightShift
	case symCtrlL:
		return backend.KeyLeftControl
	case symCtrlR:
		return backend.KeyRightControl
	case symCaps:
		return backend.KeyCapsLock
	case symAltL:
		return backend.KeyLeftAlt
	case symAltR, symISOLevel3:
		return backend.KeyRightAlt
	case symSuperL:
		return backend.KeyLeftSuper
	case symSuperR:
		return backend.KeyRightSuper
	}

	switch {
	case sym >= symF1 && sym <= symF12:
		return backend.KeyF1 + backend.Key(sym-symF1)
	case sym >= symKP0 && sym <= symKP9:
		return backend.Key0 + backend.Key(sym-symKP0)
	case sym < 0x80:
		return backend.Key(sym) // ASCII punctuation
	}
	return backend.KeyUnknown
}

// deadPair keys the compose table by the accent and the letter it lands on.
type deadPair struct {
	dead uint32
	base rune
}

// composeDead is what makes an ABNT2 keyboard work: "~" then "a" is one
// character, not two. Only the accents a Latin-1 font can actually draw are
// here, because composing to something the font has no glyph for would put a
// '?' on screen where the user typed a letter.
var composeDead = map[deadPair]rune{
	// grave
	{0xFE50, 'a'}: 0xE0, {0xFE50, 'e'}: 0xE8, {0xFE50, 'i'}: 0xEC,
	{0xFE50, 'o'}: 0xF2, {0xFE50, 'u'}: 0xF9,
	{0xFE50, 'A'}: 0xC0, {0xFE50, 'E'}: 0xC8, {0xFE50, 'I'}: 0xCC,
	{0xFE50, 'O'}: 0xD2, {0xFE50, 'U'}: 0xD9,
	// acute
	{0xFE51, 'a'}: 0xE1, {0xFE51, 'e'}: 0xE9, {0xFE51, 'i'}: 0xED,
	{0xFE51, 'o'}: 0xF3, {0xFE51, 'u'}: 0xFA, {0xFE51, 'y'}: 0xFD,
	{0xFE51, 'A'}: 0xC1, {0xFE51, 'E'}: 0xC9, {0xFE51, 'I'}: 0xCD,
	{0xFE51, 'O'}: 0xD3, {0xFE51, 'U'}: 0xDA,
	// circumflex
	{0xFE52, 'a'}: 0xE2, {0xFE52, 'e'}: 0xEA, {0xFE52, 'i'}: 0xEE,
	{0xFE52, 'o'}: 0xF4, {0xFE52, 'u'}: 0xFB,
	{0xFE52, 'A'}: 0xC2, {0xFE52, 'E'}: 0xCA, {0xFE52, 'I'}: 0xCE,
	{0xFE52, 'O'}: 0xD4, {0xFE52, 'U'}: 0xDB,
	// tilde
	{0xFE53, 'a'}: 0xE3, {0xFE53, 'n'}: 0xF1, {0xFE53, 'o'}: 0xF5,
	{0xFE53, 'A'}: 0xC3, {0xFE53, 'N'}: 0xD1, {0xFE53, 'O'}: 0xD5,
	// diaeresis
	{0xFE57, 'a'}: 0xE4, {0xFE57, 'e'}: 0xEB, {0xFE57, 'i'}: 0xEF,
	{0xFE57, 'o'}: 0xF6, {0xFE57, 'u'}: 0xFC,
	{0xFE57, 'A'}: 0xC4, {0xFE57, 'E'}: 0xCB, {0xFE57, 'I'}: 0xCF,
	{0xFE57, 'O'}: 0xD6, {0xFE57, 'U'}: 0xDC,
	// cedilla
	{0xFE59, 'c'}: 0xE7, {0xFE59, 'C'}: 0xC7,
}

// deadAlone is the accent on its own, which is what a space after it means.
var deadAlone = map[uint32]rune{
	0xFE50: '`', 0xFE51: '\'', 0xFE52: '^', 0xFE53: '~', 0xFE57: '"',
}

// compose joins a pending dead accent to the character that followed it,
// returning 0 when the pair does not make a letter this font can draw.
func compose(dead uint32, base rune) rune {
	if r, ok := composeDead[deadPair{dead, base}]; ok {
		return r
	}
	if base == ' ' {
		return deadAlone[dead]
	}
	return 0
}

func modsFromState(state uint16) backend.Mod {
	var mods backend.Mod
	if state&stateShift != 0 {
		mods |= backend.ModShift
	}
	if state&stateCtrl != 0 {
		mods |= backend.ModControl
	}
	if state&stateAlt != 0 {
		mods |= backend.ModAlt
	}
	if state&stateSuper != 0 {
		mods |= backend.ModSuper
	}
	return mods
}

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
		press := packet[0]&0x7F == xButtonPress
		state := binary.LittleEndian.Uint16(packet[28:])
		ev := backend.Event{
			Mods: modsFromState(state),
			X:    int(int16(binary.LittleEndian.Uint16(packet[24:]))),
			Y:    int(int16(binary.LittleEndian.Uint16(packet[26:]))),
		}
		switch button := packet[1]; button {
		case 4, 5: // the wheel arrives as a button, and only on press
			if press {
				ev.Type = backend.EventMouseWheel
				ev.Wheel = 1
				if button == 5 {
					ev.Wheel = -1
				}
				win.Push(ev)
			}
			return
		case 1:
			ev.Button = backend.MouseLeft
		case 2:
			ev.Button = backend.MouseMiddle
		case 3:
			ev.Button = backend.MouseRight
		default:
			return
		}
		ev.Type = backend.EventMouseUp
		if press {
			ev.Type = backend.EventMouseDown
		}
		win.Push(ev)

	case xMotionNotify:
		win.Push(backend.Event{
			Type: backend.EventMouseMove,
			X:    int(int16(binary.LittleEndian.Uint16(packet[24:]))),
			Y:    int(int16(binary.LittleEndian.Uint16(packet[26:]))),
			Mods: modsFromState(binary.LittleEndian.Uint16(packet[28:])),
		})

	case xExpose:
		win.Push(backend.Event{Type: backend.EventExpose})

	case xMapNotify:
		// The window has just appeared, which means everything in it has to
		// be drawn. Some servers send an Expose for that and some do not —
		// a bare X server with no window manager is the case that does not —
		// so this says it either way rather than leaving a program that has
		// only just opened waiting to be told to draw. Pushing the expose is
		// the whole of it: that is what marks the frame as a full redraw.
		win.Push(backend.Event{Type: backend.EventExpose})

	case xConfigureNotify:
		// Where the window is, which a drag needs: XDND gives the pointer's
		// place on the root and everything here counts from the window.
		x.windowX = int(int16(binary.LittleEndian.Uint16(packet[16:])))
		x.windowY = int(int16(binary.LittleEndian.Uint16(packet[18:])))

		w := int(binary.LittleEndian.Uint16(packet[20:]))
		h := int(binary.LittleEndian.Uint16(packet[22:]))
		cv := win.Canvas()
		if w > 0 && h > 0 && (w != cv.Width || h != cv.Height) {
			if win.ResizeCanvas(w, h) {
				win.Push(backend.Event{Type: backend.EventResize, Width: w, Height: h})
			}
		}

	case xSelectionClear, xSelectionRequest:
		x.answerSelection(packet)

	case xClientMessage:
		typeAtom := binary.LittleEndian.Uint32(packet[8:])
		data0 := binary.LittleEndian.Uint32(packet[12:])
		if typeAtom == x.atomWMProtocols && data0 == x.atomWMDeleteWindow {
			win.PushSimple(backend.EventClose)
			return
		}
		x.handleDrop(packet)

	case xFocusIn, xFocusOut:
		win.Push(backend.Event{Type: backend.EventFocus, Focused: packet[0]&0x7F == xFocusIn})

	case xMappingNotify:
		// The layout changed under us; the old table would now type the
		// wrong letters.
		x.loadKeymap()
	}
}

// handleKey turns one key packet into a key event and, when the key stands
// for a character, the text event that goes with it.
func (x *Driver) handleKey(packet []byte) {
	win := x.win
	if win == nil {
		return
	}
	press := packet[0]&0x7F == xKeyPress
	state := binary.LittleEndian.Uint16(packet[28:])
	sym := x.keysym(packet[1], state)

	ev := backend.Event{
		Type: backend.EventKeyUp,
		Key:  keysymToKey(sym),
		Mods: modsFromState(state),
		X:    int(int16(binary.LittleEndian.Uint16(packet[24:]))),
		Y:    int(int16(binary.LittleEndian.Uint16(packet[26:]))),
	}
	if press {
		ev.Type = backend.EventKeyDown
	}
	win.Push(ev)

	if !press {
		return
	}

	// A dead accent produces no text of its own: it waits for the next key.
	if sym >= symDeadLo && sym <= symDeadHi {
		x.deadKey = sym
		return
	}

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
	if r != 0 && ev.Mods&backend.ModControl == 0 {
		win.Push(backend.Event{
			Type: backend.EventText,
			Rune: r,
			Text: string(r),
			Mods: ev.Mods,
		})
	}
}